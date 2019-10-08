package events

import (
	"github.com/jumale/gooster/pkg/log"
	"sort"
	"sync"
)

type ManagerConfig struct {
	DelayedStart bool
	Log          log.Logger
}

type DefaultManager struct {
	cfg ManagerConfig
	sub map[EventId][]Subscriber
	ext map[EventId][]Extension
	log log.Logger
	mu  *sync.Mutex

	// support for delayed start
	buffer  []Event
	started bool
}

func NewManager(cfg ManagerConfig) (*DefaultManager, error) {
	em := &DefaultManager{
		cfg:     cfg,
		sub:     make(map[EventId][]Subscriber),
		ext:     make(map[EventId][]Extension),
		mu:      &sync.Mutex{},
		log:     log.EmptyLogger{},
		started: !cfg.DelayedStart,
	}
	if cfg.Log != nil {
		em.log = cfg.Log
	}
	em.log.Info("Event Manager: initialized")

	return em, nil
}

func (em *DefaultManager) Dispatch(e Event) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if !em.started {
		em.buffer = append(em.buffer, e)
		return
	}

	em.log.DebugF("Event Manager: dispatched %s %s", e.Id, e.formattedData())
	e = em.extendEvent(e)
	em.log.DebugF("Event Manager: modified %s to %s", e.Id, e.formattedData())

	if subscribers, ok := em.sub[e.Id]; ok {
		for _, sub := range subscribers {
			sub.Handle(e)
		}
	}
}

func (em *DefaultManager) Subscribe(id EventId, es Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, ok := em.sub[id]; !ok {
		em.sub[id] = []Subscriber{}
	}
	em.sub[id] = append(em.sub[id], es)

	sort.SliceStable(em.sub[id], func(i, j int) bool {
		return em.sub[id][i].Priority > em.sub[id][j].Priority
	})
}

func (em *DefaultManager) Extend(id EventId, ext Extension) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, ok := em.ext[id]; !ok {
		em.ext[id] = []Extension{}
	}
	em.ext[id] = append(em.ext[id], ext)

	sort.SliceStable(em.ext[id], func(i, j int) bool {
		return em.ext[id][i].Priority > em.ext[id][j].Priority
	})
}

func (em *DefaultManager) Start() {
	em.started = true
	for _, event := range em.buffer {
		em.Dispatch(event)
	}
	em.log.Info("Event Manager: started")
}

func (em *DefaultManager) Close() error {
	em.log.Info("Event Manager: closed")
	return nil
}

func (em *DefaultManager) extendEvent(originalEvent Event) (extendedEvent Event) {
	data := originalEvent.Data
	if extensions, ok := em.ext[originalEvent.Id]; ok {
		for _, extension := range extensions {
			data = extension.Extend(data)
		}
	}
	return Event{
		Id:   originalEvent.Id,
		Data: data,
	}
}
