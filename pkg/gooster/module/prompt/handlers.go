package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/command"
	"github.com/jumale/gooster/pkg/convert"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	"regexp"
)

func (m *Module) handleEventSetPrompt(event EventSetPrompt) {
	m.view.SetText(event.Input)
}

func (m *Module) handleEventClearPrompt() {
	m.clearPrompt()
}

var ignoreCommandErrors = regexList{
	regexp.MustCompile("^exit status"),
	regexp.MustCompile("^signal: killed"),
}

func (m *Module) handleEventExecCommand(event EventExecCommand) {
	if m.cmd != nil {
		m.Log().ErrorF("Previous command '%s' is still running. Wait for it to finish or cancel it.", m.cmd.Command())
		return
	}

	command := convert.ToString(event.Cmd)
	m.history.Add(command)

	if m.cfg.PrintCommand {
		m.Output().WriteF("[%s]%s%s[-]\n", getColorName(m.cfg.Colors.Command), m.cfg.Label, command)
	}
	// If it's exit command
	if command == "exit" {
		go func() { m.Events().Dispatch(gooster.EventExit{}) }()
		return
	}
	m.clearPrompt()

	// If it looks like "cd" command:
	if path := detectWorkDirPath(m.fs, command); path != "" {
		m.Events().Dispatch(workdir.EventChangeDir{Path: path})
		return
	}

	m.cmd = NewCommand(command).SetOutput(m.Output())
	go func() {
		m.Log().DebugF("Starting command `%s`", command)
		if err := m.cmd.Run(); err != nil {
			if !ignoreCommandErrors.MatchString(err.Error()) {
				m.Log().Error(err)
			}
		}
		m.Log().DebugF("Command finished `%s`", command)
		m.clearCommand()
	}()
}

const newLine byte = 10

func (m *Module) handleEventSendUserInput(event EventSendUserInput) {
	if m.cmd == nil {
		m.Log().Error("Could not send user input - there is no current command running")
		return
	}
	m.clearPrompt()
	if _, err := m.cmd.Write(append([]byte(event.Input), newLine)); err != nil {
		m.Log().Error(err)
	}
}

func (m *Module) handleEventInterruptCommand() {
	if m.cmd == nil {
		return
	}
	m.cmd.Cancel()
	m.clearCommand()
}

func (m *Module) handleKeyHistoryPrev(event *tcell.EventKey) *tcell.EventKey {
	if !m.history.IsActive() {
		m.latestInput = m.view.GetText()
	}
	m.Events().Dispatch(EventSetPrompt{Input: m.history.Prev()})
	return event
}

func (m *Module) handleKeyHistoryNext(event *tcell.EventKey) *tcell.EventKey {
	input := m.history.Next()
	if !m.history.IsActive() {
		input = m.latestInput
	}
	m.Events().Dispatch(EventSetPrompt{Input: input})
	return event
}

func (m *Module) handleCompletion(input string) {
	if m.tabPressed {
		m.tabPressed = false
		m.Events().Dispatch(gooster.EventSetFocusByName{TargetName: "completion"})
		return
	}

	commands, err := command.ParseCommands(input)
	if err != nil {
		m.Log().DebugF("ParseCommands error: %s", err)
	}
	m.Events().Dispatch(gooster.EventSetCompletion{Commands: commands})
	m.tabPressed = true
}
