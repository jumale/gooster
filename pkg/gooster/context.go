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

type AppContextConfig struct {
	LogLevel          log.Level
	DelayEventManager bool
}

type AppContext struct {
	cfg AppContextConfig
	log log.Logger
	act Actions
	em  events.Manager
}

func NewAppContext(cfg AppContextConfig) (ctx *AppContext, err error) {
	var em events.Manager
	var logger log.Logger

	em, err = events.NewManager(events.ManagerConfig{
		DelayedStart: cfg.DelayEventManager,
		BeforeEvent: func(event events.Event) bool {
			logEventToOutput(logger, event)
			return true
		},
	})
	if err != nil {
		return nil, errors.WithMessage(err, "init event manager")
	}

	logger = log.NewSimpleLogger(cfg.LogLevel, &outputWriter{em: em})

	ctx = &AppContext{
		cfg: cfg,
		em:  em,
		log: logger,
	}
	ctx.act = Actions{ctx}

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

func (ctx *AppContext) AppActions() Actions {
	return ctx.act
}

func (ctx *AppContext) Close() error {
	ctx.log.Info("Closing context")

	// substitute logger with an stdout logger,
	// in case if some modules will try to send logs after everything is closed
	ctx.log = log.NewSimpleLogger(ctx.cfg.LogLevel, os.Stdout)

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

func logEventToOutput(logger log.Logger, event events.Event) bool {
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
		logger.Debug(msg)
	}
	return true
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
