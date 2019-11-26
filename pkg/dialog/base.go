package dialog

import (
	"github.com/jumale/gooster/pkg/log"
	"github.com/rivo/tview"
)

const (
	buttonHeight = 1
	//minHeight    = 2
	minWidth  = 30
	maxHeight = 10
	maxWidth  = 50
)

type Base struct {
	Title   string
	Border  bool
	Buttons []Button
	Log     log.Logger
}

func (b Base) CreateBox(content tview.Primitive, cfg Config, onDone ActionHandler) tview.Primitive {
	box := tview.NewGrid()
	box.SetBackgroundColor(cfg.Colors.Bg.Origin())
	if b.Border {
		box.SetBorder(true)
	}
	if b.Title != "" {
		box.SetTitle(" " + b.Title + " ")
		box.SetBorder(true)
	}

	var form *tview.Form
	switch content.(type) {
	case *tview.Form:
		form = content.(*tview.Form)
		content = nil
	default:
		form = tview.NewForm()
		form.SetRect(0, 0, 0, 0)
	}

	width := 0
	height := 0
	if content != nil && form != nil {
		box.SetColumns(-1)
		box.SetRows(-1, buttonHeight)
		box.AddItem(content, 0, 0, 1, 1, 0, 0, false)
		box.AddItem(form, 1, 0, 1, 1, 0, 0, true)
		b.initForm(form, cfg, onDone)
		width, height = stackRows(content, form)

		cw, ch := viewSize(content)
		fw, fh := viewSize(form)
		b.debug("content size: %dx%d, form size: %dx%d, sum: %dx%d", cw, ch, fw, fh, width, height)

	} else if content != nil {
		box.SetColumns(-1)
		box.SetRows(-1)
		box.AddItem(content, 0, 0, 1, 1, 0, 0, false)
		width, height = viewSize(content)
		b.debug("content size: %dx%d", width, height)

	} else if form != nil {
		box.SetColumns(-1)
		box.SetRows(-1)
		box.AddItem(form, 0, 0, 1, 1, 0, 0, true)
		b.initForm(form, cfg, onDone)
		width, height = viewSize(form)
		b.debug("form size: %dx%d", width, height)
	}

	if b.Border || b.Title != "" {
		width += 2
		height += 2
	}
	b.debug("base box size: %dx%d", width, height)
	box.SetRect(0, 0, width, height)

	return box
}

func (b Base) initForm(form *tview.Form, cfg Config, onDone ActionHandler) {
	form.SetBorderPadding(0, 0, 0, 0)
	form.SetBackgroundColor(cfg.Colors.Bg.Origin())

	if len(b.Buttons) > 0 {
		_, _, w, h := form.GetRect()
		form.SetRect(0, 0, w, h+buttonHeight)

		form.SetButtonBackgroundColor(cfg.Colors.Bg.Origin())
		form.SetButtonsAlign(tview.AlignRight)

		for _, btn := range b.Buttons {
			form.AddButton(btn.Label, combineFunc(form, btn.Action, onDone))
			if btn.Focus {
				form.SetFocus(form.GetFormItemIndex(btn.Label))
			}
		}
	}
}

func (b Base) debug(msg string, args ...interface{}) {
	if b.Log != nil {
		b.Log.DebugF(msg, args...)
	}
}

func combineFunc(form *tview.Form, fn ...ActionHandler) func() {
	return func() {
		for _, item := range fn {
			if item != nil {
				item(form)
			}
		}
	}
}

func viewSize(view tview.Primitive) (width int, height int) {
	_, _, width, height = view.GetRect()

	if width < minWidth {
		width = minWidth
	}
	if width > maxWidth {
		width = maxWidth
	}

	if height > maxHeight {
		height = maxHeight
	}

	return width, height
}

func stackRows(rows ...tview.Primitive) (width int, height int) {
	for _, view := range rows {
		w, h := viewSize(view)
		height += h
		if width == 0 || width < w {
			width = w
		}
	}
	return width, height
}
