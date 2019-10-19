package output

import (
	"github.com/jumale/gooster/pkg/convert"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
)

const (
	ActionWrite events.EventId = "output:write"
)

type Actions struct {
	*gooster.AppContext
}

func (a Actions) Write(data interface{}) {
	a.Events().Dispatch(events.Event{Id: ActionWrite, Payload: convert.ToBytes(data)})
}

func (a Actions) WriteLine(data interface{}) {
	a.Events().Dispatch(events.Event{Id: ActionWrite, Payload: append(convert.ToBytes(data), []byte("\n")...)})
}
