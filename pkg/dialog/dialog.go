package dialog

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Dialog interface {
	View(cfg Config) tview.Primitive
	Size() (width int, height int)
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
