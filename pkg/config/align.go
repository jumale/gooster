package config

import (
	"errors"
	"fmt"
	"github.com/rivo/tview"
)

type Align int

func (a *Align) UnmarshalJSON(b []byte) error {
	val := string(b)
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
