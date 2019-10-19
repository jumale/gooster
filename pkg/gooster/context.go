package gooster

import (
	"fmt"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

type AppContext struct {
	cfg AppConfig
	log log.Logger
	act Actions
	em  *events.DefaultManager
}

func newAppContext(cfg AppConfig, drawFunc func()) (ctx *AppContext, err error) {
	hook := &eventHook{draw: drawFunc}
	em, err := events.NewManager(events.ManagerConfig{
		DelayedStart: true,
		BeforeEvent:  hook.beforeEvent,
		AfterEvent:   hook.afterEvent,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "init event manager")
	}

	ctx = &AppContext{
		cfg: cfg,
		em:  em,
		log: log.NewSimpleLogger(cfg.LogLevel, &outputWriter{em: em}),
	}
	hook.log = ctx.log
	ctx.act = Actions{ctx}

	return ctx, nil
}

// ------------------------------------------------------------ //

type eventHook struct {
	draw func()
	log  log.Logger
}

func (eh *eventHook) beforeEvent(event events.Event) bool {
	if string(event.Id) != "output:write" {
		data := ""
		switch event.Payload.(type) {
		case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			data = fmt.Sprintf(" %v", event.Payload)
		default:
			dataType := fmt.Sprintf("%T", event.Payload)
			if strings.Contains(dataType, ".Payload") {
				data = fmt.Sprintf(" %+v", event.Payload)
			} else {
				data = " " + dataType
			}
		}
		if event.Payload == nil {
			data = ""
		}
		msg := fmt.Sprintf("[gold]Event:[-] [lightseagreen]%s[-]%s", event.Id, data)
		if len(msg) > 130 {
			msg = msg[:130] + "..."
		}
		eh.log.Debug(msg)
	}
	return true
}

func (eh *eventHook) afterEvent(event events.Event) {
	eh.draw()
}

// ------------------------------------------------------------ //

func (ctx *AppContext) Config() AppConfig {
	return ctx.cfg
}

func (ctx *AppContext) Log() log.Logger {
	return ctx.log
}

func (ctx *AppContext) Events() events.Manager {
	return ctx.em
}

func (ctx *AppContext) AppActions() Actions {
	return ctx.act
}

func (ctx *AppContext) Output() io.Writer {
	return &outputWriter{em: ctx.em}
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

// ------------------------------------------------------------ //

type outputWriter struct {
	em events.Manager
}

func (o *outputWriter) Write(p []byte) (n int, err error) {
	o.em.Dispatch(events.Event{
		Id:      "output:write", // @todo use Actions
		Payload: p,
	})
	return len(p), nil
}
