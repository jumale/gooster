package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/helper"
	"github.com/jumale/gooster/pkg/history"
	"github.com/rivo/tview"
)

type Module struct {
	*gooster.BaseModule
	cfg     Config
	fs      filesys.FileSys
	view    *tview.InputField
	history *history.Manager
	actions Actions
	helper  helper.Actions
	cmd     *Command
}

func NewModule(cfg Config) *Module {
	return newModule(cfg, filesys.Default{})
}

func newModule(cfg Config, fs filesys.FileSys) *Module {
	return &Module{cfg: cfg, fs: fs}
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	m.view = tview.NewInputField()
	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, m.view, m.view.Box)
	m.actions = Actions{ctx}
	m.helper = helper.Actions{AppContext: ctx}

	m.history = history.NewManager(history.Config{
		HistoryFile: m.cfg.HistoryFile,
		Log:         ctx.Log(),
	})

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
		events.Subscriber{Id: ActionSendUserInput, Fn: m.handleSendUserInput},
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
	case tcell.KeyTab:
		m.helper.SetCompletion(input, nil)

	case tcell.KeyEnter:
		if m.cmd == nil {
			m.actions.ExecCommand(input)
		} else {
			m.actions.SendUserInput(input)
		}
	}
}
