package prompt

import (
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
)

const (
	ActionSetPrompt        events.EventId = "prompt:set"
	ActionClearPrompt                     = "prompt:clear"
	ActionExecCommand                     = "command:exec"
	ActionInterruptCommand                = "command:interrupt"
)

type Actions struct {
	*gooster.AppContext
}

func (a Actions) SetPrompt(input string) {
	a.Events().Dispatch(events.Event{Id: ActionSetPrompt, Payload: input})
}

func (a Actions) ClearPrompt() {
	a.Events().Dispatch(events.Event{Id: ActionClearPrompt})
}

func (a Actions) ExecCommand(cmd string) {
	a.Events().Dispatch(events.Event{Id: ActionExecCommand, Payload: cmd})
}

func (a Actions) InterruptCommand() {
	a.Events().Dispatch(events.Event{Id: ActionInterruptCommand})
}
