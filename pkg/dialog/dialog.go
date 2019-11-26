package dialog

import (
	"github.com/jumale/gooster/pkg/config"
	"github.com/rivo/tview"
)

type Dialog interface {
	View(cfg Config, onDone ActionHandler) tview.Primitive
}

type Config struct {
	Colors ColorsConfig `json:"colors"`
}

type ColorsConfig struct {
	Bg        config.Color `json:"bg"`
	Btn       config.Color `json:"btn"`
	BtnActive config.Color `json:"btn_active"`
}

type Context interface {
	CloseDialog()
}
