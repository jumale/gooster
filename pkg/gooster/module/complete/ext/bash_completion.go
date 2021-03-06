package ext

import (
	"github.com/jumale/gooster/pkg/completion"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
)

type Completions = []string

type BashCompletionConfig struct {
	Completer completion.BashCompleterConfig
}

type BashCompletion struct {
	cfg       BashCompletionConfig
	completer *completion.BashCompleter
}

func NewBashCompletion() gooster.Extension {
	return &BashCompletion{
		cfg: BashCompletionConfig{
			Completer: completion.BashCompleterConfig{
				CompleteBin: "complete",
				CompgenBin:  "compgen",
			},
		},
	}
}

func (ext *BashCompletion) Name() string {
	return "bash_completion"
}

func (ext *BashCompletion) Init(_ gooster.Module, ctx gooster.Context) error {
	if err := ctx.LoadConfig(&ext.cfg); err != nil {
		return err
	}
	ext.completer = completion.NewBashCompleter(ext.cfg.Completer)

	ctx.Events().Subscribe(events.HandleWithPrio(10, func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case gooster.EventSetCompletion:
			// skip if command already has completions defined by someone else
			if len(event.Completion.Suggested) > 0 {
				return e
			}

			if len(event.Commands) == 0 {
				return e
			}

			var err error
			if event.Completion, err = ext.completer.Get(event.Commands[len(event.Commands)-1]); err != nil {
				ctx.Log().Debug(err)
			}
			return event
		}
		return e
	}))

	return nil
}
