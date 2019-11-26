package gooster

import (
	"fmt"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

type Context interface {
	LoadConfig(target interface{}) error
	Log() log.Logger
	Events() events.Manager
	Output() *output
	Fs() filesys.FileSys
}

type AppContextConfig struct {
	LogLevel          log.Level
	LogFormat         string
	LogTarget         io.Writer
	DelayEventManager bool
	ConfigReader      config.Reader
	FileSys           filesys.FileSys
}

type AppContext struct {
	cfg     AppContextConfig
	log     log.Logger
	em      events.Manager
	out     *output
	cfgPath string
}

func NewAppContext(cfg AppContextConfig) (ctx *AppContext, err error) {
	var em events.Manager
	var logger log.Logger

	em, err = events.NewManager(events.ManagerConfig{DelayedStart: cfg.DelayEventManager})
	if err != nil {
		return nil, errors.WithMessage(err, "init event manager")
	}

	em.Subscribe(events.HandleWithPrio(events.AfterAllOtherChanges, func(event events.IEvent) events.IEvent {
		logEventToOutput(logger, event)

		if drawable, ok := event.(DrawableEvent); ok {
			if drawable.NeedsDraw() {
				em.Dispatch(EventDraw{})
			}
		}

		return event
	}))

	output := &output{em}

	if cfg.LogTarget == nil {
		cfg.LogTarget = output
	}
	logger = log.NewSimpleLogger(cfg.LogTarget, log.SimpleLoggerConfig{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})

	if cfg.FileSys == nil {
		cfg.FileSys = filesys.Default{}
	}

	ctx = &AppContext{
		cfgPath: "$",
		cfg:     cfg,
		em:      em,
		log:     logger,
		out:     output,
	}

	return ctx, nil
}

// ------------------------------------------------------------ //

func (ctx *AppContext) LoadConfig(target interface{}) error {
	return ctx.cfg.ConfigReader.Read(ctx.cfgPath, target)
}

func (ctx *AppContext) Log() log.Logger {
	return ctx.log
}

func (ctx *AppContext) Events() events.Manager {
	return ctx.em
}

func (ctx *AppContext) Output() *output {
	return ctx.out
}

func (ctx *AppContext) Fs() filesys.FileSys {
	return ctx.cfg.FileSys
}

// ------------------------------------------------------------ //

func (ctx *AppContext) SetCfgPath(path string) {
	ctx.cfgPath = path
}

const (
	moduleCfgPath = "modules[?(@.#id == '%s')][0]"
	extCfgPath    = "extensions[?(@.#id == '%s')][0]"
)

func (ctx *AppContext) forModule(mod Module) *AppContext {
	newCtx := *ctx
	newCtx.cfgPath = joinPath("$", fmt.Sprintf(moduleCfgPath, mod.Name()))
	return &newCtx
}

func (ctx *AppContext) forExtension(ext Extension, target Module) *AppContext {
	newCtx := *ctx
	newCtx.cfgPath = joinPath("$", fmt.Sprintf(moduleCfgPath, target.Name()), fmt.Sprintf(extCfgPath, ext.Name()))
	return &newCtx
}

func joinPath(vals ...string) string {
	return strings.Join(vals, ".")
}

// ------------------------------------------------------------ //

func (ctx *AppContext) close() error {
	ctx.log.Info("Closing context")

	// substitute logger with an stdout logger,
	// in case if some modules will try to send logs after everything is closed
	ctx.log = log.NewSimpleLogger(os.Stdout, log.SimpleLoggerConfig{
		Level: ctx.cfg.LogLevel,
	})

	if err := ctx.closeService(ctx.em); err != nil {
		return errors.WithMessage(err, "closing event manager")
	}

	return nil
}

func (ctx *AppContext) closeService(s interface{}) error {
	if closer, ok := s.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------ //

func logEventToOutput(logger log.Logger, event events.IEvent) {
	switch event.(type) {
	case EventOutput:
		return
	case EventDraw:
		return
	default:
		msg := fmt.Sprintf("[gold]Event:[-] [lightseagreen]%T[-]%+v", event, event)
		if len(msg) > 130 {
			msg = msg[:130] + "..."
		}
		logger.Debug(msg)
	}
}
