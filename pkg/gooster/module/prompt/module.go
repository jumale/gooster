package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/command"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/history"
	"github.com/rivo/tview"
	"strings"
)

type Module struct {
	*gooster.BaseModule
	cfg         Config
	fs          filesys.FileSys
	view        *tview.InputField
	history     *history.Manager
	cmd         *Command
	latestInput string
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
		case gooster.EventSetCompletion:
			if len(event.Completion) == 1 {
				m.view.SetText(command.Complete(m.view.GetText(), event.Completion[0]))
				return events.StopPropagation
			}
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
		m.handleCompletion(input)

	case tcell.KeyEnter:
		if m.cmd == nil {
			m.Events().Dispatch(EventExecCommand{Cmd: input})
		} else {
			m.Events().Dispatch(EventSendUserInput{Input: input})
		}
	}
}

func (m *Module) clearCommand() {
	lineBreak := ""
	if m.cmd != nil && m.cmd.LastChar() != newLine {
		lineBreak = "\n"
	}

	m.cmd = nil
	if m.cfg.PrintDivider {
		_, _, width, _ := m.view.GetInnerRect()
		m.Output().WriteF(
			"%s[%s]%s[-]\n",
			lineBreak,
			getColorName(m.cfg.Colors.Divider),
			strings.Repeat("-", width-2),
		)
	}
}

func (m *Module) clearPrompt() {
	m.view.SetText("")
	m.history.Reset()
}
