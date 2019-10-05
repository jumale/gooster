package output

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	Colors               ColorsConfig
}

type ColorsConfig struct {
	Bg   tcell.Color
	Text tcell.Color
}

func NewModule(cfg Config) *Module {
	return &Module{
		cfg:   cfg,
		write: func([]byte) {},
	}
}

type Module struct {
	cfg   Config
	view  *tview.TextView
	write func([]byte)
	*gooster.AppContext
}

func (w *Module) Name() string {
	return "output"
}

func (w *Module) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
	w.AppContext = ctx

	w.view = tview.NewTextView()
	w.view.SetBorder(false)
	w.view.SetDynamicColors(true)
	w.view.SetScrollable(true)
	w.view.SetBorderPadding(0, 0, 1, 1)

	w.view.SetBackgroundColor(w.cfg.Colors.Bg)
	w.view.SetTextColor(w.cfg.Colors.Text)

	return w.view, w.cfg.ModuleConfig, nil
}

func (w *Module) OnOutputWrite(content []byte) {
	w.write(content)
}

func (w *Module) OutputWriteCallback(callback func(interface{})) {
	w.write = func(bytes []byte) {
		if _, err := w.view.Write(bytes); err != nil {
			w.Log().Error(errors.WithMessage(err, "write to output"))
		}
		callback(bytes)
	}
}
