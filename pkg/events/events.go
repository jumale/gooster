package events

type EventId string
type subscriberId string

type Event struct {
	Id   EventId
	Data interface{}
}

type Handler func(Event)

type Manager interface {
	Dispatch(Event)
	Subscribe(EventId, Subscriber)
}

type Subscriber struct {
	Handler  Handler
	Priority int
}

func (e Event) formattedData() string {
	return truncateString(toString(e.Data), 80)
}
