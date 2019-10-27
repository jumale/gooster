package helper

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	Colors               ColorsConfig
	Keys                 KeysConfig
}

type ColorsConfig struct {
	Bg tcell.Color
}

type KeysConfig struct {
}

func NewModule(cfg Config) *Module {
	return &Module{
		cfg: cfg,
	}
}

type Module struct {
	*gooster.BaseModule
	cfg      Config
	complete *tview.Table
	actions  *Actions
}

func (m *Module) Config() gooster.ModuleConfig {
	return m.cfg.ModuleConfig
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	view := tview.NewPages()
	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, view, view.Box)
	m.actions = &Actions{ctx}

	view.SetBackgroundColor(m.cfg.Colors.Bg)

	m.complete = tview.NewTable()
	view.AddPage("complete", m.complete, true, true)

	m.Events().Subscribe(
		events.Subscriber{Id: ActionSetCompletion, Fn: m.handleSetCompletion},
	)

	return nil
}
