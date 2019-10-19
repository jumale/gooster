package output

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/convert"
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
	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, view, view.Box)

	view.SetBorder(false)
	view.SetDynamicColors(true)
	view.SetScrollable(true)
	view.SetBorderPadding(0, 0, 1, 1)
	view.SetBackgroundColor(m.cfg.Colors.Bg)
	view.SetTextColor(m.cfg.Colors.Text)

	m.Events().Subscribe(
		events.Subscriber{Id: ActionWrite, Fn: func(event events.Event) {
			if _, err := view.Write(convert.ToBytes(event.Payload)); err != nil {
				m.Log().Error(errors.WithMessage(err, "write to output"))
			}
		}},
	)
	return nil
}
