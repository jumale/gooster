package complete

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

const CompletionView = "completion"

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	Colors               ColorsConfig
	Keys                 KeysConfig
}

type ColorsConfig struct {
	Bg tcell.Color
}

type KeysConfig struct {
	NextItem tcell.Key
}

func NewModule(cfg Config) *Module {
	return &Module{
		cfg: cfg,
	}
}

type Module struct {
	*gooster.BaseModule
	cfg  Config
	view *tview.Table
}

func (m *Module) Config() gooster.ModuleConfig {
	return m.cfg.ModuleConfig
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	m.view = tview.NewTable()
	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, m.view, m.view.Box)

	m.view.SetBackgroundColor(m.cfg.Colors.Bg)
	m.view.SetSelectable(true, true)

	m.Events().Subscribe(events.HandleFunc(func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case gooster.EventSetCompletion:
			m.handleSetCompletion(event)
		case gooster.EventSetFocusByName:
			if event.TargetName == CompletionView {
				m.Events().Dispatch(gooster.EventSetFocus{Target: m.view})
			}
		}
		return e
	}))

	//m.HandleKeyEvents(gooster.KeyEventHandlers{
	//	tcell.KeyUp: m.handleMoveUp,
	//	tcell.KeyDown: m.handleMoveDown,
	//	tcell.KeyLeft: m.handleMoveLeft,
	//	tcell.KeyRight: m.handleMoveRight,
	//})

	return nil
}
