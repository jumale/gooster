package output

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/ansi"
	"github.com/jumale/gooster/pkg/events"
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
	return &Module{cfg: cfg}
}

type Module struct {
	*gooster.BaseModule
	cfg Config
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	view := tview.NewTextView()
	output := ansi.NewWriter(view, ansi.WriterConfig{
		DefaultFg: m.cfg.Colors.Text,
		DefaultBg: m.cfg.Colors.Bg,
	})

	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, view, view.Box)

	view.SetBorder(false)
	view.SetDynamicColors(true)
	view.SetScrollable(true)
	view.SetBorderPadding(0, 0, 1, 1)
	view.SetBackgroundColor(m.cfg.Colors.Bg)
	view.SetTextColor(m.cfg.Colors.Text)

	m.Events().Subscribe(events.HandleFunc(func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case gooster.EventOutput:
			if _, err := output.Write(event.Data); err != nil {
				m.Log().Error(errors.WithMessage(err, "write to output"))
			}
		}
		return e
	}))

	return nil
}
