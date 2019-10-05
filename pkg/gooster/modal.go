package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/events"
	"github.com/rivo/tview"
)

type modalManger struct {
	cfg     dialog.Config
	ctx     *AppContext
	pages   *tview.Pages
	pageId  string
	isOpen  bool
	onClose func()
}

func newModalManger(cfg dialog.Config, ctx *AppContext, pages *tview.Pages) *modalManger {
	mng := &modalManger{
		cfg:    cfg,
		ctx:    ctx,
		pages:  pages,
		pageId: "gooster_dialog_box",
	}

	ctx.actions.Subscribe(ctx.actions.withSeed(actionAppOpenDialog), events.Subscriber{
		Handler: func(event events.Event) {
			dg := event.Data.(dialog.Dialog)
			view := dg.View(mng.cfg, func(form *tview.Form) {
				mng.ctx.actions.CloseDialog()
			})
			_, _, width, height := view.GetRect()

			modal := tview.NewGrid().
				SetColumns(0, width, 0).
				SetRows(0, height, 0).
				AddItem(view, 1, 1, 1, 1, 5, 5, true)
			modal.SetBackgroundColor(tcell.ColorDefault)

			mng.pages.AddPage(mng.pageId, modal, true, true)
			mng.isOpen = true
		},
	})

	ctx.actions.Subscribe(ctx.actions.withSeed(actionAppCloseDialog), events.Subscriber{
		Handler: func(event events.Event) {
			mng.pages.RemovePage(mng.pageId)
			mng.isOpen = false
			if mng.onClose != nil {
				mng.onClose()
			}
		},
	})

	return mng
}
