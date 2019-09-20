package dialog

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Dialog interface {
	View(cfg Config, onDone ActionHandler) tview.Primitive
}

type Config struct {
	Colors ColorsConfig
}

type ColorsConfig struct {
	Bg             tcell.Color
	BtnColor       tcell.Color
	BtnActiveColor tcell.Color
}

type Context interface {
	CloseDialog()
}
