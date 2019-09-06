package gooster

import (
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/log"
	"io"
)

type AppContext struct {
	EventManager *events.Manager
	Log          log.Logger
	Actions      *actions
	Output       io.Writer
}

func NewAppContext(cfg AppConfig) (ctx *AppContext, err error) {
	ctx = &AppContext{}

	ctx.EventManager, err = events.NewManager(events.ManagerConfig{
		SubscriberStackLevel: 3,
		LogFile:              cfg.EventsLogPath,
	})
	if err != nil {
		return nil, err
	}

	ctx.Actions = &actions{
		em:          ctx.EventManager,
		afterAction: func(e events.Event) {},
	}
	ctx.Output = &outputWriter{actions: ctx.Actions}
	ctx.Log = log.NewSimpleLogger(cfg.LogLevel, ctx.Output)

	return ctx, nil
}

type outputWriter struct {
	actions *actions
}

func (o *outputWriter) Write(p []byte) (n int, err error) {
	o.actions.Write(p)
	return 0, nil
}
