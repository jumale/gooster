package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/config"
	"github.com/rivo/tview"
)

type ModuleConfig struct {
	Position `json:",inline"`
	Enabled  bool       `json:"enabled"`
	Focused  bool       `json:"focused"`
	FocusKey config.Key `json:"focus_key"`
}

type Position struct {
	Col    int `json:"col"`
	Row    int `json:"row"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ------------------------------------------------------------ //

var defaultModConfig = ModuleConfig{
	Enabled: true,
	Position: Position{
		Width:  10,
		Height: 10,
	},
}

// ------------------------------------------------------------ //

type Module interface {
	Name() string
	View() ModuleView
	Init(Context) error
}

type ModuleView interface {
	tview.Primitive
	tview.Boxed
}

type handlerGetter interface {
	GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey
}

type appHandlerSetter interface {
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Application
}

type boxHandlerSetter interface {
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Box
}
