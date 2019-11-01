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
}

type DefaultManager struct {
	cfg ManagerConfig
	sub []ISubscriber
	mu  *sync.Mutex

	filter func(IEvent) bool

	// support for delayed start
	buffer  []IEvent
	started bool
}

func NewManager(cfg ManagerConfig) (*DefaultManager, error) {
	em := &DefaultManager{
		cfg:     cfg,
		mu:      &sync.Mutex{},
		filter:  func(IEvent) bool { return true },
		started: !cfg.DelayedStart,
	}

	return em, nil
}

func (em *DefaultManager) Dispatch(e IEvent) {
	if !em.started {
		em.mu.Lock()
		em.buffer = append(em.buffer, e)
		em.mu.Unlock()
		return
	}

	em.mu.Lock()
	subscribers := em.sub
	em.mu.Unlock()

	for _, sub := range subscribers {
		e = sub.Handler()(e)
		if e == nil {
			return
		}
	}
}

func (em *DefaultManager) Subscribe(subscriber ISubscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.sub = append(em.sub, subscriber)
	sort.SliceStable(em.sub, func(i, j int) bool {
		return em.sub[i].Priority() > em.sub[j].Priority()
	})
}

func (em *DefaultManager) Init() error {
	em.started = true
	for _, event := range em.buffer {
		em.Dispatch(event)
	}
	return nil
}

func (em *DefaultManager) Close() error {
	return nil
}
