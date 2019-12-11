package command

type CompletionType uint8

const (
	CompleteCommand CompletionType = iota
	CompleteArg
	CompleteVar
	CompleteDir
	CompleteFile
	CompleteCustom
)

type Completion struct {
	Type   CompletionType
	Values []string
}

func (e Completion) Empty() bool {
	return len(e.Values) == 0
}

func (e Completion) HasSingle() bool {
	return len(e.Values) == 1
}

type Completer interface {
	Get(cmd Definition) (Completion, error)
}

func ApplyCompletion(cmd string, completion string, completionType CompletionType) string {
	var prefix string
	switch completionType {
	case CompleteDir:
		prefix = "/"
	case CompleteCommand, CompleteArg, CompleteVar, CompleteCustom:
		prefix = " "
	}

	for i := len(cmd) - 1; i >= 0; i -= 1 {
		if rune(cmd[i]) == space && (i-1 < 0 || rune(cmd[i-1]) != escape) {
			return cmd[:i+1] + completion + prefix
		}
	}
	return completion + prefix
}
