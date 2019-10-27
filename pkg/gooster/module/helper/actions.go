package helper

import (
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
)

const (
	ActionSetCompletion events.EventId = "helper:set_completion"
)

type Actions struct {
	*gooster.AppContext
}

type PayloadSetCompletion struct {
	Input      string
	Completion []string
}

func (a Actions) SetCompletion(input string, completion []string) {
	a.Events().Dispatch(events.Event{Id: ActionSetCompletion, Payload: PayloadSetCompletion{
		Input:      input,
		Completion: completion,
	}})
}
