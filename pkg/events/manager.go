package events

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type ManagerConfig struct {
	LogFile              string
	SubscriberStackLevel int
	DelayedStart         bool
}

type DefaultManager struct {
	cfg      ManagerConfig
	sub      map[EventId][]Subscriber
	buffer   []Event
	started  bool
	mu       *sync.Mutex
	eventLog io.WriteCloser
}

func NewManager(cfg ManagerConfig) (*DefaultManager, error) {
	em := &DefaultManager{
		cfg:     cfg,
		sub:     make(map[EventId][]Subscriber),
		mu:      &sync.Mutex{},
		started: !cfg.DelayedStart,
	}
	if cfg.LogFile != "" {
		if err := em.initEventLog(cfg.LogFile); err != nil {
			return nil, errors.WithMessage(err, "init event log")
		}
	}

	return em, nil
}

func (em *DefaultManager) Dispatch(e Event) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if !em.started {
		em.buffer = append(em.buffer, e)
		return
	}

	em.log("\033[33mdispatched \033[32m%s \033[94m%s \033[0m", e.Id, e.formattedData())

	if subscribers, ok := em.sub[e.Id]; ok {
		for _, sub := range subscribers {
			sub.Handler(e)
		}
	}
}

func (em *DefaultManager) Subscribe(id EventId, es Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, ok := em.sub[id]; !ok {
		em.sub[id] = []Subscriber{}
	}

	subId := getSubscriberIdFromCaller(em.cfg.SubscriberStackLevel)
	em.log("\033[34m%s \033[33msubscribes to \033[32m%s \033[0m", subId, id)

	em.sub[id] = append(em.sub[id], Subscriber{
		Handler: func(e Event) {
			em.log("\033[34m%s \033[33mconsumes \033[32m%s \033[94m%s \033[0m", subId, e.Id, e.formattedData())
			es.Handler(e)
		},
		Priority: es.Priority,
	})

	sort.SliceStable(em.sub[id], func(i, j int) bool {
		return em.sub[id][i].Priority > em.sub[id][j].Priority
	})
}

func (em *DefaultManager) Start() {
	em.started = true
	for _, event := range em.buffer {
		em.Dispatch(event)
	}
}

func (em *DefaultManager) Close() error {
	em.log("Closing Event Manager")
	if em.eventLog != nil {
		if err := em.eventLog.Close(); err != nil {
			return errors.WithMessage(err, "close event-log file")
		}
	}
	return nil
}

func (em *DefaultManager) initEventLog(path string) (err error) {
	if em.eventLog, err = os.Create(path); err != nil {
		return errors.WithMessage(err, "open event-log file")
	}
	em.log("Initializing Event Manager")
	return nil
}

func (em *DefaultManager) log(msg string, args ...interface{}) {
	if em.eventLog == nil {
		return
	}

	prefix := fmt.Sprintf("[%s] ", time.Now().Format("2006-01-02 15:04:05"))
	msg = fmt.Sprintf(prefix+" "+msg, args...)
	msg = strings.ReplaceAll(msg, "\n", "")
	msg += "\n"

	_, _ = fmt.Fprint(em.eventLog, msg)
}
