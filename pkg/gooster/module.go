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

type KeyEventHandler func(event *tcell.EventKey) *tcell.EventKey
type KeyEventHandlers map[tcell.Key]KeyEventHandler

func HandleKeyEvents(target handlerGetter, handlers KeyEventHandlers) {
	prev := target.GetInputCapture()
	capture := func(ev *tcell.EventKey) *tcell.EventKey {
		if prev != nil {
			if ev = prev(ev); ev == nil {
				return nil
			}
		}

		if handler, ok := handlers[ev.Key()]; ok {
			if ev = handler(ev); ev == nil {
				return nil
			}
		}

		return ev
	}

	switch t := target.(type) {
	case appHandlerSetter:
		t.SetInputCapture(capture)
	case boxHandlerSetter:
		t.SetInputCapture(capture)
	}
}
