package complete

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

type Config struct {
	Colors ColorsConfig `json:"colors"`
	Keys   KeysConfig   `json:"keys"`
}

type ColorsConfig struct {
	Bg config.Color `json:"bg"`
}

type KeysConfig struct {
	NextItem config.Key `json:"next_item"`
	Select   config.Key `json:"select"`
}

type Module struct {
	gooster.Context
	cfg     Config
	view    *tview.Table
	current gooster.EventSetCompletion
}

func NewModule() *Module {
	return &Module{cfg: Config{
		Colors: ColorsConfig{
			Bg: config.Color(tcell.NewHexColor(0x333333)),
		},
		Keys: KeysConfig{
			NextItem: config.NewKey(tcell.KeyTab),
			Select:   config.NewKey(tcell.KeyEnter),
		},
	}}
}

func (m *Module) Name() string {
	return "complete"
}

func (m *Module) View() gooster.ModuleView {
	return m.view
}

func (m *Module) Init(ctx gooster.Context) (err error) {
	m.Context = ctx
	if err = ctx.LoadConfig(&m.cfg); err != nil {
		return err
	}

	m.view = tview.NewTable()

	m.view.SetBackgroundColor(m.cfg.Colors.Bg.Origin())
	m.view.SetSelectable(true, true)

	m.Events().Subscribe(events.HandleFunc(func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case gooster.EventSetCompletion:
			m.handleSetCompletion(event)
		}
		return e
	}))

	gooster.HandleKeyEvents(m.view, gooster.KeyEventHandlers{
		m.cfg.Keys.NextItem: m.handleNextItem,
		m.cfg.Keys.Select:   m.handleSelectItem,
	})

	return nil
}
