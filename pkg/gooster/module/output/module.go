package output

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/ansi"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type Config struct {
	Colors ColorsConfig
}

type ColorsConfig struct {
	Bg   config.Color
	Text config.Color
}

type Module struct {
	gooster.Context
	cfg  Config
	view *tview.TextView
}

func NewModule() *Module {
	return &Module{cfg: Config{
		Colors: ColorsConfig{
			Bg:   config.Color(tcell.NewHexColor(0x222222)),
			Text: config.Color(tcell.ColorDefault),
		},
	}}
}

func (m *Module) Name() string {
	return "output"
}

func (m *Module) View() gooster.ModuleView {
	return m.view
}

func (m *Module) Init(ctx gooster.Context) error {
	m.Context = ctx
	if err := ctx.LoadConfig(&m.cfg); err != nil {
		return err
	}

	m.view = tview.NewTextView()
	output := ansi.NewWriter(m.view, ansi.WriterConfig{
		DefaultFg: m.cfg.Colors.Text.Origin(),
		DefaultBg: m.cfg.Colors.Bg.Origin(),
	})

	m.view.SetBorder(false)
	m.view.SetDynamicColors(true)
	m.view.SetScrollable(true)
	m.view.SetBorderPadding(0, 0, 1, 1)
	m.view.SetBackgroundColor(m.cfg.Colors.Bg.Origin())
	m.view.SetTextColor(m.cfg.Colors.Text.Origin())

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
