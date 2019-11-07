package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/cmd"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/history"
	"github.com/rivo/tview"
)

type Module struct {
	*gooster.BaseModule
	cfg     Config
	fs      filesys.FileSys
	view    *tview.InputField
	history *history.Manager
	cmd     *Command
}

func NewModule(cfg Config) *Module {
	return newModule(cfg, filesys.Default{})
}

func newModule(cfg Config, fs filesys.FileSys) *Module {
	if cfg.Label == "" {
		cfg.Label = " > "
	}
	return &Module{cfg: cfg, fs: fs}
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	m.view = tview.NewInputField()
	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, m.view, m.view.Box)

	m.history = history.NewManager(history.Config{
		HistoryFile: m.cfg.HistoryFile,
		Log:         ctx.Log(),
		FileSys:     m.fs,
	})

	m.view.SetLabel(m.cfg.Label)
	m.view.SetFieldWidth(m.cfg.FieldWidth)
	m.view.SetBorder(false)
	m.view.SetLabelColor(m.cfg.Colors.Label)
	m.view.SetBackgroundColor(m.cfg.Colors.Bg)
	m.view.SetFieldBackgroundColor(m.cfg.Colors.Bg)
	m.view.SetFieldTextColor(m.cfg.Colors.Text)

	m.Events().Subscribe(events.HandleFunc(func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case EventSetPrompt:
			m.handleEventSetPrompt(event)
		case EventClearPrompt:
			m.handleEventClearPrompt()
		case EventExecCommand:
			m.handleEventExecCommand(event)
		case EventSendUserInput:
			m.handleEventSendUserInput(event)
			m.Events().Dispatch(gooster.EventOutput{Data: []byte(event.Input + "\n")})
		case gooster.EventInterrupt:
			m.handleEventInterruptCommand()
		}
		return e
	}))

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
		// ignore parsing errors, because non-complete commands are also allowed
		commands, _ := cmd.ParseCommands(input)
		m.Events().Dispatch(gooster.EventSetCompletion{Commands: commands})

	case tcell.KeyEnter:
		if m.cmd == nil {
			m.Events().Dispatch(EventExecCommand{Cmd: input})
		} else {
			m.Events().Dispatch(EventSendUserInput{Input: input})
		}
	}
}
