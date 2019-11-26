package status

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

type Config struct {
	Colors ColorsConfig
}

type ColorsConfig struct {
	Bg config.Color
}

type Module struct {
	gooster.Context
	cfg  Config
	view *tview.Table
}

func NewModule() *Module {
	return &Module{cfg: Config{
		Colors: ColorsConfig{
			Bg: config.Color(tcell.NewHexColor(0x555555)),
		},
	}}
}

func (m *Module) Name() string {
	return "status"
}

func (m *Module) View() gooster.ModuleView {
	return m.view
}

func (m *Module) Init(ctx gooster.Context) error {
	m.Context = ctx
	if err := ctx.LoadConfig(&m.cfg); err != nil {
		return err
	}

	m.view = tview.NewTable()
	m.view.SetBorder(false)
	m.view.SetBorders(false)
	m.view.SetBackgroundColor(m.cfg.Colors.Bg.Origin())

	m.Events().Subscribe(events.HandleWithPrio(events.AfterAllOtherChanges, func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case EventShowInStatus:
			cell := tview.NewTableCell(event.Value)
			cell.SetExpansion(2)
			cell.SetAlign(event.Align)
			m.view.SetCell(0, event.Col, cell)
		}
		return e
	}))
	return nil
}
