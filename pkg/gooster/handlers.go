package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/convert"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/events"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

func (app *App) handleExitEvent(events.Event) {
	app.Log().Info("Stopping app")
	if err := app.AppContext.Close(); err != nil {
		app.Log().Error(errors.WithMessage(err, "stopping app"))
	}
	app.root.Stop()
}

func (app *App) handleSetFocusEvent(event events.Event) {
	if event.Payload == nil {
		return
	}

	if view, ok := event.Payload.(tview.Primitive); ok {
		app.Log().DebugF("Focusing view: %T", view)
		app.root.SetFocus(view)
		app.lastFocus = view
	} else {
		app.Log().ErrorF("gooster.Actions.SetFocus event expects tview.Primitive as payload. Found %T", event.Payload)
	}
}

const dialogPageId = "gooster_dialog_box"

func (app *App) handleEventOpenDialog(event events.Event) {
	dlg := event.Payload.(dialog.Dialog)
	view := dlg.View(app.cfg.Dialog, func(form *tview.Form) {
		app.AppActions().CloseDialog()
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

func (app *App) handleEventCloseDialog(event events.Event) {
	if !app.pages.HasPage(dialogPageId) {
		return
	}

	app.Log().Debug("Closing dialog")
	app.pages.RemovePage(dialogPageId)
	app.AppActions().SetFocus(app.lastFocus)
}

func (app *App) handleEventAddTab(event events.Event) {
	tab, ok := event.Payload.(Tab)
	if !ok {
		app.Log().ErrorF("gooster.Actions.AddTab event expects gooster.Tab as payload. Found %T", event.Payload)
		return
	}

	pageId := tab.pageId()
	if app.pages.HasPage(pageId) {
		app.Log().ErrorF("Can not add tab with ID '%s'. The ID must be unique, but such tab already exists.", tab.Id)
		return
	}

	if tab.View == nil {
		tab.View = app.createMainGrid()
	}

	// @todo: implement tab title

	app.log.DebugF("Creating a new tab '%s'", tab.Id)
	app.pages.AddPage(pageId, tab.View, true, true)
}

func (app *App) handleEventShowTab(event events.Event) {
	tabId := convert.ToString(event.Payload)
	pageId := Tab{Id: tabId}.pageId()
	if !app.pages.HasPage(pageId) {
		app.Log().ErrorF("Could not show tab with ID '%s'. Not found.", tabId)
		return
	}
	app.pages.ShowPage(pageId)
}

const initialTabId = "initial"

func (app *App) handleEventRemoveTab(event events.Event) {
	tabId := convert.ToString(event.Payload)
	if tabId == initialTabId {
		app.Log().Warn("Can not remove the initial tab.")
		return
	}

	pageId := Tab{Id: tabId}.pageId()
	if !app.pages.HasPage(pageId) {
		app.Log().ErrorF("Could not remove tab with ID '%s'. Not found.", tabId)
		return
	}

	app.pages.RemovePage(pageId)
}

func (app *App) handleKeyCtrlC(event *tcell.EventKey) *tcell.EventKey {
	app.Log().Debug("Interrupting latest command")
	app.Events().Dispatch(events.Event{Id: "command:interrupt"}) // @todo use Actions
	return nil
}

func (app *App) handleKeyEscape(event *tcell.EventKey) *tcell.EventKey {
	app.AppActions().CloseDialog()
	return event
}

func (app *App) handleKeyExit(event *tcell.EventKey) *tcell.EventKey {
	app.AppActions().Exit()
	return nil
}
