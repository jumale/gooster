package prompt

type history struct {
	list []string
}

func newHistory(historyFile string) *history {
	return &history{}
}

func (h *history) add(cmd string) {

}
