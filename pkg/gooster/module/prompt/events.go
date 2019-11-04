package prompt

type EventSetPrompt struct {
	Input string
}

func (e EventSetPrompt) NeedsDraw() bool {
	return true
}

type EventClearPrompt struct{}

func (e EventClearPrompt) NeedsDraw() bool {
	return true
}

type EventExecCommand struct {
	Cmd string
}

type EventSendUserInput struct {
	Input string
}
