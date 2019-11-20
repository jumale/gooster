package status

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	Colors               ColorsConfig
}

type ColorsConfig struct {
	Bg tcell.Color
}

func NewModule(cfg Config) *Module {
	return &Module{cfg: cfg}
}

type Module struct {
	*gooster.BaseModule
	cfg  Config
	view *tview.TableCell
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	view := tview.NewTable()
	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, view, view.Box)

	view.SetBorder(false)
	view.SetBorders(false)
	view.SetBackgroundColor(m.cfg.Colors.Bg)

	m.Events().Subscribe(events.HandleWithPrio(events.AfterAllOtherChanges, func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case EventShowInStatus:
			cell := tview.NewTableCell(event.Value)
			cell.SetExpansion(2)
			cell.SetAlign(event.Align)
			view.SetCell(0, event.Col, cell)
		}
		return e
	}))
	return nil
}
