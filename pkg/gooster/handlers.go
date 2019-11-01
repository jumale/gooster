package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

func (app *App) handleExitEvent() {
	app.Log().Info("Stopping app")
	if err := app.AppContext.Close(); err != nil {
		app.Log().Error(errors.WithMessage(err, "stopping app"))
	}
	app.root.Stop()
}

func (app *App) handleDrawEvent() {
	app.root.Draw()
}

func (app *App) handleSetFocusEvent(event EventSetFocus) {
	if event.Target != nil {
		app.Log().DebugF("Focusing view: %T", event.Target)
		app.root.SetFocus(event.Target)
		app.lastFocus = event.Target
	}
}

const dialogPageId = "gooster_dialog_box"

func (app *App) handleEventOpenDialog(event EventOpenDialog) {
	view := event.Dialog.View(app.cfg.Dialog, func(form *tview.Form) {
		app.Events().Dispatch(EventCloseDialog{})
	})
	_, _, width, height := view.GetRect()

	modal := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(view, 1, 1, 1, 1, 5, 5, true)
	modal.SetBackgroundColor(tcell.ColorDefault)

	app.Log().Debug("Opening dialog")
	app.pages.AddPage(dialogPageId, modal, true, true)
}

func (app *App) handleEventCloseDialog() {
	if !app.pages.HasPage(dialogPageId) {
		return
	}
	app.Log().Debug("Closing dialog")
	app.pages.RemovePage(dialogPageId)
	app.Events().Dispatch(EventSetFocus{Target: app.lastFocus})
}

func (app *App) handleEventAddTab(event EventAddTab) {
	pageId := event.pageId()
	if app.pages.HasPage(pageId) {
		app.Log().ErrorF("Can not add tab with ID '%s'. The ID must be unique, but such tab already exists.", event.Id)
		return
	}

	if event.View == nil {
		event.View = app.createMainGrid()
	}

	// @todo: implement tab title

	app.log.DebugF("Creating a new tab '%s'", event.Id)
	app.pages.AddPage(pageId, event.View, true, true)
}

func (app *App) handleEventShowTab(event EventShowTab) {
	tabId := event.TabId
	pageId := EventAddTab{Id: tabId}.pageId()
	if !app.pages.HasPage(pageId) {
		app.Log().ErrorF("Could not show tab with ID '%s'. Not found.", tabId)
		return
	}
	app.pages.ShowPage(pageId)
}

const initialTabId = "initial"

func (app *App) handleEventRemoveTab(event EventRemoveTab) {
	tabId := event.TabId
	if tabId == initialTabId {
		app.Log().Warn("Can not remove the initial tab.")
		return
	}

	pageId := EventAddTab{Id: tabId}.pageId()
	if !app.pages.HasPage(pageId) {
		app.Log().ErrorF("Could not remove tab with ID '%s'. Not found.", tabId)
		return
	}

	app.pages.RemovePage(pageId)
}

func (app *App) handleKeyCtrlC(_ *tcell.EventKey) *tcell.EventKey {
	app.Log().Debug("Interrupting latest command")
	app.Events().Dispatch(EventInterrupt{})
	return nil
}

func (app *App) handleKeyEscape(event *tcell.EventKey) *tcell.EventKey {
	app.Events().Dispatch(EventCloseDialog{})
	return event
}

func (app *App) handleKeyExit(_ *tcell.EventKey) *tcell.EventKey {
	app.Events().Dispatch(EventExit{})
	return nil
}
