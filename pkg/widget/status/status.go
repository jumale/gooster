package status

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

type Config struct {
	gooster.WidgetConfig `json:",inline"`
}

func NewWidget(cfg Config) *Widget {
	return &Widget{cfg: cfg}
}

type Widget struct {
	cfg  Config
	view *tview.Table
}

func (w *Widget) Name() string {
	return "Status"
}

func (w *Widget) Init(ctx *gooster.AppContext) error {
	w.view = tview.NewTable()
	w.view.SetBorder(false)
	w.view.SetBorders(false)
	w.view.SetBackgroundColor(tcell.ColorGray)

	wd := tview.NewTableCell("/Some/path")
	wd.SetTextColor(tcell.ColorYellow)
	wd.SetExpansion(2)
	wd.SetAlign(tview.AlignLeft)
	w.view.SetCell(0, 0, wd)

	branch := tview.NewTableCell("master")
	branch.SetTextColor(tcell.ColorLightGreen)
	branch.SetExpansion(1)
	branch.SetAlign(tview.AlignCenter)
	w.view.SetCell(0, 1, branch)

	kubeContext := tview.NewTableCell("some.long-context.preview.ams1.example.com")
	kubeContext.SetTextColor(tcell.ColorLightBlue)
	kubeContext.SetExpansion(2)
	kubeContext.SetAlign(tview.AlignRight)
	w.view.SetCell(0, 2, kubeContext)

	return nil
}

func (w *Widget) View() tview.Primitive {
	return w.view
}

func (w *Widget) Config() gooster.WidgetConfig {
	return w.cfg.WidgetConfig
}
