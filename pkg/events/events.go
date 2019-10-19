package events

type EventId string
type EventPayload interface{}
type Event struct {
	Id      EventId
	Payload EventPayload
}

func (e Event) formattedData() string {
	return truncateString(toString(e.Payload), 80)
}

type Manager interface {
	Dispatch(Event)
	Subscribe(...Subscriber)
	Extend(...Extension)
}

type EventHandler func(Event)

type Subscriber struct {
	Id    EventId
	Fn    EventHandler
	Order float64 // higher value == earlier called
}

type PayloadHandler func(data EventPayload) (newData EventPayload)

type Extension struct {
	Id    EventId
	Fn    PayloadHandler
	Order float64 // higher value == earlier called
}
