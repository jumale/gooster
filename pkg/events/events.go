package events

// EventHandler performs some action when an event is dispatched.
// It must return a new (or the same event).
// The returned event will be passed to the next subscribers,
// this feature can be used to extend events.
// By returning nil you stop the propagation, so the next
// subscribers will not be called for this event.
type EventHandler func(event IEvent) IEvent

var StopPropagation IEvent = nil

type IEvent interface{}

type Manager interface {
	Dispatch(IEvent)
	Subscribe(ISubscriber)
}

type ISubscriber interface {
	Handler() EventHandler
	Priority() float64 // higher value == earlier called
}

func HandleFunc(fn EventHandler) ISubscriber {
	return subscriber{handler: fn}
}

func HandleWithPrio(prio float64, fn EventHandler) ISubscriber {
	return subscriber{handler: fn, prio: prio}
}

type subscriber struct {
	handler EventHandler
	prio    float64
}

func (s subscriber) Handler() EventHandler {
	return s.handler
}

func (s subscriber) Priority() float64 {
	return s.prio
}

type GenericEvent struct {
	ID string
}
