package config

import (
	"encoding/json"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Color tcell.Color

func (c *Color) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}

	if val, ok := tcell.ColorNames[name]; ok {
		*c = Color(val)
		return nil
	}

	hex := name
	if strings.HasPrefix(name, "#") {
		hex = hex[1:]
	}
	if hexInt, err := strconv.ParseInt(hex, 16, 32); err == nil {
		*c = Color(tcell.NewHexColor(int32(hexInt)))
		return nil
	}

	return errors.Errorf("Can not unmarshal color '%s' into tcell.Color. Invalid value.", name)
}

func (c Color) Origin() tcell.Color {
	return tcell.Color(c)
}
