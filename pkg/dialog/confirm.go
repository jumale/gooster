package dialog

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"math"
)

const buttonHeight = 3

type ConfirmDialog struct {
	Title    string
	Text     string
	Width    int
	Height   int
	Border   bool
	OnOk     func()
	OnCancel func()
	Context  Context
}

func (d ConfirmDialog) View(cfg Config) tview.Primitive {
	view := tview.NewTextView()
	view.SetBackgroundColor(tcell.ColorDefault)
	if d.Text != "" {
		view.SetText(d.Text)
	}

	halfWidth := float64(d.Width) / 2
	box := tview.NewGrid().
		SetColumns(int(math.Ceil(halfWidth)), int(math.Floor(halfWidth))).
		SetRows(d.Height, buttonHeight).
		AddItem(view, 0, 0, 1, 2, 0, 0, false).
		AddItem(d.btn("Cancel", d.OnCancel).View(cfg), 1, 0, 1, 1, 0, 0, true).
		AddItem(d.btn("OK", d.OnOk).View(cfg), 1, 1, 1, 1, 0, 0, false)

	box.SetBackgroundColor(cfg.Colors.Bg)
	if d.Border {
		box.SetBorder(true)
	}
	if d.Title != "" {
		box.SetTitle(" " + d.Title + " ")
	}

	return box
}

func (d ConfirmDialog) Size() (width int, height int) {
	return d.Width, d.Height + buttonHeight
}

func (d ConfirmDialog) btn(label string, handler func()) Button {
	return &SimpleButton{
		Label: label,
		OnClick: func() {
			handler()
			d.closeDialog()
		},
	}
}

func (d ConfirmDialog) closeDialog() {
	if d.Context != nil {
		d.Context.CloseDialog()
	}
}
