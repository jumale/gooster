package config

import (
	"bytes"
	"github.com/jumale/gooster/pkg/filesys/fstub"
	"github.com/stretchr/testify/require"
	"testing"
)

const defaultConfig = `
app:
  name: my_app
  debug: false
modules:
  - '#id': foo
    enabled: true
  - '#id': bar
    enabled: false
`

const customConfig = `
app:
  name: custom_name
  debug: true
modules:
  - '#id': bar
    enabled: true
`

func TestYamlReader(t *testing.T) {
	assert := require.New(t)

	t.Run("Load", func(t *testing.T) {
		t.Run("should load a config file", func(t *testing.T) {
			fs := fstub.New(fstub.Config{})
			fs.Root().Add("/config.yaml", fstub.NewFile(customConfig))
			reader := NewYamlReader(YamlReaderConfig{Fs: fs, Defaults: bytes.NewBufferString(defaultConfig)})

			err := reader.LoadFile("/config.yaml")
			assert.NoError(err)

			t.Run("then Read", func(t *testing.T) {
				t.Run("should read the corresponding configs", func(t *testing.T) {
					appConfig := &appCfg{}
					err = reader.Read("$.app", appConfig)
					assert.NoError(err)
					assert.Equal("custom_name", appConfig.Name)
					assert.Equal(true, appConfig.Debug)

					fooConfig := &moduleCfg{}
					err = reader.Read("$.modules[?(@.#id == 'foo')][0]", fooConfig)
					assert.NoError(err)
					assert.Equal(true, fooConfig.Enabled)

					barConfig := &moduleCfg{}
					err = reader.Read("$.modules[?(@.#id == 'bar')][0]", barConfig)
					assert.NoError(err)
					assert.Equal(true, barConfig.Enabled)
				})
			})
		})

		// @todo: test converting yaml maps to json maps
	})
}

func TestMergeConfigs(t *testing.T) {
	assert := require.New(t)

	t.Run("should just take 'a' values, if 'b' is empty", func(t *testing.T) {
		a := jsonMap{"foo": "bar"}
		assert.Equal(a, mergeConfigs(a, nil))
	})

	t.Run("should just take 'b' values, if 'a' is empty", func(t *testing.T) {
		b := jsonMap{"foo": "bar"}
		assert.Equal(b, mergeConfigs(nil, b))
	})

	t.Run("should override values", func(t *testing.T) {
		a := jsonMap{"foo": "bar", "baz": false}
		b := jsonMap{"baz": true}
		assert.Equal(jsonMap{"foo": "bar", "baz": true}, mergeConfigs(a, b))
	})

	t.Run("should add values if not exist", func(t *testing.T) {
		a := jsonMap{"foo": "bar"}
		b := jsonMap{"baz": true}
		assert.Equal(jsonMap{"foo": "bar", "baz": true}, mergeConfigs(a, b))
	})

	t.Run("should merge sub-level maps", func(t *testing.T) {
		a := jsonMap{"foo": jsonMap{"bar": jsonMap{"baz": 123}, "cat": "sad"}}
		b := jsonMap{"foo": jsonMap{"bar": jsonMap{"baz": 456}}}
		assert.Equal(jsonMap{"foo": jsonMap{"bar": jsonMap{"baz": 456}, "cat": "sad"}}, mergeConfigs(a, b))
	})

	t.Run("should override value if types do not match", func(t *testing.T) {
		a := jsonMap{"foo": jsonMap{"bar": "baz"}}
		b := jsonMap{"foo": jsonArr{1, 2, 3}}
		assert.Equal(jsonMap{"foo": jsonArr{1, 2, 3}}, mergeConfigs(a, b))

		t.Run("and in another direction", func(t *testing.T) {
			a = jsonMap{"foo": jsonArr{1, 2, 3}}
			b = jsonMap{"foo": jsonMap{"bar": "baz"}}
			assert.Equal(jsonMap{"foo": jsonMap{"bar": "baz"}}, mergeConfigs(a, b))
		})
	})

	t.Run("should override array values", func(t *testing.T) {
		a := jsonMap{"foo": jsonArr{1, 2, 3}}
		b := jsonMap{"foo": jsonArr{4, 5}}
		assert.Equal(jsonMap{"foo": jsonArr{4, 5}}, mergeConfigs(a, b))
	})

	t.Run("should merge array of object with '#id' properties", func(t *testing.T) {
		a := jsonMap{"parent": jsonArr{
			jsonMap{idKey: "foo", "val": 111},
			jsonMap{idKey: "bar", "val": 222},
		}}
		b := jsonMap{"parent": jsonArr{
			jsonMap{idKey: "foo", "val": 333}, // should be merged with existing item in 'a'
		}}
		expected := jsonMap{"parent": []jsonMap{
			{idKey: "foo", "val": 333},
			{idKey: "bar", "val": 222},
		}}
		assert.Equal(expected, mergeConfigs(a, b))
	})

	t.Run("should override nil values", func(t *testing.T) {
		a := jsonMap{"foo": nil}
		b := jsonMap{"foo": "bar"}
		assert.Equal(jsonMap{"foo": "bar"}, mergeConfigs(a, b))
	})

	t.Run("should skip if new value is nil", func(t *testing.T) {
		a := jsonMap{"foo": "bar"}
		b := jsonMap{"foo": nil}
		assert.Equal(jsonMap{"foo": "bar"}, mergeConfigs(a, b))
	})
}

