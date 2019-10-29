package gooster

import (
	"fmt"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/events"
	"github.com/rivo/tview"
)

const (
	ActionExit        events.EventId = "app:exit"
	ActionSetFocus                   = "app:set_focus"
	ActionDraw                       = "app:draw"
	ActionOpenDialog                 = "app:open_dialog"
	ActionCloseDialog                = "app:close_dialog"
	ActionAddTab                     = "app:add_tab"
	ActionShowTab                    = "app:show_tab"
	ActionRemoveTab                  = "app:remove_tab"
)

type Actions struct {
	ctx *AppContext
}

func (a Actions) Exit() {
	a.ctx.Events().Dispatch(events.Event{Id: ActionExit})
}

func (a Actions) Draw() {
	a.ctx.Events().Dispatch(events.Event{Id: ActionDraw})
}

func (a Actions) SetFocus(view tview.Primitive) {
	a.ctx.Events().Dispatch(events.Event{Id: ActionSetFocus, Payload: view})
}

func (a Actions) OpenDialog(diag dialog.Dialog) {
	a.ctx.Events().Dispatch(events.Event{Id: ActionOpenDialog, Payload: diag})
}

func (a Actions) CloseDialog() {
	a.ctx.Events().Dispatch(events.Event{Id: ActionCloseDialog})
}

func (a Actions) AddTab(tab Tab) {
	a.ctx.Events().Dispatch(events.Event{Id: ActionAddTab, Payload: tab})
}

func (a Actions) ShowTab(tabId string) {
	a.ctx.Events().Dispatch(events.Event{Id: ActionShowTab, Payload: tabId})
}

func (a Actions) RemoveTab(tabId string) {
	a.ctx.Events().Dispatch(events.Event{Id: ActionRemoveTab, Payload: tabId})
}

type Tab struct {
	Id    string
	Title string
	View  tview.Primitive
}

func (tab Tab) pageId() string {
	return fmt.Sprintf("gooster_tab_%s", tab.Id)
}
