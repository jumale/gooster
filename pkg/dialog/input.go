package dialog

import (
	"github.com/jumale/gooster/pkg/log"
	"github.com/rivo/tview"
)

type Input struct {
	Title  string
	Label  string
	Value  string
	Width  int
	Border bool
	OnOk   func(val string)
	Log    log.Logger
}

func (d Input) View(cfg Config, onDone ActionHandler) tview.Primitive {
	inputText := ""
	if d.Width == 0 {
		d.Width = 20
	}

	form := tview.NewForm()
	form.SetRect(0, 0, d.Width+len(d.Label)+1, 2)
	form.AddInputField(d.Label, d.Value, d.Width, nil, func(text string) {
		inputText = text
	})

	return Base{
		Title:  d.Title,
		Border: d.Border,
		Buttons: []Button{
			{
				Label: "OK",
				Action: func(form *tview.Form) {
					d.OnOk(inputText)
				},
				Focus: true,
			},
		},
		Log: d.Log,
	}.CreateBox(form, cfg, onDone)
}
