package events

import (
	"fmt"
	_assert "github.com/stretchr/testify/assert"
	"testing"
)

const (
	EventBird EventId = "bird"
	EventFish         = "fish"
)

func TestSubscribe(t *testing.T) {
	assert := _assert.New(t)
	cfg := ManagerConfig{}

	t.Run("should let subscriber to handle an event", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(EventBird, Subscriber{
			Handler: handled.withName("bob"),
		})

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})

		assert.Len(handled.events, 1)
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[0])
	})

	t.Run("should let multiple subscribers to handle multiple events", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(EventBird, Subscriber{
			Handler: handled.withName("bob"),
		})
		mng.Subscribe(EventBird, Subscriber{
			Handler: handled.withName("john"),
		})

		mng.Subscribe(EventFish, Subscriber{
			Handler: handled.withName("eric"),
		})

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})
		mng.Dispatch(Event{Id: EventFish, Data: "bass"})
		mng.Dispatch(Event{Id: EventFish, Data: "tuna"})

		assert.Len(handled.events, 4)
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[0])
		assert.Equal("'john' handled 'bird/eagle'", handled.events[1])
		assert.Equal("'eric' handled 'fish/bass'", handled.events[2])
		assert.Equal("'eric' handled 'fish/tuna'", handled.events[3])
	})

	t.Run("should consider priority while calling subscribers", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(EventBird, Subscriber{
			Handler:  handled.withName("bob"),
			Priority: -10,
		})
		mng.Subscribe(EventBird, Subscriber{
			Handler: handled.withName("john"),
		})
		mng.Subscribe(EventBird, Subscriber{
			Handler:  handled.withName("eric"),
			Priority: 9999,
		})

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})

		assert.Len(handled.events, 3)
		assert.Equal("'eric' handled 'bird/eagle'", handled.events[0])
		assert.Equal("'john' handled 'bird/eagle'", handled.events[1])
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[2])
	})

	t.Run("should not change order of non-prioritised subscribers", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(EventBird, Subscriber{
			Handler: handled.withName("bob"),
		})
		mng.Subscribe(EventBird, Subscriber{
			Handler: handled.withName("john"),
		})
		mng.Subscribe(EventBird, Subscriber{
			Handler:  handled.withName("eric"),
			Priority: 9999,
		})

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})

		assert.Len(handled.events, 3)
		assert.Equal("'eric' handled 'bird/eagle'", handled.events[0])
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[1])
		assert.Equal("'john' handled 'bird/eagle'", handled.events[2])
	})
}

func TestDelayedStart(t *testing.T) {
	assert := _assert.New(t)

	t.Run("should consume events immediately if delay is not defined", func(t *testing.T) {
		cfg := ManagerConfig{}
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(EventBird, Subscriber{
			Handler: handled.withName("bob"),
		})

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})

		assert.Len(handled.events, 1)
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[0])
	})

	t.Run("should delay handling events if delay is defined", func(t *testing.T) {
		cfg := ManagerConfig{
			DelayedStart: true,
		}

		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(EventBird, Subscriber{
			Handler: handled.withName("bob"),
		})

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})
		assert.Len(handled.events, 0, "should not handle events before delayed manager has started")

		// let's add another EventBird handler after the actual event has been dispatched
		mng.Subscribe(EventBird, Subscriber{
			Handler: handled.withName("john"),
		})

		// now all events should be handled by all subscribers
		mng.Start()

		assert.Len(handled.events, 2)
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[0])
		assert.Equal("'john' handled 'bird/eagle'", handled.events[1], "john should be still able to handle the delayed events even though he was subscribed after the event has been dispatched")
	})
}

type handler struct {
	events []string
}

func (h *handler) withName(name string) Handler {
	return func(event Event) {
		e := fmt.Sprintf("'%s' handled '%s/%s'", name, event.Id, event.Data)
		h.events = append(h.events, e)
	}
}
