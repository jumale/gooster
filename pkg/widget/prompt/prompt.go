package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

type Config struct {
	gooster.WidgetConfig `json:",inline"`
	Colors               ColorsConfig
}

type ColorsConfig struct {
	Bg    tcell.Color
	Label tcell.Color
	Text  tcell.Color
}

func NewWidget(cfg Config) *Widget {
	return &Widget{cfg: cfg}
}

type Widget struct {
	cfg  Config
	view *tview.InputField
	cmd  *Command
	*gooster.AppContext
}

func (w *Widget) Name() string {
	return "prompt"
}

func (w *Widget) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.WidgetConfig, error) {
	w.AppContext = ctx
	w.cmd = &Command{
		ctx:    ctx,
		Stdout: ctx.Output,
		Stderr: ctx.Output,
	}

	w.view = tview.NewInputField()
	w.view.SetLabel("> ")
	w.view.SetBorder(false)

	w.view.SetLabelColor(w.cfg.Colors.Label)
	w.view.SetBackgroundColor(w.cfg.Colors.Bg)
	w.view.SetFieldBackgroundColor(w.cfg.Colors.Bg)
	w.view.SetFieldTextColor(w.cfg.Colors.Text)

	//w.view.SetAutocompleteFunc(func(currentText string) (entries []string) {
	//	return []string{"foo", "bar", "baz"}
	//})

	w.view.SetDoneFunc(func(key tcell.Key) {
		w.Log.Debug(tcell.KeyNames[key])

		switch key {
		case tcell.KeyEnter:
			cmd := w.view.GetText()
			w.view.SetText("")
			err := w.cmd.Run(cmd)
			if err != nil {
				w.Log.Error(err)
			}
		}
	})

	return w.view, w.cfg.WidgetConfig, nil
}
