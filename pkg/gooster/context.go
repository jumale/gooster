package gooster

import (
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
	"io"
	"os"
)

type AppContext struct {
	cfg     AppConfig
	em      *events.DefaultManager
	log     log.Logger
	actions *actions
	output  io.Writer
}

func NewAppContext(cfg AppConfig) (ctx *AppContext, err error) {
	ctx = &AppContext{cfg: cfg}

	ctx.em, err = events.NewManager(events.ManagerConfig{
		SubscriberStackLevel: 4,
		LogFile:              cfg.EventsLogPath,
		DelayedStart:         true,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "init event manager")
	}

	ctx.actions = newActions(ctx.em)
	ctx.output = &outputWriter{actions: ctx.actions}
	ctx.log = log.NewSimpleLogger(cfg.LogLevel, ctx.output)

	return ctx, nil
}

func (ctx *AppContext) EventManager() events.Manager {
	return ctx.em
}

func (ctx *AppContext) Log() log.Logger {
	return ctx.log
}

func (ctx *AppContext) Actions() *actions {
	return ctx.actions
}

func (ctx *AppContext) Output() io.Writer {
	return ctx.output
}

func (ctx *AppContext) Close() error {
	ctx.log.Info("Closing context")

	// substitute logger with an stdout logger,
	// in case if some modules will try to send logs after everything is closed
	ctx.log = log.NewSimpleLogger(ctx.cfg.LogLevel, os.Stdout)

	err := ctx.em.Close()
	if err != nil {
		return errors.WithMessage(err, "closing event manager")
	}

	return nil
}

type outputWriter struct {
	actions *actions
}

func (o *outputWriter) Write(p []byte) (n int, err error) {
	o.actions.Write(p)
	return len(p), nil
}
