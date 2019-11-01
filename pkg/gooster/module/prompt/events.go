package prompt

type EventSetPrompt struct {
	Input string
}

type EventClearPrompt struct{}

type EventExecCommand struct {
	Cmd string
}

type EventSendUserInput struct {
	Input string
}
