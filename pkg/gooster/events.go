package gooster

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/command"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/rivo/tview"
)

type DelayedEventManager interface {
	Init() error
}

type DrawableEvent interface {
	NeedsDraw() bool
}

// ------------------------------------------------------------ //

type EventExit struct{}

type EventInterrupt struct{}

type EventSetFocus struct {
	Target tview.Primitive
}

type EventSetFocusByName struct {
	TargetName string
}

type EventDraw struct{}

type EventOutput struct {
	Data []byte
}

func (e EventOutput) NeedsDraw() bool {
	return true
}

type EventSetCompletion struct {
	Input      string
	Commands   []command.Definition
	Completion command.Completion
}

type EventOpenDialog struct {
	Dialog dialog.Dialog
}

type EventCloseDialog struct{}

type EventAddTab struct {
	Id    string
	Title string
	View  tview.Primitive
}

func (e EventAddTab) pageId() string {
	return fmt.Sprintf("gooster_tab_%s", e.Id)
}

type EventShowTab struct {
	TabId string
}

type EventRemoveTab struct {
	TabId string
}

// ------------------------------------------------------------ //

type KeyEventHandler func(event *tcell.EventKey) *tcell.EventKey
type KeyEventHandlers map[config.Key]KeyEventHandler

func HandleKeyEvents(target handlerGetter, handlers KeyEventHandlers) {
	keyMap := make(map[[3]int16]KeyEventHandler)
	for k, handler := range handlers {
		keyMap[keyDef(k.Type, k.Rune, k.Mod)] = handler
	}

	prev := target.GetInputCapture()
	capture := func(ev *tcell.EventKey) *tcell.EventKey {
		if prev != nil {
			if ev = prev(ev); ev == nil {
				return nil
			}
		}

		if handler, ok := keyMap[keyDef(ev.Key(), ev.Rune(), ev.Modifiers())]; ok {
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

func keyDef(key tcell.Key, char rune, mod tcell.ModMask) [3]int16 {
	return [3]int16{int16(key), int16(char), int16(mod)}
}
