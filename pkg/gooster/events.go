package gooster

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	EventWidgetUpdate  EventId = "widget_update"
	EventWorkDirChange         = "work_dir_change"
	EventOutputMessage         = "output_message"
)

type EventId string
type subscriberId string

type Event struct {
	Id   EventId
	Data interface{}
}

type EventSubscriber func(Event)

type subscribers map[subscriberId]EventSubscriber

type EventManager struct {
	events   chan Event
	sub      map[EventId]subscribers
	mu       *sync.Mutex
	eventLog *os.File
}

func NewEventManager(logFile string) (*EventManager, error) {
	em := &EventManager{
		events: make(chan Event, 10),
		sub:    make(map[EventId]subscribers),
		mu:     &sync.Mutex{},
	}
	if logFile != "" {
		if err := em.initEventLog(logFile); err != nil {
			return nil, err
		}
	}

	return em, nil
}

func (em *EventManager) Start() {
	go func() {
		for e := range em.events {
			for _, subscriber := range em.subscribers(e.Id) {
				subscriber(e)
			}
		}
	}()
}

func (em *EventManager) Dispatch(e Event) {
	em.log("dispatch   %-20s %s", e.Id, e.formattedData())
	em.events <- e
}

func (em *EventManager) Subscribe(id EventId, es EventSubscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, ok := em.sub[id]; !ok {
		em.sub[id] = make(subscribers)
	}

	subId := em.getSubscriberId()
	em.log("subscribe  %-20s to %s", subId, id)

	em.sub[id][subId] = func(e Event) {
		em.log("consume    %-20s by %-20s %s", e.Id, subId, e.formattedData())
		es(e)
	}

}

func (em *EventManager) Close() error {
	em.log("Closing Event Manager")
	close(em.events)
	return em.eventLog.Close()
}

func (em *EventManager) subscribers(id EventId) subscribers {
	em.mu.Lock()
	defer em.mu.Unlock()

	if s, ok := em.sub[id]; ok {
		return s

	} else {
		em.sub[id] = make(subscribers)
		return em.sub[id]
	}
}

func (em *EventManager) initEventLog(path string) (err error) {
	if em.eventLog, err = os.Create(path); err != nil {
		return err
	}
	em.log("Initializing Event Manager")
	return nil
}

func (em *EventManager) log(msg string, args ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	args = append([]interface{}{now}, args...)
	_, _ = fmt.Fprintf(em.eventLog, "[%s] "+msg+"\n", args...)
}

func (em *EventManager) getSubscriberId() subscriberId {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}

	parts := strings.Split(file, "/")
	return subscriberId(fmt.Sprintf("%s:%d", parts[len(parts)-1], line))
}

func (e Event) formattedData() string {
	return truncateString(fmt.Sprintf("%+v", e.Data), 100)
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}
