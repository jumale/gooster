package testtools

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
	"strings"
)

func TestableModule(m gooster.Module) *ModuleTester {
	ctx, err := gooster.NewAppContext(
		gooster.AppContextConfig{
			LogLevel:          log.Debug,
			DelayEventManager: false,
		},
	)
	if err != nil {
		panic(errors.WithMessagef(err, "could not create a module tester for %T", m))
	}

	err = m.Init(ctx)
	if err != nil {
		panic(errors.WithMessagef(err, "could not init module %T", m))
	}

	tester := &ModuleTester{
		AppContext: ctx,
		module:     m,
		screen:     NewScreenStub(10, 10),
	}

	ctx.Events().Subscribe(events.Subscriber{
		Id:    "app:draw",
		Fn:    func(event events.Event) { tester.draw() },
		Order: -9999,
	})

	return tester
}

type ModuleTester struct {
	*gooster.AppContext
	module gooster.Module
	screen *ScreenStub
}

func (t *ModuleTester) SetSize(width, height int) *ModuleTester {
	t.screen = NewScreenStub(width, height)
	return t
}

func (t *ModuleTester) PressKey(key tcell.Key) *ModuleTester {

	return t
}

func (t *ModuleTester) SendEvent(id events.EventId, payload ...events.EventPayload) *ModuleTester {
	return t
}

func (t *ModuleTester) GetDisplay() []byte {
	return []byte(t.GetDisplayString())
}

func (t *ModuleTester) GetDisplayString() string {
	var data []string
	for _, row := range t.screen.data {
		data = append(data, string(row))
	}
	return strings.Join(data, "\n")
}

func (t *ModuleTester) draw() {
	t.module.Draw(t.screen)
}
