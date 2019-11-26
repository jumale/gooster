package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/oliveagle/jsonpath"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"reflect"
)

type jsonMap = map[string]interface{}
type jsonArr = []interface{}
type yamlMap = map[interface{}]interface{}

const idKey = "#id"

type YamlReaderConfig struct {
	Fs       filesys.FileSys
	Defaults io.Reader
}

type YamlReader struct {
	cfg  YamlReaderConfig
	fs   filesys.FileSys
	data jsonMap
}

func NewYamlReader(cfg YamlReaderConfig) *YamlReader {
	fs := cfg.Fs
	if fs == nil {
		fs = filesys.Default{}
	}
	return &YamlReader{cfg: cfg, fs: fs}
}

func (r *YamlReader) Read(jsonPath string, target interface{}) error {
	return readJsonPath(r.data, jsonPath, target)
}

func (r *YamlReader) Load(yamlCfg io.Reader) (err error) {
	var config, defaults yamlMap
	if err = yaml.NewDecoder(yamlCfg).Decode(&config); err != nil {
		return errors.WithMessage(err, "Could not decode config")
	}
	if r.cfg.Defaults != nil {
		if err = yaml.NewDecoder(r.cfg.Defaults).Decode(&defaults); err != nil {
			return errors.WithMessage(err, "Could not decode default config")
		}
	}

	r.data = mergeConfigs(
		yamlValToJsonVal(defaults).(jsonMap),
		yamlValToJsonVal(config).(jsonMap),
	)

	return nil
}

func (r *YamlReader) LoadString(yamlCfg string) error {
	return r.Load(bytes.NewBufferString(yamlCfg))
}

func (r *YamlReader) LoadFile(path string) (err error) {
	var config io.Reader
	if config, err = r.fs.Open(path); err != nil {
		return errors.WithMessagef(err, "Could not open config file %s", path)
	}
	return r.Load(config)
}

func readJsonPath(data interface{}, path string, target interface{}) error {
	result, err := jsonpath.JsonPathLookup(data, path)
	if err != nil {
		return errors.WithMessagef(err, "Could not read jsonPath '%s'", path)
	}

	encoded := bytes.NewBuffer(nil)
	if err = json.NewEncoder(encoded).Encode(result); err != nil {
		return errors.WithMessagef(err, "Could not encode data from jsonPath '%s' to JSON format", path)
	}

	if err = json.NewDecoder(encoded).Decode(target); err != nil {
		return errors.WithMessagef(err, "Could not decode data from jsonPath '%s' to the target", path)
	}
	return nil
}

func mergeConfigs(a jsonMap, b jsonMap) jsonMap {
	result := make(jsonMap)

	// copy "a"
	for key := range a {
		result[key] = a[key]
	}

	// apply "b" on top
	for key := range b {
		// just set it, if there is no such value yet
		if _, ok := result[key]; !ok {
			result[key] = b[key]
			continue
		}
		// just set it if the old value is nil
		if result[key] == nil {
			result[key] = b[key]
			continue
		}
		// do nothing if the new value is nil
		if b[key] == nil {
			continue
		}
		// just set it if new value has a different type
		if reflect.ValueOf(result[key]).Type() != reflect.ValueOf(b[key]).Type() {
			result[key] = b[key]
			continue
		}

		switch b[key].(type) {
		case jsonMap:
			result[key] = mergeConfigs(a[key].(jsonMap), b[key].(jsonMap))
		case jsonArr:
			aVal, aIsArrOfMaps := arrToArrOfMaps(a[key].(jsonArr))
			bVal, bIsArrOfMaps := arrToArrOfMaps(b[key].(jsonArr))
			if aIsArrOfMaps && bIsArrOfMaps {
				result[key] = mergeArrOfMaps(aVal, bVal)
			} else {
				result[key] = b[key]
			}

		default:
			result[key] = b[key]
		}

	}
	return result
}

func yamlValToJsonVal(val interface{}) interface{} {
	switch typedVal := val.(type) {
	case jsonMap:
		newMap := jsonMap{}
		for key, value := range typedVal {
			newMap[key] = yamlValToJsonVal(value)
		}
		return newMap

	case jsonArr:
		newArr := make(jsonArr, len(typedVal))
		for i := range typedVal {
			newArr[i] = yamlValToJsonVal(typedVal[i])
		}
		return newArr

	case yamlMap:
		newMap := jsonMap{}
		for key, value := range typedVal {
			newMap[fmt.Sprintf("%v", key)] = yamlValToJsonVal(value)
		}
		return newMap

	default:
		return val
	}
}

func arrToArrOfMaps(a jsonArr) (result []jsonMap, success bool) {
	result = make([]jsonMap, len(a))
	for i := range a {
		switch v := a[i].(type) {
		case jsonMap:
			result[i] = v
		default:
			return nil, false
		}
	}
	return result, true
}

func mergeArrOfMaps(a []jsonMap, b []jsonMap) []jsonMap {
	indexById := make(map[string]int)
	//noinspection GoPreferNilSlice
	result := []jsonMap{}

	for i := range a {
		result = append(result, a[i])
		if id, ok := a[i][idKey]; ok {
			indexById[fmt.Sprintf("%v", id)] = i
		}
	}

	for i := range b {
		if idx, ok := indexById[fmt.Sprintf("%v", b[i][idKey])]; ok {
			result[idx] = mergeConfigs(result[idx], b[i])
		} else {
			result = append(result, b[i])
		}
	}

	return result
}
