package prompt

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"strings"
)

type Config struct {
	gooster.WidgetConfig `json:",inline"`
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

func NewWidget(cfg Config) *Widget {
	return &Widget{
		cfg:     cfg,
		history: newHistory(cfg.HistoryFile),
	}
}

type Widget struct {
	cfg     Config
	view    *tview.InputField
	cmd     *Command
	history *history
	*gooster.AppContext
}

func (w *Widget) Name() string {
	return "prompt"
}

func (w *Widget) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.WidgetConfig, error) {
	w.AppContext = ctx
	w.cmd = &Command{
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

	//w.view.SetAutocompleteFunc(func(currentText string) (entries []string) {
	//	return []string{"foo", "bar", "baz"}
	//})
	w.view.SetDoneFunc(w.processKeyPress)

	return w.view, w.cfg.WidgetConfig, nil
}

func (w *Widget) processKeyPress(key tcell.Key) {
	input := w.view.GetText()
	if input == "" {
		return
	}

	if w.cfg.PrintDivider {
		_, _, width, _ := w.view.GetInnerRect()
		div := strings.Repeat("-", width-2)
		w.Actions().Write(fmt.Sprintf("[%s]%s[-]\n", w.getColorName(w.cfg.Colors.Divider), div))
	}

	if w.cfg.PrintCommand {
		w.Actions().Write(fmt.Sprintf("[%s]> %s[-]\n", w.getColorName(w.cfg.Colors.Command), input))
	}

	switch key {
	case tcell.KeyEnter:
		w.view.SetText("")
		err := w.cmd.Run(input)
		if err != nil {
			w.Log().Error(err)
		}
	}
}

func (w *Widget) getColorName(c tcell.Color) string {
	for name, value := range tcell.ColorNames {
		if value == c {
			return name
		}
	}
	return "black"
}
