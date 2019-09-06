package gooster

import (
	"github.com/jumale/gooster/pkg/events"
)

const (
	EventChangeWorkDir events.EventId = "change_work_dir"
	EventSendOutput                   = "send_output"
)

type actions struct {
	em          *events.Manager
	afterAction func(e events.Event)
}

func (a *actions) Dispatch(e events.Event) {
	a.em.Dispatch(e)
	a.afterAction(e)
}

func (a *actions) Subscribe(id events.EventId, es events.Subscriber) {
	a.em.Subscribe(id, es)
}

func (a *actions) SetWorkDir(path string) {
	a.Dispatch(events.Event{
		Id:   EventChangeWorkDir,
		Data: path,
	})
}

func (a *actions) OnWorkDirChange(fn func(newPath string)) {
	a.Subscribe(EventChangeWorkDir, func(event events.Event) {
		fn(event.Data.(string))
	})
}

func (a *actions) Write(data interface{}) {
	a.Dispatch(events.Event{
		Id:   EventSendOutput,
		Data: data,
	})
}

func (a *actions) Writeln(data interface{}) {
	a.Dispatch(events.Event{
		Id:   EventSendOutput,
		Data: append(toBytes(data), []byte(`\n`)...),
	})
}

func (a *actions) OnOutput(fn func(data []byte)) {
	a.Subscribe(EventSendOutput, func(event events.Event) {
		fn(toBytes(event.Data))
	})
}