func TestMergeArrOfMaps(t *testing.T) {
	assert := require.New(t)

	t.Run("should just take 'a' values if 'b' is empty", func(t *testing.T) {
		a := []jsonMap{{"foo": "bar"}}
		assert.Equal(a, mergeArrOfMaps(a, nil))
	})

	t.Run("should just take 'b' values if 'a' is empty", func(t *testing.T) {
		b := []jsonMap{{"foo": "bar"}}
		assert.Equal(b, mergeArrOfMaps(nil, b))
	})

	t.Run("should merge items by ID", func(t *testing.T) {
		a := []jsonMap{
			{idKey: "foo", "val": 111},
			{idKey: "bar", "val": 222},
		}
		b := []jsonMap{
			{idKey: "bar", "val": 333},
		}
		expected := []jsonMap{
			{idKey: "foo", "val": 111},
			{idKey: "bar", "val": 333},
		}
		assert.Equal(expected, mergeArrOfMaps(a, b))
	})

	t.Run("should merge nested maps", func(t *testing.T) {
		a := []jsonMap{{idKey: "foo", "child": jsonMap{"val": 111, "bar": "baz"}}}
		b := []jsonMap{{idKey: "foo", "child": jsonMap{"val": 222}}}
		expected := []jsonMap{
			{idKey: "foo", "child": jsonMap{"val": 222, "bar": "baz"}},
		}
		assert.Equal(expected, mergeArrOfMaps(a, b))
	})

	t.Run("should add new values from 'b'", func(t *testing.T) {
		a := []jsonMap{
			{idKey: "foo", "val": 111},
			{idKey: "bar", "val": 222},
		}
		b := []jsonMap{
			{idKey: "baz", "val": 444},
		}
		expected := []jsonMap{
			{idKey: "foo", "val": 111},
			{idKey: "bar", "val": 222},
			{idKey: "baz", "val": 444},
		}
		assert.Equal(expected, mergeArrOfMaps(a, b))
	})

	t.Run("should add items as is, if they do not have IDs", func(t *testing.T) {
		a := []jsonMap{
			{idKey: "foo", "val": 111},
			{"val": 222},
		}
		b := []jsonMap{
			{"val": 333},
			{idKey: "bar", "val": 444},
		}
		expected := []jsonMap{
			{idKey: "foo", "val": 111},
			{"val": 222},
			{"val": 333},
			{idKey: "bar", "val": 444},
		}
		assert.Equal(expected, mergeArrOfMaps(a, b))
	})
}

func TestReadJsonPath(t *testing.T) {
	assert := require.New(t)

	t.Run("should read value to the target", func(t *testing.T) {
		data := jsonMap{
			"parent": jsonMap{
				"list": jsonArr{
					jsonMap{"name": "cat", "meows": true},
				},
			},
		}
		target := &ReadTarget{}

		err := readJsonPath(data, "$.parent.list[0]", target)
		assert.NoError(err)

		assert.Equal("cat", target.Name)
		assert.Equal(true, target.Meows)
	})

	t.Run("should find value by property", func(t *testing.T) {
		data := jsonMap{
			"list": jsonArr{
				jsonMap{"name": "cat", "meows": true},
				jsonMap{"name": "dog", "meows": false},
			},
		}
		target := &ReadTarget{}

		err := readJsonPath(data, "$.list[?(@.name == 'dog')][0]", target)
		assert.NoError(err)

		assert.Equal("dog", target.Name)
		assert.Equal(false, target.Meows)
	})

	t.Run("should keep default data if no config provided", func(t *testing.T) {
		data := jsonMap{
			"parent": jsonMap{"meows": true},
		}
		target := &ReadTarget{Name: "bat"}

		err := readJsonPath(data, "$.parent", target)
		assert.NoError(err)

		assert.Equal("bat", target.Name)
		assert.Equal(true, target.Meows)
	})
}

type ReadTarget struct {
	Name  string `json:"name"`
	Meows bool   `json:"meows"`
}

type moduleCfg struct {
	Enabled bool `json:"enabled"`
}

type appCfg struct {
	Name  string `json:"name"`
	Debug bool   `json:"debug"`
}
