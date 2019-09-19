package prompt

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"strings"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	Colors               ColorsConfig
	PrintDivider         bool
	PrintCommand         bool
	HistoryFile          string
}

type ColorsConfig struct {
	Bg      tcell.Color
	Label   tcell.Color
	Text    tcell.Color
	Divider tcell.Color
	Command tcell.Color
}

func NewModule(cfg Config) *Module {
	return &Module{
		cfg:     cfg,
		history: newHistory(cfg.HistoryFile),
	}
}

type Module struct {
	cfg     Config
	view    *tview.InputField
	cmd     *CmdRunner
	history *history
	*gooster.AppContext
}

func (w *Module) Name() string {
	return "prompt"
}

func (w *Module) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
	w.AppContext = ctx
	w.cmd = &CmdRunner{
		ctx:    ctx,
		Stdout: ctx.Output(),
		Stderr: ctx.Output(),
	}

	w.view = tview.NewInputField()
	w.view.SetLabel(" > ")
	w.view.SetBorder(false)

	w.view.SetLabelColor(w.cfg.Colors.Label)
	w.view.SetBackgroundColor(w.cfg.Colors.Bg)
	w.view.SetFieldBackgroundColor(w.cfg.Colors.Bg)
	w.view.SetFieldTextColor(w.cfg.Colors.Text)

	w.Actions().OnSetPrompt(func(input string) {
		w.view.SetText(input)
	})

	w.Actions().OnCommandInterrupt(func() {
		w.view.SetText("")
	})

	//w.view.SetAutocompleteFunc(func(currentText string) (entries []string) {
	//	return []string{"foo", "bar", "baz"}
	//})
	w.view.SetDoneFunc(w.processKeyPress)

	return w.view, w.cfg.ModuleConfig, nil
}

func (w *Module) processKeyPress(key tcell.Key) {
	input := w.view.GetText()
	if input == "" {
		return
	}

	switch key {
	case tcell.KeyEnter:
		w.executeCommand(input)
	}
}

func (w *Module) executeCommand(input string) {
	if w.cfg.PrintDivider {
		_, _, width, _ := w.view.GetInnerRect()
		div := strings.Repeat("-", width-2)
		w.Actions().Write(fmt.Sprintf("[%s]%s[-]\n", w.getColorName(w.cfg.Colors.Divider), div))
	}

	if w.cfg.PrintCommand {
		w.Actions().Write(fmt.Sprintf("[%s]> %s[-]\n", w.getColorName(w.cfg.Colors.Command), input))
	}

	w.view.SetText("")
	err := w.cmd.Run(Command{
		Cmd:   input,
		Async: true,
		Ctx:   nil,
	})
	if err != nil {
		w.Log().Error(err)
	}
}

func (w *Module) getColorName(c tcell.Color) string {
	for name, value := range tcell.ColorNames {
		if value == c {
			return name
		}
	}
	return "black"
}