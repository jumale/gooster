package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Module interface {
	Name() string
	Init(*AppContext) (tview.Primitive, ModuleConfig, error)
}

type Position struct {
	Col    int `json:"col"`
	Row    int `json:"row"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ModuleConfig struct {
	Position `json:",inline"`
	Enabled  bool `json:"enabled"`
	Focused  bool `json:"focused"`
	FocusKey tcell.Key
}

type moduleDefinition struct {
	module Module
	view   tview.Primitive
	cfg    ModuleConfig
}
