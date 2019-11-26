package testtools

import (
	"github.com/pkg/errors"
	"reflect"
)

type ConfigReader struct {
	stubs map[string]interface{}
}

func (c *ConfigReader) Read(jsonPath string, target interface{}) error {
	cfg, ok := c.stubs[jsonPath]
	if !ok {
		return errors.Errorf("Could not find config stub for path '%s'", jsonPath)
	}
	if cfg == nil {
		return nil
	}
	if err, ok := cfg.(error); ok {
		return err
	}

	reflect.ValueOf(target).Elem().Set(reflect.ValueOf(cfg))
	return nil
}

func (c *ConfigReader) ShouldReturn(cfgName string, val interface{}) {
	c.stubs[cfgName] = val
}
