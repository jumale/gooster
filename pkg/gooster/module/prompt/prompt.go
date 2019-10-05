package prompt

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/history"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"strings"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	Colors               ColorsConfig
	PrintDivider         bool
	PrintCommand         bool
	HistoryFile          string
	Keys                 KeysConfig
}

type ColorsConfig struct {
	Bg      tcell.Color
	Label   tcell.Color
	Text    tcell.Color
	Divider tcell.Color
	Command tcell.Color
}

type KeysConfig struct {
	HistoryNext tcell.Key
	HistoryPrev tcell.Key
}

func NewModule(cfg Config) *Module {
	return &Module{
		cfg:       cfg,
		interrupt: func() {},
		setPrompt: func(string) {},
	}
}

type Module struct {
	cfg       Config
	view      *tview.InputField
	cmd       *CmdRunner
	history   *history.Manager
	interrupt func()
	setPrompt func(input string)
	*gooster.AppContext
}

func (m *Module) Name() string {
	return "prompt"
}

func (m *Module) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
	m.AppContext = ctx
	m.history = history.NewManager(history.Config{
		HistoryFile: m.cfg.HistoryFile,
		Log:         ctx.Log(),
	})

	m.cmd = &CmdRunner{
		ctx:    ctx,
		Stdout: ctx.Output(),
		Stderr: ctx.Output(),
	}

	m.view = tview.NewInputField()
	m.view.SetLabel(" > ")
	m.view.SetBorder(false)

	m.view.SetLabelColor(m.cfg.Colors.Label)
	m.view.SetBackgroundColor(m.cfg.Colors.Bg)
	m.view.SetFieldBackgroundColor(m.cfg.Colors.Bg)
	m.view.SetFieldTextColor(m.cfg.Colors.Text)

	m.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case m.cfg.Keys.HistoryNext:
			m.view.SetText(m.history.Next())
			return &tcell.EventKey{}

		case m.cfg.Keys.HistoryPrev:
			m.view.SetText(m.history.Prev())
			return &tcell.EventKey{}
		}
		return event
	})

	//m.view.SetAutocompleteFunc(func(currentText string) (entries []string) {
	//	return []string{"foo", "bar", "baz"}
	//})
	m.view.SetDoneFunc(m.processKeyPress)

	return m.view, m.cfg.ModuleConfig, nil
}

func (m *Module) OnCommandInterrupt() {
	m.interrupt()
}

func (m *Module) CommandInterruptCallback(callback func()) {
	m.interrupt = func() {
		m.setPrompt("")
		m.history.Reset()
		callback()
	}
}

func (m *Module) OnPromptSet(prompt string) {
	m.setPrompt(prompt)
}

func (m *Module) PromptChangeCallback(callback func(string)) {
	m.setPrompt = func(input string) {
		m.view.SetText(input)
		callback(input)
	}
}

func (m *Module) processKeyPress(key tcell.Key) {
	input := m.view.GetText()
	if input == "" {
		return
	}

	switch key {
	case tcell.KeyEnter:
		m.executeCommand(input)
	}
}

func (m *Module) executeCommand(input string) {
	if m.cfg.PrintDivider {
		_, _, width, _ := m.view.GetInnerRect()
		div := strings.Repeat("-", width-2)
		m.Actions().Write(fmt.Sprintf("[%s]%s[-]\n", m.getColorName(m.cfg.Colors.Divider), div))
	}

	if m.cfg.PrintCommand {
		m.Actions().Write(fmt.Sprintf("[%s]> %s[-]\n", m.getColorName(m.cfg.Colors.Command), input))
	}

	m.history.Add(input)
	m.view.SetText("")
	err := m.cmd.Run(Command{
		Cmd:   input,
		Async: true,
		Ctx:   nil,
	})
	m.check(err)
}

func (m *Module) getColorName(c tcell.Color) string {
	for name, value := range tcell.ColorNames {
		if value == c {
			return name
		}
	}
	return "black"
}

func (m *Module) check(err error, msg ...string) {
	if err == nil {
		return
	}
	if len(msg) > 0 {
		m.Log().Error(errors.WithMessage(err, msg[0]))
	} else {
		m.Log().Error(err)
	}
}
