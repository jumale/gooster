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
	*gooster.AppContext
}

func (w *Widget) Name() string {
	return "Status"
}

func (w *Widget) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.WidgetConfig, error) {
	w.AppContext = ctx

	w.view = tview.NewTable()
	w.view.SetBorder(false)
	w.view.SetBorders(false)
	w.view.SetBackgroundColor(tcell.ColorGray)

	wd := tview.NewTableCell("")
	wd.SetTextColor(tcell.ColorYellow)
	wd.SetExpansion(2)
	wd.SetAlign(tview.AlignLeft)
	w.view.SetCell(0, 0, wd)
	w.Actions.OnWorkDirChange(func(newPath string) {
		wd.SetText(newPath)
	})

	branch := tview.NewTableCell("master")
	branch.SetTextColor(tcell.ColorLightGreen)
	branch.SetExpansion(1)
	branch.SetAlign(tview.AlignCenter)
	w.view.SetCell(0, 1, branch)

	kubeCtx := tview.NewTableCell("some.long-context.preview.ams1.example.com")
	kubeCtx.SetTextColor(tcell.ColorLightBlue)
	kubeCtx.SetExpansion(2)
	kubeCtx.SetAlign(tview.AlignRight)
	w.view.SetCell(0, 2, kubeCtx)

	return w.view, w.cfg.WidgetConfig, nil
}
