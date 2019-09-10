package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"os"
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
		return nil, errors.WithMessage(err, "init app context")
	}
	ctx.Actions.afterAction = func(e events.Event) {
		root.Draw()
	}

	ctx.Log.Debug("Start initializing app")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(tcell.ColorDefault)
	grid.SetColumns(cfg.Grid.Cols...)
	grid.SetRows(cfg.Grid.Rows...)

	app := &App{
		cfg:      cfg,
		root:     root,
		grid:     grid,
		ctx:      ctx,
		focusMap: make(map[tcell.Key]tview.Primitive),
	}

	ctx.Actions.OnWorkDirChange(func(newPath string) {
		if err := os.Chdir(newPath); err != nil {
			ctx.Log.Error(errors.WithMessage(err, "change work dir"))
		}
	})
	if cfg.InitDir == "" {
		cfg.InitDir = getWd()
	}
	ctx.Actions.SetWorkDir(cfg.InitDir)
	ctx.Log.Debug("App is initialized")

	return app, nil
}

type App struct {
	cfg      AppConfig
	root     *tview.Application
	grid     *tview.Grid
	widgets  []Widget
	ctx      *AppContext
	focusMap map[tcell.Key]tview.Primitive
}

func (app *App) AddWidget(w Widget) {
	view, cfg, err := w.Init(app.ctx)
	if err != nil {
		panic(errors.WithMessage(err, "init widget"))
	}

	if !cfg.Enabled {
		return
	}

	app.widgets = append(app.widgets, w)
	if cfg.FocusKey != 0 {
		app.focusMap[cfg.FocusKey] = view
	}

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

	defer func() {
		if err := app.Close(); err != nil {
			panic(errors.WithMessage(err, "closing app"))
		}
	}()

	app.ctx.EventManager.Start()

	app.root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if view, ok := app.focusMap[event.Key()]; ok {
			app.root.SetFocus(view)
		}

		return event
	})

	app.root.SetRoot(app.grid, true)
	if err := app.root.Run(); err != nil {
		panic(errors.WithMessage(err, "run app"))
	}
}

func (app *App) Close() error {
	app.ctx.Log.Debug("Closing App")
	if err := app.ctx.EventManager.Close(); err != nil {
		return errors.WithMessage(err, "closing event manager")
	}

	return nil
}
