package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/convert"
	"github.com/jumale/gooster/pkg/events"
	"regexp"
)

func (m *Module) handleEventSetPrompt(event events.Event) {
	m.view.SetText(convert.ToString(event.Payload))
}

func (m *Module) handleEventClearPrompt(event events.Event) {
	m.view.SetText("")
	m.history.Reset()
}

var ignoreCommandErrors = regexList{
	regexp.MustCompile("^exit status"),
	regexp.MustCompile("^signal: killed"),
}

func (m *Module) handleEventExecCommand(event events.Event) {
	if m.cmd != nil {
		m.Log().ErrorF("Previous command '%s' is still running. Wait for it to finish or cancel it.", m.cmd.Command())
		return
	}

	command := convert.ToString(event.Payload)
	m.history.Add(command)

	if m.cfg.PrintCommand {
		m.actions.writeOutputF("[%s]> %s[-]\n", getColorName(m.cfg.Colors.Command), command)
	}
	// If it's exit command
	if command == "exit" {
		go m.AppActions().Exit()
		return
	}
	// If it looks like "cd" command:
	if path := detectWorkDirPath(m.fs, command); path != "" {
		m.actions.changeDir(path)
		return
	}

	m.view.SetText("")
	m.cmd = NewCommand(command).SetOutput(m.actions.outputWriter())
	go func() {
		m.Log().DebugF("Starting command `%s`", command)
		if err := m.cmd.Run(); err != nil {
			if !ignoreCommandErrors.MatchString(err.Error()) {
				m.Log().Error(err)
			}
		}
		m.Log().DebugF("Command finished", command)
		m.clearCommand()
	}()
}

const newLine byte = 10

func (m *Module) handleSendUserInput(event events.Event) {
	if m.cmd == nil {
		m.Log().Error("Could not send user input - there is no current command running")
		return
	}
	m.view.SetText("")
	input := convert.ToBytes(event.Payload)
	if _, err := m.cmd.Write(append(input, newLine)); err != nil {
		m.Log().Error(err)
	}
}

func (m *Module) handleInterruptCommand(event events.Event) {
	if m.cmd == nil {
		return
	}
	m.cmd.Cancel()
	m.clearCommand()
}

func (m *Module) handleKeyHistoryPrev(event *tcell.EventKey) *tcell.EventKey {
	m.actions.SetPrompt(m.history.Prev())
	return event
}

func (m *Module) handleKeyHistoryNext(event *tcell.EventKey) *tcell.EventKey {
	m.actions.SetPrompt(m.history.Next())
	return event
}
