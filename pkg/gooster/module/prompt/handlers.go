package prompt

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/convert"
	"github.com/jumale/gooster/pkg/events"
	"strings"
)

func (m *Module) handleEventSetPrompt(event events.Event) {
	m.view.SetText(convert.ToString(event.Payload))
}

func (m *Module) handleEventClearPrompt(event events.Event) {
	m.view.SetText("")
	m.history.Reset()
}

func (m *Module) handleEventExecCommand(event events.Event) {
	command := convert.ToString(event.Payload)
	cmdContext, cancelFunc := context.WithCancel(context.Background())
	m.cancelFunc = cancelFunc

	if m.cfg.PrintDivider {
		_, _, width, _ := m.view.GetInnerRect()
		div := strings.Repeat("-", width-2)
		m.Events().Dispatch(events.Event{
			Id:      "output:write", // @todo use Actions
			Payload: fmt.Sprintf("[%s]%s[-]\n", getColorName(m.cfg.Colors.Divider), div),
		})
	}

	if m.cfg.PrintCommand {
		m.Events().Dispatch(events.Event{
			Id:      "output:write", // @todo use Actions
			Payload: fmt.Sprintf("[%s]> %s[-]\n", getColorName(m.cfg.Colors.Command), command),
		})
	}

	m.history.Add(command)
	m.view.SetText("")
	err := m.runner.Run(Command{
		Cmd:   command,
		Async: true,
		Ctx:   cmdContext,
	})
	m.check(err)
}

func (m *Module) handleInterruptCommand(event events.Event) {
	if m.cancelFunc != nil {
		m.cancelFunc()
	}
}

func (m *Module) handleKeyHistoryPrev(event *tcell.EventKey) *tcell.EventKey {
	m.actions.SetPrompt(m.history.Prev())
	return event
}

func (m *Module) handleKeyHistoryNext(event *tcell.EventKey) *tcell.EventKey {
	m.actions.SetPrompt(m.history.Next())
	return event
}
