package gooster

import (
	"github.com/jumale/gooster/pkg/events"
	"github.com/rivo/tview"
)

const (
	eventChangeWorkDir    events.EventId = "change_work_dir"
	eventSendOutput                      = "send_output"
	eventSetPrompt                       = "set_prompt"
	eventCommandInterrupt                = "command_interrupt"
	eventSetFocus                        = "set_focus"
	eventAppExit                         = "app_exit"
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
		Id:   eventChangeWorkDir,
		Data: path,
	})
}

func (a *actions) OnWorkDirChange(fn func(newPath string)) {
	a.Subscribe(eventChangeWorkDir, events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(string))
		},
	})
}

func (a *actions) SetPrompt(input string) {
	a.Dispatch(events.Event{
		Id:   eventSetPrompt,
		Data: input,
	})
}

func (a *actions) OnSetPrompt(fn func(input string)) {
	a.Subscribe(eventSetPrompt, events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(string))
		},
	})
}

func (a *actions) SetFocus(view tview.Primitive) {
	a.Dispatch(events.Event{
		Id:   eventSetFocus,
		Data: view,
	})
}

func (a *actions) OnSetFocus(fn func(view tview.Primitive)) {
	a.Subscribe(eventSetFocus, events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(tview.Primitive))
		},
	})
}

func (a *actions) Write(data interface{}) {
	a.Dispatch(events.Event{
		Id:   eventSendOutput,
		Data: data,
	})
}

func (a *actions) Writeln(data interface{}) {
	a.Dispatch(events.Event{
		Id:   eventSendOutput,
		Data: append(toBytes(data), []byte(`\n`)...),
	})
}

func (a *actions) OnOutput(fn func(data []byte)) {
	a.Subscribe(eventSendOutput, events.Subscriber{
		Handler: func(event events.Event) {
			fn(toBytes(event.Data))
		},
	})
}

func (a *actions) InterruptLatestCommand() {
	a.Dispatch(events.Event{Id: eventCommandInterrupt})
}

func (a *actions) OnCommandInterrupt(fn func()) {
	a.Subscribe(eventCommandInterrupt, events.Subscriber{
		Handler: func(events.Event) { fn() },
	})
}

func (a *actions) Exit() {
	a.Dispatch(events.Event{Id: eventAppExit})
}

func (a *actions) OnAppExit(fn func()) {
	a.Subscribe(eventAppExit, events.Subscriber{
		Handler: func(event events.Event) { fn() },
	})
}
