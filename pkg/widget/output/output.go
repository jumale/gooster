package output

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type Config struct {
	gooster.WidgetConfig `json:",inline"`
	Colors               ColorsConfig
}

type ColorsConfig struct {
	Bg   tcell.Color
	Text tcell.Color
}

func NewWidget(cfg Config) *Widget {
	return &Widget{cfg: cfg}
}

type Widget struct {
	cfg  Config
	view *tview.TextView
	*gooster.AppContext
}

func (w *Widget) Name() string {
	return "output"
}

func (w *Widget) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.WidgetConfig, error) {
	w.AppContext = ctx

	w.view = tview.NewTextView()
	w.view.SetBorder(false)
	w.view.SetDynamicColors(true)
	w.view.SetScrollable(true)
	w.view.SetBorderPadding(0, 0, 1, 1)

	w.view.SetBackgroundColor(w.cfg.Colors.Bg)
	w.view.SetTextColor(w.cfg.Colors.Text)

	w.Actions().OnOutput(func(data []byte) {
		if _, err := w.view.Write(data); err != nil {
			w.Log().Error(errors.WithMessage(err, "write to output"))
		}
	})

	return w.view, w.cfg.WidgetConfig, nil
}
