package dialog

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Dialog interface {
	View(cfg Config, onDone ActionHandler) tview.Primitive
}

type Config struct {
	Colors ColorsConfig `json:"colors"`
}

type ColorsConfig struct {
	Bg        tcell.Color `json:"bg"`
	Btn       tcell.Color `json:"btn"`
	BtnActive tcell.Color `json:"btn_active"`
}

type Context interface {
	CloseDialog()
}
