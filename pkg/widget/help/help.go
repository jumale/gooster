package help

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

type Config struct {
	gooster.WidgetConfig `json:",inline"`
}

func NewWidget(cfg Config) *Widget {
	return &Widget{cfg: cfg}
}

type Widget struct {
	cfg  Config
	view *tview.Grid
	*gooster.AppContext
}

func (w *Widget) Name() string {
	return "help"
}

func (w *Widget) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.WidgetConfig, error) {
	w.AppContext = ctx

	w.view = tview.NewGrid()
	w.view.SetBorder(false)
	w.view.SetBorders(false)
	w.view.SetBackgroundColor(tcell.ColorDefault)

	w.view.SetColumns(-1)
	w.view.SetRows(-1, -1)

	_ = w.addWidget(ctx, NewColorNamesWidget, gooster.Position{
		Col: 0, Row: 0,
		Width: 1, Height: 1,
	})
	_ = w.addWidget(ctx, NewKeyNamesWidget, gooster.Position{
		Col: 0, Row: 1,
		Width: 1, Height: 1,
	})

	return w.view, w.cfg.WidgetConfig, nil
}

func (w *Widget) addWidget(
	ctx *gooster.AppContext,
	constructor func(gooster.WidgetConfig) gooster.Widget,
	pos gooster.Position,
) error {
	widget := constructor(gooster.WidgetConfig{
		Position: pos,
		Enabled:  true,
		Focused:  false,
	})
	view, _, err := widget.Init(ctx)
	if err != nil {
		return err
	}
	w.view.AddItem(view, pos.Row, pos.Col, pos.Height, pos.Width, 0, 0, false)
	return nil
}
