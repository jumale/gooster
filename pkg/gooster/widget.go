package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Widget interface {
	Name() string
	Init(*AppContext) (tview.Primitive, WidgetConfig, error)
}

type Position struct {
	Col    int `json:"col"`
	Row    int `json:"row"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type WidgetConfig struct {
	Position `json:",inline"`
	Enabled  bool `json:"enabled"`
	Focused  bool `json:"focused"`
	FocusKey tcell.Key
}
