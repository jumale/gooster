package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"io"
)

type App struct {
	*AppContext
	cfg       AppConfig
	root      *tview.Application
	pages     *tview.Pages
	modules   []moduleDefinition
	focusMap  map[config.Key]tview.Primitive
	lastFocus tview.Primitive
}

func NewApp(cfgSource io.Reader, defaultCfgSource io.Reader) (*App, error) {
	fs := filesys.Default{}
	configReader := config.NewYamlReader(config.YamlReaderConfig{
		Fs:       fs,
		Defaults: defaultCfgSource,
	})
	if err := configReader.Load(cfgSource); err != nil {
		return nil, errors.WithMessage(err, "Failed to load gooster config")
	}

	appCfg := defaultConfig
	if err := configReader.Read("$.app", &appCfg); err != nil {
		return nil, errors.WithMessage(err, "Failed to read app config")
	}

	root := tview.NewApplication()

	ctx, err := NewAppContext(AppContextConfig{
		LogLevel:          appCfg.LogLevel,
		DelayEventManager: true,
		FileSys:           filesys.Default{},
		ConfigReader:      configReader,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "init app context")
	}

	ctx.log.Info("Start initializing app")

	pages := tview.NewPages()
	pages.SetBackgroundColor(tcell.ColorDefault)

	app := &App{
		AppContext: ctx,
		cfg:        appCfg,
		root:       root,
		pages:      pages,
		focusMap:   make(map[config.Key]tview.Primitive),
	}

	ctx.log.Info("App is initialized")

	return app, nil
}

func (app *App) RegisterModule(mod Module, extensions ...Extension) {
	app.modules = append(app.modules, moduleDefinition{
		module:     mod,
		extensions: extensions,
	})
}

func (app *App) Run() {
	// init event handlers
	app.Events().Subscribe(events.HandleWithPrio(events.AfterAllOtherChanges, func(e events.IEvent) events.IEvent {
		switch e.(type) {
		case EventExit:
			app.handleExitEvent()
		}
		return e
	}))
	app.Events().Subscribe(events.HandleFunc(func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case EventSetFocus:
			app.handleSetFocusEvent(event)
		case EventDraw:
			app.handleDrawEvent()
		case EventOpenDialog:
			app.handleEventOpenDialog(event)
		case EventCloseDialog:
			app.handleEventCloseDialog()
		case EventAddTab:
			app.handleEventAddTab(event)
		case EventShowTab:
			app.handleEventShowTab(event)
		case EventRemoveTab:
			app.handleEventRemoveTab(event)
		}
		return e
	}))

	app.Events().Dispatch(EventAddTab{Id: initialTabId, View: app.createMainGrid()})

	// init services and views
	if em, ok := app.Events().(DelayedEventManager); ok {
		if err := em.Init(); err != nil {
			panic(errors.WithMessage(err, "init event manager"))
		}
	}

	// init key handlers
	HandleKeyEvents(app.root, app.withFocusKeys(KeyEventHandlers{
		config.NewKey(tcell.KeyCtrlC):  app.handleKeyCtrlC,
		config.NewKey(tcell.KeyEscape): app.handleKeyEscape,
		app.cfg.Keys.Exit:              app.handleKeyExit,
	}))

	// debug keys
	prev := app.root.GetInputCapture()
	app.root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		app.Log().DebugF("Key press: %s [%d, %d, %d]", event.Name(), event.Key(), event.Rune(), event.Modifiers())
		return prev(event)
	})

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

	for _, def := range app.modules {
		mod, cfg, err := app.initModule(def.module, def.extensions...)
		if err != nil {
			panic(err)
		}

		focusKey := cfg.FocusKey
		if !focusKey.Empty() {
			app.focusMap[focusKey] = mod.View()
		}

		if mod != nil {
			grid.AddItem(
				mod.View(),
				cfg.Row, cfg.Col,
				cfg.Height, cfg.Width,
				0, 0,
				cfg.Focused,
			)
		}

	}
	return grid
}

func (app *App) initModule(mod Module, extensions ...Extension) (Module, *ModuleConfig, error) {
	modCtx := app.AppContext.forModule(mod)
	modCfg := defaultModConfig
	err := modCtx.LoadConfig(&modCfg)
	if err != nil {
		return nil, nil, errors.WithMessagef(err, "Failed to load config for module %T", mod)
	}

	err = mod.Init(modCtx)
	if err != nil {
		return nil, nil, errors.WithMessagef(err, "Failed to init module %T", mod)
	}
	for _, ext := range extensions {
		extCtx := app.AppContext.forExtension(ext, mod)

		extCfg := defaultExtConfig
		err = extCtx.LoadConfig(&extCfg)
		if err != nil {
			return nil, nil, errors.WithMessagef(err, "Failed to load config for extension %T of module %T", ext, mod)
		}

		if !extCfg.Enabled {
			continue
		}

		if err = ext.Init(mod, extCtx); err != nil {
			return nil, nil, errors.WithMessagef(err, "Failed to init extension %T of module %T", ext, mod)
		}
	}

	if !modCfg.FocusKey.Empty() {
		app.focusMap[modCfg.FocusKey] = mod.View()
	}

	app.Log().InfoF("Initialized module [lightgreen]'%T'[-]", mod)
	return mod, &modCfg, nil
}

// withFocusKeys adds handles for every focus key
func (app *App) withFocusKeys(keyHandlers KeyEventHandlers) KeyEventHandlers {
	for focusKey, view := range app.focusMap {
		v := view
		keyHandlers[focusKey] = func(event *tcell.EventKey) *tcell.EventKey {
			app.Events().Dispatch(EventSetFocus{Target: v})
			return nil
		}
	}
	return keyHandlers
}

type moduleDefinition struct {
	module     Module
	extensions []Extension
}
