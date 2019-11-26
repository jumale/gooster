package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rivo/tview"
)

type Align int

func (a *Align) UnmarshalJSON(b []byte) error {
	var val string
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}

	switch val {
	case "left":
		*a = tview.AlignLeft
	case "right":
		*a = tview.AlignRight
	case "center", "":
		*a = tview.AlignCenter
	default:
		return errors.New(fmt.Sprintf("Unknown align value '%s'. Expected: left, right, center", val))
	}
	return nil
}

func (a Align) Origin() int {
	return int(a)
}
