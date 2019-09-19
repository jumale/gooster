package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/events"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type dialogManger struct {
	cfg       dialog.Config
	ctx       *AppContext
	pages     *tview.Pages
	pageId    string
	hasDialog bool
}

func newDialogManger(cfg dialog.Config, ctx *AppContext, pages *tview.Pages) *dialogManger {
	mng := &dialogManger{
		cfg:    cfg,
		ctx:    ctx,
		pages:  pages,
		pageId: "gooster_dialog_box",
	}

	ctx.actions.Subscribe(ctx.actions.withSeed(eventOpenDialog), events.Subscriber{
		Handler: func(event events.Event) {
			mng.open(event.Data.(dialog.Dialog))
		},
	})

	ctx.actions.Subscribe(ctx.actions.withSeed(eventCloseDialog), events.Subscriber{
		Handler: func(event events.Event) {
			if err := mng.Close(); err != nil {
				ctx.log.Error(errors.WithMessage(err, "closing dialog"))
			}
		},
	})

	return mng
}

func (mng *dialogManger) open(dialog dialog.Dialog) {
	width, height := dialog.Size()

	modal := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(dialog.View(mng.cfg), 1, 1, 1, 1, 0, 0, true)

	modal.SetBackgroundColor(tcell.ColorDefault)

	mng.pages.AddPage(mng.pageId, modal, true, true)
	mng.hasDialog = true
}

func (mng *dialogManger) Close() error {
	mng.pages.RemovePage(mng.pageId)
	mng.hasDialog = false
	return nil
}
