package ext

import (
	"github.com/jumale/gooster/pkg/command"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
)

type Completions = []string

type BashCompletionConfig struct {
	gooster.ExtensionConfig
	Completer command.BashCompleterConfig
}

type BashCompletion struct {
	cfg       BashCompletionConfig
	completer *command.BashCompleter
}

func NewBashCompletion(cfg BashCompletionConfig) gooster.Extension {
	return &BashCompletion{
		cfg:       cfg,
		completer: command.NewBashCompleter(cfg.Completer),
	}
}

func (b *BashCompletion) Config() gooster.ExtensionConfig {
	return b.cfg.ExtensionConfig
}

func (b *BashCompletion) Init(_ gooster.Module, ctx *gooster.AppContext) error {
	ctx.Events().Subscribe(events.HandleWithPrio(10, func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case gooster.EventSetCompletion:
			// skip if command already has completions defined by someone else
			if len(event.Completion) > 0 {
				return e
			}

			if len(event.Commands) == 0 {
				return e
			}

			var err error
			if event.Completion, err = b.completer.Get(event.Commands[len(event.Commands)-1]); err != nil {
				ctx.Log().Debug(err)
			}
			return event
		}
		return e
	}))

	return nil
}
