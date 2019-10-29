package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type App struct {
	*AppContext
	cfg       AppConfig
	root      *tview.Application
	pages     *tview.Pages
	modules   []Module
	focusMap  map[tcell.Key]tview.Primitive
	lastFocus tview.Primitive
}

func NewApp(cfg AppConfig) (*App, error) {
	root := tview.NewApplication()
	if cfg.Debug {
		root.SetScreen(NewScreenStub(10, 10))
	}

	ctx, err := newAppContext(cfg, func() { root.Draw() })
	if err != nil {
		return nil, errors.WithMessage(err, "init app context")
	}

	ctx.log.Info("Start initializing app")

	pages := tview.NewPages()
	pages.SetBackgroundColor(tcell.ColorDefault)

	app := &App{
		AppContext: ctx,
		cfg:        cfg,
		root:       root,
		pages:      pages,
		focusMap:   make(map[tcell.Key]tview.Primitive),
	}

	ctx.log.Info("App is initialized")

	return app, nil
}

func (app *App) RegisterModule(mod Module, extensions ...Extension) {
	if err := mod.Init(app.AppContext); err != nil {
		panic(errors.WithMessage(err, "init module"))
	}

	for _, ext := range extensions {
		if err := ext.Init(mod, app.AppContext); err != nil {
			panic(errors.WithMessage(err, "init extension"))
		}
	}

	cfg := mod.Config()
	if !cfg.Enabled {
		return
	}

	app.modules = append(app.modules, mod)
	if cfg.FocusKey != 0 {
		app.focusMap[cfg.FocusKey] = mod
	}

	app.Log().InfoF("Initialized module [lightgreen]'%T'[-] with config [lightblue]%+v[-]", mod, cfg)
}

func (app *App) Run() {
	// init event handlers
	app.Events().Subscribe(
		events.Subscriber{Id: ActionExit, Fn: app.handleExitEvent, Order: -9999}, // as late as possible
		events.Subscriber{Id: ActionSetFocus, Fn: app.handleSetFocusEvent},
		events.Subscriber{Id: ActionDraw, Fn: app.handleDrawEvent},
		events.Subscriber{Id: ActionOpenDialog, Fn: app.handleEventOpenDialog},
		events.Subscriber{Id: ActionCloseDialog, Fn: app.handleEventCloseDialog},
		events.Subscriber{Id: ActionAddTab, Fn: app.handleEventAddTab},
		events.Subscriber{Id: ActionShowTab, Fn: app.handleEventShowTab},
		events.Subscriber{Id: ActionRemoveTab, Fn: app.handleEventRemoveTab},
	)

	// init key handlers
	handleKeyEvents(&appInputAdaptor{app.root}, app.withFocusKeys(KeyEventHandlers{
		tcell.KeyCtrlC:    app.handleKeyCtrlC,
		tcell.KeyEscape:   app.handleKeyEscape,
		app.cfg.Keys.Exit: app.handleKeyExit,
	}))

	// init services and views
	app.em.Start()
	app.AppActions().AddTab(Tab{Id: initialTabId, View: app.createMainGrid()})
	app.root.SetRoot(app.pages, true)

	// start the app
	app.Log().Info("Starting App")
	if err := app.root.Run(); err != nil {
		panic(errors.WithMessage(err, "run app"))
	}
}

func (app *App) createMainGrid() tview.Primitive {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(tcell.ColorDefault)
	grid.SetColumns(app.cfg.Grid.Cols...)
	grid.SetRows(app.cfg.Grid.Rows...)

	for _, mod := range app.modules {
		cfg := mod.Config()
		grid.AddItem(
			mod,
			cfg.Row, cfg.Col,
			cfg.Height, cfg.Width,
			0, 0,
			cfg.Focused,
		)
	}
	return grid
}

// withFocusKeys adds handles for every focus key
func (app *App) withFocusKeys(keyHandlers KeyEventHandlers) KeyEventHandlers {
	for focusKey, view := range app.focusMap {
		v := view
		keyHandlers[focusKey] = func(event *tcell.EventKey) *tcell.EventKey {
			app.AppActions().SetFocus(v)
			return nil
		}
	}
	return keyHandlers
}
