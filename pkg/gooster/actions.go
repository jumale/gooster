package gooster

import (
	"fmt"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/log"
	"github.com/rivo/tview"
	"math/rand"
)

type actions struct {
	em          events.Manager
	log         log.Logger
	afterAction func(e events.Event)
	eventIdSeed int
}

func newActions(em events.Manager, log log.Logger) *actions {
	return &actions{
		em:          em,
		log:         log,
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

// ------------------------- APP LEVEL ACTIONS ------------------------- //

const (
	// This is a list of native app events.
	// Do not use these raw event names to dispatch or subscribe
	// to the native events, instead, use the `Dispatch...` and `On...`
	// methods provided by `actions` struct for every event.
	actionAppSetFocus    events.EventId = "app:set_focus"
	actionAppOpenDialog                 = "app:open_dialog"
	actionAppCloseDialog                = "app:close_dialog"
	actionAppExit                       = "app:exit"
)

func (a *actions) OpenDialog(dialog dialog.Dialog) {
	a.Dispatch(events.Event{
		Id:   a.withSeed(actionAppOpenDialog),
		Data: dialog,
	})
}

func (a *actions) CloseDialog() {
	a.Dispatch(events.Event{
		Id: a.withSeed(actionAppCloseDialog),
	})
}

func (a *actions) SetFocus(view tview.Primitive) {
	a.Dispatch(events.Event{
		Id:   a.withSeed(actionAppSetFocus),
		Data: view,
	})
}

func (a *actions) OnSetFocus(fn func(view tview.Primitive)) {
	a.Subscribe(a.withSeed(actionAppSetFocus), events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(tview.Primitive))
		},
	})
}

func (a *actions) Exit() {
	a.Dispatch(events.Event{Id: a.withSeed(actionAppExit)})
}

func (a *actions) OnAppExit(fn func()) {
	a.Subscribe(a.withSeed(actionAppExit), events.Subscriber{
		Handler: func(event events.Event) { fn() },
	})
}

// ---------------------- ACTIONS OWNED BY MODULES --------------------- //

const (
	actionWorkDirChange          events.EventId = "work_dir:change"
	actionOutputWrite                           = "output:write"
	actionPromptChange                          = "prompt:on_change"
	actionPromptInterruptCommand                = "prompt:interrupt_command"
)

type WorkDirOwner interface {
	OnWorkDirSet(path string)
	WorkDirChangeCallback(fn func(path string))
}

type PromptOwner interface {
	OnPromptSet(prompt string)
	PromptChangeCallback(fn func(prompt string))
}

type CommandOwner interface {
	OnCommandInterrupt()
	CommandInterruptCallback(fn func())
}

type OutputOwner interface {
	OnOutputWrite(content []byte)
	OutputWriteCallback(fn func(content interface{}))
}

func (a *actions) registerActionOwners(modules []moduleDefinition) {
	for _, def := range modules {
		switch m := def.module.(type) {

		case WorkDirOwner:
			a.Subscribe(a.internal(actionWorkDirChange), events.Subscriber{
				Handler: func(event events.Event) {
					m.OnWorkDirSet(event.Data.(string))
				},
			})
			m.WorkDirChangeCallback(func(path string) {
				a.Dispatch(events.Event{
					Id:   a.withSeed(actionWorkDirChange),
					Data: path,
				})
			})

		case PromptOwner:
			a.Subscribe(a.internal(actionPromptChange), events.Subscriber{
				Handler: func(event events.Event) {
					m.OnPromptSet(event.Data.(string))
				},
			})
			m.PromptChangeCallback(func(prompt string) {
				a.Dispatch(events.Event{
					Id:   a.withSeed(actionPromptChange),
					Data: prompt,
				})
			})

		case CommandOwner:
			a.Subscribe(a.internal(actionPromptInterruptCommand), events.Subscriber{
				Handler: func(event events.Event) {
					m.OnCommandInterrupt()
				},
			})
			m.CommandInterruptCallback(func() {
				a.Dispatch(events.Event{
					Id: a.withSeed(actionPromptInterruptCommand),
				})
			})

		case OutputOwner:
			a.Subscribe(a.internal(actionOutputWrite), events.Subscriber{
				Handler: func(event events.Event) {
					m.OnOutputWrite(toBytes(event.Data))
				},
			})
			m.OutputWriteCallback(func(content interface{}) {
				a.Dispatch(events.Event{
					Id:   a.withSeed(actionOutputWrite),
					Data: content,
				})
			})
		}
	}
}

func (a *actions) SetWorkDir(path string) {
	a.Dispatch(events.Event{
		Id:   a.internal(actionWorkDirChange),
		Data: path,
	})
}

func (a *actions) SetPrompt(input string) {
	a.Dispatch(events.Event{
		Id:   a.internal(actionPromptChange),
		Data: input,
	})
}

func (a *actions) InterruptLatestCommand() {
	a.Dispatch(events.Event{
		Id: a.internal(actionPromptInterruptCommand),
	})
}

func (a *actions) Write(data interface{}) {
	a.Dispatch(events.Event{
		Id:   a.internal(actionOutputWrite),
		Data: toBytes(data),
	})
}

func (a *actions) Writeln(data interface{}) {
	a.Dispatch(events.Event{
		Id:   a.internal(actionOutputWrite),
		Data: append(toBytes(data), []byte(`\n`)...),
	})
}

func (a *actions) OnWorkDirChange(fn func(newPath string)) {
	a.Subscribe(a.withSeed(actionWorkDirChange), events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(string))
		},
	})
}

func (a *actions) OnPromptChange(fn func(input string)) {
	a.Subscribe(a.withSeed(actionPromptChange), events.Subscriber{
		Handler: func(event events.Event) {
			fn(event.Data.(string))
		},
	})
}

func (a *actions) OnCommandInterrupt(fn func()) {
	a.Subscribe(a.withSeed(actionPromptInterruptCommand), events.Subscriber{
		Handler: func(events.Event) { fn() },
	})
}

func (a *actions) OnOutput(fn func(data []byte)) {
	a.Subscribe(a.withSeed(actionOutputWrite), events.Subscriber{
		Handler: func(event events.Event) {
			fn(toBytes(event.Data))
		},
	})
}

// ------------------------------ HELPERS ------------------------------ //

func (a *actions) withSeed(id events.EventId) events.EventId {
	return events.EventId(fmt.Sprintf("%s#%d", id, a.eventIdSeed))
}

func (a *actions) internal(id events.EventId) events.EventId {
	return events.EventId(fmt.Sprintf("%s_internal#%d", id, a.eventIdSeed))
}

//
//func (a *actions) checkOwnerNotExists(target ownerInterface, m Module) {
//	if prev, exists := a.actionOwners[target]; exists {
//		a.ownerConflict(target, m, prev.(Module))
//	}
//}
//
//func (a *actions) ownerNotFound(event string, target ownerInterface) {
//	a.log.WarnF(
//		"Could not find an owner for triggered event '%s'. One of the registered modules must implement '%s' interface and handle this event.",
//		event,
//		target,
//	)
//}
//
//func (a *actions) ownerConflict(target ownerInterface, newModule Module, existingModule Module) {
//	a.log.ErrorF(
//		"Can not register module '%s' as '%s'. Module '%s' is already took this responsibility. "+
//			"It's only allowed to register one module which implements '%s'",
//		newModule.Name(),
//		target,
//		existingModule.Name(),
//		target,
//	)
//}
