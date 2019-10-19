package help

import (
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
}

func NewModule(cfg Config) *Module {
	return &Module{cfg: cfg}
}

type Module struct {
	cfg Config
	*tview.Grid
	*gooster.AppContext
}

func (w *Module) Name() string {
	return "help"
}

//func (w *Module) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
//	w.AppContext = ctx
//
//	w.Grid = tview.NewGrid()
//	w.Grid.SetBorder(false)
//	w.Grid.SetBorders(false)
//	w.Grid.SetBackgroundColor(tcell.ColorDefault)
//
//	w.Grid.SetColumns(-1)
//	w.Grid.SetRows(-1, -1)
//
//	_ = w.addModule(ctx, NewColorNamesModule, gooster.Position{
//		Col: 0, Row: 0,
//		Width: 1, Height: 1,
//	})
//	_ = w.addModule(ctx, NewKeyNamesModule, gooster.Position{
//		Col: 0, Row: 1,
//		Width: 1, Height: 1,
//	})
//
//	return w.view, w.cfg.ModuleConfig, nil
//}
//
//func (w *Module) addModule(
//	ctx *gooster.AppContext,
//	constructor func(gooster.ModuleConfig) gooster.Module,
//	pos gooster.Position,
//) error {
//	module := constructor(gooster.ModuleConfig{
//		Position: pos,
//		Enabled:  true,
//		Focused:  false,
//	})
//	view, _, err := module.Init(ctx)
//	if err != nil {
//		return err
//	}
//	w.view.AddItem(view, pos.Row, pos.Col, pos.Height, pos.Width, 0, 0, false)
//	return nil
//}
