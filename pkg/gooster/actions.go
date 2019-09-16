package gooster

import (
	"github.com/jumale/gooster/pkg/events"
)

const (
	EventChangeWorkDir    events.EventId = "change_work_dir"
	EventSendOutput                      = "send_output"
	EventCommandInterrupt                = "command_interrupt"
	EventAppExit                         = "app_exit"
)

type actions struct {
	em          events.Manager
	afterAction func(e events.Event)
}

func (a *actions) Dispatch(e events.Event) {
	go func() {
		a.em.Dispatch(e)
		a.afterAction(e)
	}()
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
	a.Subscribe(EventChangeWorkDir, events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(string))
		},
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
	a.Subscribe(EventSendOutput, events.Subscriber{
		Handler: func(event events.Event) {
			fn(toBytes(event.Data))
		},
	})
}

func (a *actions) InterruptLatestCommand() {
	a.Dispatch(events.Event{Id: EventCommandInterrupt})
}

func (a *actions) OnCommandInterrupt(fn func()) {
	a.Subscribe(EventCommandInterrupt, events.Subscriber{
		Handler: func(events.Event) { fn() },
	})
}

func (a *actions) Exit() {
	a.Dispatch(events.Event{Id: EventAppExit})
}

func (a *actions) OnAppExit(fn func()) {
	a.Subscribe(EventAppExit, events.Subscriber{
		Handler: func(event events.Event) { fn() },
	})
}
