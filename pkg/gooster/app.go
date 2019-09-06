package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/log"
	"github.com/rivo/tview"
	"time"
)

type AppConfig struct {
	Grid          GridConfig
	InitDir       string
	LogLevel      log.Level
	EventsLogPath string
}

type GridConfig struct {
	Cols []int
	Rows []int
}

func NewApp(cfg AppConfig) (*App, error) {
	root := tview.NewApplication()
	ctx, err := NewAppContext(cfg)
	if err != nil {
		return nil, err
	}
	ctx.Actions.afterAction = func(e events.Event) {
		root.Draw()
	}

	ctx.Log.Debug("Start initializing app")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(tcell.ColorDefault)
	grid.SetColumns(cfg.Grid.Cols...)
	grid.SetRows(cfg.Grid.Rows...)
	root.SetRoot(grid, true)

	app := &App{
		cfg:  cfg,
		root: root,
		grid: grid,
		ctx:  ctx,
	}

	ctx.Actions.SetWorkDir(cfg.InitDir)
	ctx.Log.Debug("App is initialized")

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
	view, cfg, err := w.Init(app.ctx)
	if err != nil {
		panic(err)
	}

	app.widgets = append(app.widgets, w)

	app.grid.AddItem(
		view,
		cfg.Row, cfg.Col,
		cfg.Height, cfg.Width,
		0, 0,
		cfg.Focused,
	)
	app.ctx.Log.DebugF("Initializing widget [lightgreen]'%s'[-] with config [lightblue]%+v[-]", w.Name(), cfg)
}

func (app *App) Run() {
	app.ctx.Log.Debug("Starting App")

	go func() {
		time.Sleep(3 * time.Second)
		app.ctx.Log.Debug("---->")
	}()

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
	app.ctx.Log.Debug("Closing App")
	if err := app.ctx.EventManager.Close(); err != nil {
		return err
	}

	return nil
}
