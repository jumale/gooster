package completion

import "github.com/jumale/gooster/pkg/command"

type Completer interface {
	Get(cmd command.Definition) (Completion, error)
}
