package status

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	"github.com/rivo/tview"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	Colors               ColorsConfig
}

type ColorsConfig struct {
	Bg         tcell.Color
	WorkDir    tcell.Color
	Branch     tcell.Color
	K8sContext tcell.Color
}

func NewModule(cfg Config) *Module {
	return &Module{cfg: cfg}
}

type Module struct {
	*gooster.BaseModule
	cfg Config
	wd  *tview.TableCell
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	view := tview.NewTable()
	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, view, view.Box)

	view.SetBorder(false)
	view.SetBorders(false)
	view.SetBackgroundColor(m.cfg.Colors.Bg)

	wd := tview.NewTableCell("")
	wd.SetTextColor(m.cfg.Colors.WorkDir)
	wd.SetExpansion(2)
	wd.SetAlign(tview.AlignLeft)
	view.SetCell(0, 0, wd)
	m.wd = wd

	branch := tview.NewTableCell("master")
	branch.SetTextColor(m.cfg.Colors.Branch)
	branch.SetExpansion(1)
	branch.SetAlign(tview.AlignCenter)
	view.SetCell(0, 1, branch)

	k8sCtx := tview.NewTableCell("some.long-context.preview.ams1.example.com")
	k8sCtx.SetTextColor(m.cfg.Colors.K8sContext)
	k8sCtx.SetExpansion(2)
	k8sCtx.SetAlign(tview.AlignRight)
	view.SetCell(0, 2, k8sCtx)

	m.Events().Subscribe(events.HandleFunc(func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case workdir.EventChangeDir:
			m.handleEventChangeDir(event)
		}
		return e
	}))

	return nil
}
