package events

type EventId string
type subscriberId string

type Event struct {
	Id   EventId
	Data interface{}
}

type Subscriber func(Event)

func (e Event) formattedData() string {
	return truncateString(toString(e.Data), 80)
}
