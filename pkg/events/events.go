package events

type EventId string
type EventPayload interface{}
type Event struct {
	Id   EventId
	Data EventPayload
}

func (e Event) formattedData() string {
	return truncateString(toString(e.Data), 80)
}

type Manager interface {
	Dispatch(Event)
	Subscribe(EventId, Subscriber)
	Extend(EventId, Extension)
}

type Handler func(Event)

type Subscriber struct {
	Handle   Handler
	Priority float64
}

type Extender func(data EventPayload) (newData EventPayload)

type Extension struct {
	Extend   Extender
	Priority float64
}
