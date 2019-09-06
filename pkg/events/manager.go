package events

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type ManagerConfig struct {
	LogFile              string
	SubscriberStackLevel int
}

type Manager struct {
	cfg      ManagerConfig
	sub      map[EventId][]Subscriber
	buffer   []Event
	started  bool
	mu       *sync.Mutex
	eventLog *os.File
}

func NewManager(cfg ManagerConfig) (*Manager, error) {
	em := &Manager{
		cfg: cfg,
		sub: make(map[EventId][]Subscriber),
		mu:  &sync.Mutex{},
	}
	if cfg.LogFile != "" {
		if err := em.initEventLog(cfg.LogFile); err != nil {
			return nil, err
		}
	}

	return em, nil
}

func (em *Manager) Dispatch(e Event) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if !em.started {
		em.buffer = append(em.buffer, e)
		return
	}

	em.log("\033[33mdispatched \033[32m%s \033[94m%s \033[0m", e.Id, e.formattedData())

	if subscribers, ok := em.sub[e.Id]; ok {
		for _, sub := range subscribers {
			sub(e)
		}
	}
}

func (em *Manager) Subscribe(id EventId, es Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, ok := em.sub[id]; !ok {
		em.sub[id] = []Subscriber{}
	}

	subId := callAsSubscriberId(em.cfg.SubscriberStackLevel)
	em.log("\033[34m%s \033[33msubscribes to \033[32m%s \033[0m", subId, id)

	em.sub[id] = append(em.sub[id], func(e Event) {
		em.log("\033[34m%s \033[33mconsumes \033[32m%s \033[94m%s \033[0m", subId, e.Id, e.formattedData())
		es(e)
	})
}

func (em *Manager) Start() {
	em.started = true
	for _, event := range em.buffer {
		em.Dispatch(event)
	}
}

func (em *Manager) Close() error {
	em.log("Closing Event Manager")
	return em.eventLog.Close()
}

func (em *Manager) initEventLog(path string) (err error) {
	if em.eventLog, err = os.Create(path); err != nil {
		return err
	}
	em.log("Initializing Event Manager")
	return nil
}

func (em *Manager) log(msg string, args ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", time.Now().Format("2006-01-02 15:04:05"))
	msg = fmt.Sprintf(prefix+" "+msg, args...)
	msg = strings.ReplaceAll(msg, "\n", "")
	msg += "\n"

	_, _ = fmt.Fprint(em.eventLog, msg)
}
