package events

import (
	"sort"
	"sync"
)

type ManagerConfig struct {
	// DelayedStart defines whether events will be dispatched immediately,
	// of buffered till Manager.Start() is called.
	// It allows to delay starting the event manager without loosing any events.
	DelayedStart bool

	// BeforeEvent is applied on every event, before it's dispatched,
	// and if returns false then the event will not be dispatched.
	BeforeEvent func(Event) bool

	// AfterEvent is applied on every event after it's dispatched.
	AfterEvent func(Event)
}

type DefaultManager struct {
	cfg ManagerConfig
	sub map[EventId][]Subscriber
	ext map[EventId][]Extension
	mu  *sync.Mutex

	filter func(Event) bool

	// support for delayed start
	buffer  []Event
	started bool
}

func NewManager(cfg ManagerConfig) (*DefaultManager, error) {
	if cfg.BeforeEvent == nil {
		cfg.BeforeEvent = func(Event) bool { return true }
	}
	if cfg.AfterEvent == nil {
		cfg.AfterEvent = func(Event) {}
	}

	em := &DefaultManager{
		cfg:     cfg,
		sub:     make(map[EventId][]Subscriber),
		ext:     make(map[EventId][]Extension),
		mu:      &sync.Mutex{},
		filter:  func(Event) bool { return true },
		started: !cfg.DelayedStart,
	}

	return em, nil
}

func (em *DefaultManager) Dispatch(e Event) {
	if !em.cfg.BeforeEvent(e) {
		return
	}

	if !em.started {
		em.mu.Lock()
		em.buffer = append(em.buffer, e)
		em.mu.Unlock()
		return
	}

	em.mu.Lock()
	subscribers, ok := em.sub[e.Id]
	em.mu.Unlock()
	if !ok {
		return
	}

	e = em.extendEvent(e)
	for _, sub := range subscribers {
		sub.Fn(e)
	}

	em.cfg.AfterEvent(e)
}

func (em *DefaultManager) Subscribe(subscribers ...Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()

	var changedEvents []EventId
	for _, es := range subscribers {
		id := es.Id
		changedEvents = append(changedEvents, id)
		if _, ok := em.sub[id]; !ok {
			em.sub[id] = []Subscriber{}
		}
		em.sub[id] = append(em.sub[id], es)
	}

	for _, id := range changedEvents {
		sort.SliceStable(em.sub[id], func(i, j int) bool {
			return em.sub[id][i].Order > em.sub[id][j].Order
		})
	}
}

func (em *DefaultManager) Extend(extensions ...Extension) {
	em.mu.Lock()
	defer em.mu.Unlock()

	var changedEvents []EventId
	for _, ext := range extensions {
		id := ext.Id
		changedEvents = append(changedEvents, id)
		if _, ok := em.ext[id]; !ok {
			em.ext[id] = []Extension{}
		}
		em.ext[id] = append(em.ext[id], ext)
	}

	for _, id := range changedEvents {
		sort.SliceStable(em.ext[id], func(i, j int) bool {
			return em.ext[id][i].Order > em.ext[id][j].Order
		})
	}
}

func (em *DefaultManager) Start() {
	em.started = true
	for _, event := range em.buffer {
		em.Dispatch(event)
	}
}

func (em *DefaultManager) Close() error {
	return nil
}

func (em *DefaultManager) extendEvent(originalEvent Event) (extendedEvent Event) {
	data := originalEvent.Payload
	if extensions, ok := em.ext[originalEvent.Id]; ok {
		for _, extension := range extensions {
			data = extension.Fn(data)
		}
	}
	return Event{
		Id:      originalEvent.Id,
		Payload: data,
	}
}
