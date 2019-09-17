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
	Debug         bool
}

type GridConfig struct {
	Cols []int
	Rows []int
}

func NewApp(cfg AppConfig) (*App, error) {
	root := tview.NewApplication()
	if cfg.Debug {
		root.SetScreen(NewScreenStub(10, 10))
	}

	ctx, err := NewAppContext(cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "init app context")
	}
	ctx.actions.afterAction = func(e events.Event) {
		root.Draw()
	}

	ctx.log.Info("Start initializing app")

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

	ctx.actions.OnWorkDirChange(func(newPath string) {
		if err := os.Chdir(newPath); err != nil {
			ctx.log.Error(errors.WithMessage(err, "change work dir"))
		}
	})
	if cfg.InitDir == "" {
		cfg.InitDir = getWd()
	}
	ctx.actions.SetWorkDir(cfg.InitDir)
	ctx.log.Info("App is initialized")

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
	app.ctx.log.InfoF("Initializing widget [lightgreen]'%s'[-] with config [lightblue]%+v[-]", w.Name(), cfg)
}

func (app *App) Run() {
	app.ctx.log.Info("Starting App")
	app.root.SetInputCapture(app.createInputHandler(
		app.handleFocusKeys,
		app.handleInterrupt,
	))

	defer func() {
		if err := app.Close(); err != nil {
			panic(errors.WithMessage(err, "closing app"))
		}
	}()

	app.ctx.actions.OnAppExit(func() {
		app.ctx.log.Info("Stopping App")
		if err := app.Close(); err != nil {
			app.ctx.Log().Error(errors.WithMessage(err, "stopping app"))
		}
		app.root.Stop()
	})

	app.ctx.actions.OnSetFocus(func(view tview.Primitive) {
		app.ctx.log.Debug("App: focusing view")
		app.root.SetFocus(view)
	})

	app.ctx.em.Start()

	app.root.SetRoot(app.grid, true)
	if err := app.root.Run(); err != nil {
		panic(errors.WithMessage(err, "run app"))
	}
}

func (app *App) Close() error {
	app.ctx.log.Info("Closing App")
	if err := app.ctx.Close(); err != nil {
		return errors.WithMessage(err, "closing app context")
	}

	return nil
}

func (app *App) createInputHandler(
	handlers ...func(event *tcell.EventKey) (newEvent *tcell.EventKey, handled bool),
) func(event *tcell.EventKey) *tcell.EventKey {

	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyRune {
			app.ctx.log.DebugF("App: pressed key %s", tcell.KeyNames[event.Key()])
		}

		for _, handler := range handlers {
			if newEvent, handled := handler(event); handled {
				return newEvent
			}
		}
		return event
	}
}

func (app *App) handleInterrupt(event *tcell.EventKey) (newEvent *tcell.EventKey, handled bool) {
	if event.Key() == tcell.KeyCtrlC {
		app.ctx.log.Debug("Interrupting latest command")
		app.ctx.actions.InterruptLatestCommand()
		return &tcell.EventKey{}, true
	}

	return event, false
}

func (app *App) handleFocusKeys(event *tcell.EventKey) (newEvent *tcell.EventKey, handled bool) {
	if view, ok := app.focusMap[event.Key()]; ok {
		app.ctx.log.Debug("App: focusing view")
		app.ctx.actions.SetFocus(view)
		return event, true
	}

	return event, false
}
