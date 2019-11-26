package gooster

import (
	"fmt"
	"github.com/jumale/gooster/pkg/command"
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
	Completion []string
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
