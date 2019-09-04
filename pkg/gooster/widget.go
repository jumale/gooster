package gooster

import "github.com/rivo/tview"

type Widget interface {
	Name() string
	Init(*AppContext) error
	View() tview.Primitive
	Config() WidgetConfig
}

type Position struct {
	Col    int `json:"col"`
	Row    int `json:"row"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type WidgetConfig struct {
	Position `json:",inline"`
	Focused  bool `json:"focused"`
}
