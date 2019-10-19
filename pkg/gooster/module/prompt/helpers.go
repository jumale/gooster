package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
)

func getColorName(c tcell.Color) string {
	for name, value := range tcell.ColorNames {
		if value == c {
			return name
		}
	}
	return "black"
}

func (m *Module) check(err error, msg ...string) {
	if err == nil {
		return
	}
	if len(msg) > 0 {
		m.Log().Error(errors.WithMessage(err, msg[0]))
	} else {
		m.Log().Error(err)
	}
}
