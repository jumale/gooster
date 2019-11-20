package gooster

import (
	"fmt"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
	"io"
	"os"
)

type AppContextConfig struct {
	LogLevel          log.Level
	LogFormat         string
	LogTarget         io.Writer
	DelayEventManager bool
}

type AppContext struct {
	cfg AppContextConfig
	log log.Logger
	em  events.Manager
	out *output
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

	var logTarget io.Writer = output
	if cfg.LogTarget != nil {
		logTarget = cfg.LogTarget
	}

	logger = log.NewSimpleLogger(logTarget, log.SimpleLoggerConfig{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})

	ctx = &AppContext{
		cfg: cfg,
		em:  em,
		log: logger,
		out: output,
	}

	return ctx, nil
}

// ------------------------------------------------------------ //

func (ctx *AppContext) GetConfig(jsonPath string, target interface{}) error {
	return nil
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

func (ctx *AppContext) Close() error {
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

type DelayedEventManager interface {
	Init() error
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
