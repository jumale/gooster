package dialog

import "github.com/rivo/tview"

type ActionHandler func(form *tview.Form)

type Button struct {
	Label  string
	Action ActionHandler
	Focus  bool
}
