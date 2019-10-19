package prompt

import (
	"context"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/history"
	"github.com/rivo/tview"
)

type Module struct {
	*gooster.BaseModule
	cfg        Config
	view       *tview.InputField
	runner     *CmdRunner
	history    *history.Manager
	actions    Actions
	cancelFunc context.CancelFunc
}

func NewModule(cfg Config) *Module {
	return &Module{
		cfg: cfg,
	}
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	m.view = tview.NewInputField()
	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, m.view, m.view.Box)
	m.actions = Actions{ctx}

	m.history = history.NewManager(history.Config{
		HistoryFile: m.cfg.HistoryFile,
		Log:         ctx.Log(),
	})
	m.runner = &CmdRunner{
		AppContext: ctx,
		Stdout:     ctx.Output(),
		Stderr:     ctx.Output(),
	}

	m.view.SetLabel(" > ")
	m.view.SetBorder(false)
	m.view.SetLabelColor(m.cfg.Colors.Label)
	m.view.SetBackgroundColor(m.cfg.Colors.Bg)
	m.view.SetFieldBackgroundColor(m.cfg.Colors.Bg)
	m.view.SetFieldTextColor(m.cfg.Colors.Text)

	m.Events().Subscribe(
		events.Subscriber{Id: ActionSetPrompt, Fn: m.handleEventSetPrompt},
		events.Subscriber{Id: ActionClearPrompt, Fn: m.handleEventClearPrompt},
		events.Subscriber{Id: ActionExecCommand, Fn: m.handleEventExecCommand},
		events.Subscriber{Id: ActionInterruptCommand, Fn: m.handleInterruptCommand},
	)
	m.HandleKeyEvents(gooster.KeyEventHandlers{
		m.cfg.Keys.HistoryPrev: m.handleKeyHistoryPrev,
		m.cfg.Keys.HistoryNext: m.handleKeyHistoryNext,
	})

	m.view.SetDoneFunc(m.submit)
	return nil
}

func (m *Module) submit(key tcell.Key) {
	input := m.view.GetText()
	if input == "" {
		return
	}
	switch key {
	case tcell.KeyEnter:
		m.actions.ExecCommand(input)
	}
}
