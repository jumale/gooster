package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/command"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/history"
	"github.com/rivo/tview"
	"strings"
)

type Module struct {
	gooster.Context
	cfg         Config
	view        *tview.InputField
	history     *history.Manager
	cmd         *Command
	latestInput string
}

func NewModule() *Module {
	return &Module{cfg: Config{
		Label:        " > ",
		PrintDivider: true,
		PrintCommand: true,
		HistoryFile:  "~/.bash_history",
		Colors: ColorsConfig{
			Bg:      config.Color(tcell.NewHexColor(0x555555)),
			Label:   config.Color(tcell.ColorLime),
			Text:    config.Color(tcell.ColorLightGray),
			Divider: config.Color(tcell.ColorLightGreen),
			Command: config.Color(tcell.ColorLightSkyBlue),
		},
		Keys: KeysConfig{
			HistoryNext: config.Key(tcell.KeyDown),
			HistoryPrev: config.Key(tcell.KeyUp),
		},
	}}
}

func (m *Module) Name() string {
	return "prompt"
}

func (m *Module) View() gooster.ModuleView {
	return m.view
}

func (m *Module) Init(ctx gooster.Context) (err error) {
	m.Context = ctx
	if err = ctx.LoadConfig(&m.cfg); err != nil {
		return err
	}

	m.history, err = history.NewManager(history.Config{
		HistoryFile: m.cfg.HistoryFile,
		Log:         ctx.Log(),
		FileSys:     ctx.Fs(),
	})
	if err != nil {
		return err
	}

	m.view = tview.NewInputField()
	m.view.SetLabel(m.cfg.Label)
	m.view.SetFieldWidth(m.cfg.FieldWidth)
	m.view.SetBorder(false)
	m.view.SetLabelColor(m.cfg.Colors.Label.Origin())
	m.view.SetBackgroundColor(m.cfg.Colors.Bg.Origin())
	m.view.SetFieldBackgroundColor(m.cfg.Colors.Bg.Origin())
	m.view.SetFieldTextColor(m.cfg.Colors.Text.Origin())

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
				m.view.SetText(command.ApplyCompletion(m.view.GetText(), event.Completion[0]+" "))
				return events.StopPropagation
			}
		}
		return e
	}))

	gooster.HandleKeyEvents(m.view, gooster.KeyEventHandlers{
		m.cfg.Keys.HistoryPrev.Origin(): m.handleKeyHistoryPrev,
		m.cfg.Keys.HistoryNext.Origin(): m.handleKeyHistoryNext,
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
			getColorName(m.cfg.Colors.Divider.Origin()),
			strings.Repeat("-", width-2),
		)
	}
}

func (m *Module) clearPrompt() {
	m.view.SetText("")
	m.history.Reset()
}
