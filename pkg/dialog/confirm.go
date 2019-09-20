package dialog

import (
	"github.com/jumale/gooster/pkg/log"
	"github.com/rivo/tview"
)

type Confirm struct {
	Title    string
	Text     string
	Width    int
	Height   int
	Border   bool
	OnOk     ActionHandler
	OnCancel ActionHandler
	FocusOk  bool
	Log      log.Logger
}

func (d Confirm) View(cfg Config, onDone ActionHandler) tview.Primitive {
	return Text{
		Title:  d.Title,
		Text:   d.Text,
		Border: d.Border,
		Buttons: []Button{
			{
				Label:  "Cancel",
				Action: d.OnCancel,
				Focus:  !d.FocusOk,
			},
			{
				Label:  "OK",
				Action: d.OnOk,
				Focus:  d.FocusOk,
			},
		},
		Log: d.Log,
	}.CreateBox(cfg, onDone)
}
