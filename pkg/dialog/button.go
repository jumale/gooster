package dialog

import (
	"github.com/rivo/tview"
)

type Button interface {
	View(cfg Config) tview.Primitive
}

type SimpleButton struct {
	OnClick func()
	Label   string
}

func (btn *SimpleButton) View(cfg Config) tview.Primitive {
	view := tview.NewButton(btn.Label)
	view.SetBorder(true)
	view.SetSelectedFunc(btn.OnClick)
	if cfg.Colors.BtnColor > 0 {
		view.SetBackgroundColor(cfg.Colors.BtnColor)
	}
	if cfg.Colors.BtnActiveColor > 0 {
		view.SetBackgroundColorActivated(cfg.Colors.BtnActiveColor)
	}

	return view
}
