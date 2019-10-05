package gooster

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"os"
)

type AppConfig struct {
	Grid          GridConfig
	Keys          KeysConfig
	InitDir       string
	LogLevel      log.Level
	EventsLogPath string
	Debug         bool
	Dialog        dialog.Config
}

type GridConfig struct {
	Cols []int
	Rows []int
}

type KeysConfig struct {
	Exit tcell.Key
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

	pages := tview.NewPages()
	pages.SetBackgroundColor(tcell.ColorDefault)

	app := &App{
		cfg:      cfg,
		root:     root,
		pages:    pages,
		ctx:      ctx,
		focusMap: make(map[tcell.Key]tview.Primitive),
		modal:    newModalManger(cfg.Dialog, ctx, pages),
	}
	app.modal.onClose = func() {
		if app.lastFocus != nil {
			app.root.SetFocus(app.lastFocus)
		}
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
	cfg       AppConfig
	root      *tview.Application
	pages     *tview.Pages
	modules   []moduleDefinition
	ctx       *AppContext
	focusMap  map[tcell.Key]tview.Primitive
	lastFocus tview.Primitive
	modal     *modalManger
}

func (app *App) RegisterModule(mod Module) {
	view, cfg, err := mod.Init(app.ctx)
	if err != nil {
		panic(errors.WithMessage(err, "init module"))
	}

	if !cfg.Enabled {
		return
	}

	app.modules = append(app.modules, moduleDefinition{
		module: mod,
		view:   view,
		cfg:    cfg,
	})
	if cfg.FocusKey != 0 {
		app.focusMap[cfg.FocusKey] = view
	}

	app.ctx.log.InfoF("Initializing module [lightgreen]'%s'[-] with config [lightblue]%+v[-]", mod.Name(), cfg)
}

func (app *App) Run() {
	app.ctx.log.Info("Starting App")
	app.ctx.actions.registerActionOwners(app.modules)
	app.root.SetInputCapture(app.createInputHandler(
		app.handleFocusKeys,
		app.handleInterrupt,
		app.handleCloseDialog,
		app.handleExit,
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
		app.ctx.log.DebugF("App: focusing view: %T", view)
		app.root.SetFocus(view)
	})

	app.ctx.em.Start()

	app.newTab()
	app.root.SetRoot(app.pages, true)

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

func (app *App) newTab() {
	app.ctx.log.Debug("Creating new tab")
	grid := tview.NewGrid()
	grid.SetBackgroundColor(tcell.ColorDefault)
	grid.SetColumns(app.cfg.Grid.Cols...)
	grid.SetRows(app.cfg.Grid.Rows...)

	for _, def := range app.modules {
		grid.AddItem(
			def.view,
			def.cfg.Row, def.cfg.Col,
			def.cfg.Height, def.cfg.Width,
			0, 0,
			def.cfg.Focused,
		)
	}

	pageId := fmt.Sprintf("gooster_tab_%d", app.pages.GetPageCount()+1)

	app.pages.AddPage(pageId, grid, true, true)
}

func (app *App) createInputHandler(handlers ...InputHandler) InputCapture {
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

func (app *App) handleCloseDialog(event *tcell.EventKey) (newEvent *tcell.EventKey, handled bool) {
	if event.Key() == tcell.KeyEscape && app.modal.isOpen {
		app.ctx.log.Debug("Closing dialog")
		app.ctx.actions.CloseDialog()

		return event, true
	}

	return event, false
}

func (app *App) handleFocusKeys(event *tcell.EventKey) (newEvent *tcell.EventKey, handled bool) {
	if view, ok := app.focusMap[event.Key()]; ok {
		app.ctx.log.DebugF("App: focusing view by key %s", tcell.KeyNames[event.Key()])
		app.ctx.actions.SetFocus(view)
		app.lastFocus = view
		return event, true
	}

	return event, false
}

func (app *App) handleExit(event *tcell.EventKey) (newEvent *tcell.EventKey, handled bool) {
	if event.Key() == app.cfg.Keys.Exit {
		app.ctx.actions.Exit()
		return event, true
	}
	return event, false
}
