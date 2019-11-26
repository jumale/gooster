package config

import (
	"encoding/json"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
)

var keyValues = getKeyValues()

type Key tcell.Key

func (k *Key) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}

	val, ok := keyValues[name]
	if !ok {
		return errors.Errorf("Can not unmarshal key name '%s' into tcell.Key. Such key does not exist.", name)
	}
	*k = Key(val)
	return nil
}

func (k Key) Origin() tcell.Key {
	return tcell.Key(k)
}

func getKeyValues() map[string]tcell.Key {
	vals := make(map[string]tcell.Key)
	for key, name := range tcell.KeyNames {
		vals[name] = key
	}
	return vals
}
