package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type AppConfig struct {
	Grid          GridConfig
	InitDir       string
	LogLevel      LogLevel
	EventsLogPath string
}

type GridConfig struct {
	Cols []int
	Rows []int
}

func NewApp(cfg AppConfig) (*App, error) {
	ctx, err := NewAppContext(cfg)
	if err != nil {
		return nil, err
	}

	ctx.Logger.Debug("Start initializing app")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(tcell.ColorDefault)
	grid.SetColumns(cfg.Grid.Cols...)
	grid.SetRows(cfg.Grid.Rows...)

	view := tview.NewApplication()
	view.SetRoot(grid, true)

	app := &App{
		cfg:  cfg,
		root: view,
		grid: grid,
		ctx:  ctx,
	}

	ctx.EventManager.Dispatch(Event{
		Id:   EventWorkDirChange,
		Data: cfg.InitDir,
	})

	ctx.Logger.Debug("App is initialized")

	return app, nil
}

type App struct {
	cfg     AppConfig
	root    *tview.Application
	grid    *tview.Grid
	widgets []Widget
	ctx     *AppContext
}

func (app *App) AddWidget(w Widget) {
	if err := w.Init(app.ctx); err != nil {
		panic(err)
	}

	app.widgets = append(app.widgets, w)

	app.grid.AddItem(
		w.View(),
		w.Config().Row,
		w.Config().Col,
		w.Config().Height,
		w.Config().Width,
		0,
		0,
		w.Config().Focused,
	)
	app.ctx.Logger.DebugF("Initializing widget [lightgreen]'%s'[-] with config [lightblue]%+v[-]\n", w.Name(), w.Config())
}

func (app *App) Run() {
	app.ctx.Logger.Debug("Starting App")
	defer func() {
		if err := app.Close(); err != nil {
			panic(err)
		}
	}()

	app.ctx.EventManager.Start()

	if err := app.root.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Close() error {
	app.ctx.Logger.Debug("Closing App")
	if err := app.ctx.EventManager.Close(); err != nil {
		return err
	}

	return nil
}
