package gooster

import (
	"fmt"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/events"
	"github.com/rivo/tview"
	"math/rand"
)

const (
	// This is a list of native app events.
	// Do not use these raw event names to dispatch or subscribe
	// to the native events, instead, use the `Dispatch...` and `On...`
	// methods provided by `actions` struct for every event.
	eventChangeWorkDir    events.EventId = "change_work_dir"
	eventSendOutput                      = "send_output"
	eventSetPrompt                       = "set_prompt"
	eventCommandInterrupt                = "command_interrupt"
	eventSetFocus                        = "set_focus"
	eventOpenDialog                      = "open_dialog"
	eventCloseDialog                     = "close_dialog"
	eventAppExit                         = "app_exit"
)

type actions struct {
	em          events.Manager
	afterAction func(e events.Event)
	eventIdSeed int
}

func newActions(em events.Manager) *actions {
	return &actions{
		em:          em,
		afterAction: func(e events.Event) {},
		eventIdSeed: rand.Int(),
	}
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
		Id:   a.withSeed(eventChangeWorkDir),
		Data: path,
	})
}

func (a *actions) OnWorkDirChange(fn func(newPath string)) {
	a.Subscribe(a.withSeed(eventChangeWorkDir), events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(string))
		},
	})
}

func (a *actions) SetPrompt(input string) {
	a.Dispatch(events.Event{
		Id:   a.withSeed(eventSetPrompt),
		Data: input,
	})
}

func (a *actions) OnSetPrompt(fn func(input string)) {
	a.Subscribe(a.withSeed(eventSetPrompt), events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(string))
		},
	})
}

func (a *actions) OpenDialog(dialog dialog.Dialog) {
	a.Dispatch(events.Event{
		Id:   a.withSeed(eventOpenDialog),
		Data: dialog,
	})
}

func (a *actions) CloseDialog() {
	a.Dispatch(events.Event{
		Id: a.withSeed(eventCloseDialog),
	})
}

func (a *actions) SetFocus(view tview.Primitive) {
	a.Dispatch(events.Event{
		Id:   a.withSeed(eventSetFocus),
		Data: view,
	})
}

func (a *actions) OnSetFocus(fn func(view tview.Primitive)) {
	a.Subscribe(a.withSeed(eventSetFocus), events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(tview.Primitive))
		},
	})
}

func (a *actions) Write(data interface{}) {
	a.Dispatch(events.Event{
		Id:   a.withSeed(eventSendOutput),
		Data: data,
	})
}

func (a *actions) Writeln(data interface{}) {
	a.Dispatch(events.Event{
		Id:   a.withSeed(eventSendOutput),
		Data: append(toBytes(data), []byte(`\n`)...),
	})
}

func (a *actions) OnOutput(fn func(data []byte)) {
	a.Subscribe(a.withSeed(eventSendOutput), events.Subscriber{
		Handler: func(event events.Event) {
			fn(toBytes(event.Data))
		},
	})
}

func (a *actions) InterruptLatestCommand() {
	a.Dispatch(events.Event{Id: a.withSeed(eventCommandInterrupt)})
}

func (a *actions) OnCommandInterrupt(fn func()) {
	a.Subscribe(a.withSeed(eventCommandInterrupt), events.Subscriber{
		Handler: func(events.Event) { fn() },
	})
}

func (a *actions) Exit() {
	a.Dispatch(events.Event{Id: a.withSeed(eventAppExit)})
}

func (a *actions) OnAppExit(fn func()) {
	a.Subscribe(a.withSeed(eventAppExit), events.Subscriber{
		Handler: func(event events.Event) { fn() },
	})
}

func (a *actions) withSeed(id events.EventId) events.EventId {
	return events.EventId(fmt.Sprintf("%s#%d", id, a.eventIdSeed))
}
