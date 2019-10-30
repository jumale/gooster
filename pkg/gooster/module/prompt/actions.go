package prompt

import (
	"fmt"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
)

const (
	ActionSetPrompt        events.EventId = "prompt:set"
	ActionClearPrompt                     = "prompt:clear"
	ActionExecCommand                     = "command:exec"
	ActionInterruptCommand                = "command:interrupt"
	ActionSendUserInput                   = "command:send_user_input"
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

func (a Actions) SendUserInput(input string) {
	a.Events().Dispatch(events.Event{Id: ActionSendUserInput, Payload: input})
}

func (a Actions) writeOutputF(format string, v ...interface{}) {
	a.Events().Dispatch(events.Event{Id: "output:write", Payload: fmt.Sprintf(format, v...)}) // @todo use Actions
}

func (a Actions) changeDir(path string) {
	a.Events().Dispatch(events.Event{Id: "workdir:change_dir", Payload: path}) // @todo use Actions
}

func (a Actions) outputWriter() *outputWriter {
	return &outputWriter{a.Events()}
}

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
