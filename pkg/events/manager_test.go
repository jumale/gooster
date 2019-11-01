package events

import (
	"fmt"
	_assert "github.com/stretchr/testify/assert"
	"testing"
)

const (
	EventBird string = "bird"
	EventFish        = "fish"
)

type Event struct {
	Id      string
	Payload string
}

func TestConstructor(t *testing.T) {
	assert := _assert.New(t)
	t.Run("should create a new instance of manager interface", func(t *testing.T) {
		mng, err := NewManager(ManagerConfig{})
		assert.NoError(err)
		assert.Implements((*Manager)(nil), mng)
	})
}

func TestSubscribe(t *testing.T) {
	assert := _assert.New(t)
	cfg := ManagerConfig{}

	t.Run("should do nothing if there are no subscribers", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)
		handled := &handler{}

		mng.Dispatch(Event{Id: EventBird, Payload: "eagle"})

		assert.Len(handled.events, 0)
	})

	t.Run("should let subscriber to handle an event", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(handled.withName(EventBird, "bob"))

		mng.Dispatch(Event{Id: EventBird, Payload: "eagle"})

		assert.Len(handled.events, 1)
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[0])
	})

	t.Run("should let multiple subscribers to handle multiple events", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(handled.withName(EventBird, "bob"))
		mng.Subscribe(handled.withName(EventBird, "john"))
		mng.Subscribe(handled.withName(EventFish, "eric"))

		mng.Dispatch(Event{Id: EventBird, Payload: "eagle"})
		mng.Dispatch(Event{Id: EventFish, Payload: "bass"})
		mng.Dispatch(Event{Id: EventFish, Payload: "tuna"})

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
		mng.Subscribe(handled.withNameAndPriority(EventBird, "bob", -10))
		mng.Subscribe(handled.withNameAndPriority(EventBird, "eric", 9999))
		mng.Subscribe(handled.withNameAndPriority(EventBird, "john", 0))

		mng.Dispatch(Event{Id: EventBird, Payload: "eagle"})

		assert.Len(handled.events, 3)
		assert.Equal("'eric' handled 'bird/eagle'", handled.events[0])
		assert.Equal("'john' handled 'bird/eagle'", handled.events[1])
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[2])
	})

	t.Run("should not change order of non-prioritised subscribers", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(handled.withName(EventBird, "bob"))
		mng.Subscribe(handled.withName(EventBird, "john"))
		mng.Subscribe(handled.withNameAndPriority(EventBird, "eric", 9999))

		mng.Dispatch(Event{Id: EventBird, Payload: "eagle"})

		assert.Len(handled.events, 3)
		assert.Equal("'eric' handled 'bird/eagle'", handled.events[0])
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[1])
		assert.Equal("'john' handled 'bird/eagle'", handled.events[2])
	})
}

func TestExtendEvents(t *testing.T) {
	assert := _assert.New(t)
	cfg := ManagerConfig{}

	t.Run("should extend events and modify event data before dispatching it", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}

		mng.Subscribe(extension.withName(EventBird, "bird_ext1"))
		mng.Subscribe(extension.withName(EventBird, "bird_ext2"))
		mng.Subscribe(handled.withName(EventBird, "bob"))

		mng.Subscribe(extension.withName(EventFish, "fish_ext1"))
		mng.Subscribe(handled.withName(EventFish, "eric"))

		mng.Dispatch(Event{Id: EventBird, Payload: "eagle"})
		mng.Dispatch(Event{Id: EventFish, Payload: "tuna"})

		assert.Len(handled.events, 2)
		assert.Equal("'bob' handled 'bird/bird_ext2/bird_ext1/eagle'", handled.events[0])
		assert.Equal("'eric' handled 'fish/fish_ext1/tuna'", handled.events[1])
	})

	t.Run("should consider priority while applying extensions", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(handled.withName(EventBird, "bob"))

		mng.Subscribe(extension.withNameAndPriority(EventBird, "ext1", 99))
		mng.Subscribe(extension.withNameAndPriority(EventBird, "ext2", 9999))
		mng.Subscribe(extension.withNameAndPriority(EventBird, "ext3", 999))

		mng.Dispatch(Event{Id: EventBird, Payload: "eagle"})

		assert.Len(handled.events, 1)
		assert.Equal("'bob' handled 'bird/ext1/ext3/ext2/eagle'", handled.events[0])
	})
}

func TestDelayedStart(t *testing.T) {
	assert := _assert.New(t)

	t.Run("should consume events immediately if delay is not defined", func(t *testing.T) {
		cfg := ManagerConfig{}
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(handled.withName(EventBird, "bob"))

		mng.Dispatch(Event{Id: EventBird, Payload: "eagle"})

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
		mng.Subscribe(handled.withName(EventBird, "bob"))

		mng.Dispatch(Event{Id: EventBird, Payload: "eagle"})
		assert.Len(handled.events, 0, "should not handle events before delayed manager has started")

		// let's add another EventBird handler after the actual event has been dispatched
		mng.Subscribe(handled.withName(EventBird, "john"))

		// now all events should be handled by all subscribers
		err = mng.Init()
		assert.NoError(err)

		assert.Len(handled.events, 2)
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[0])
		assert.Equal("'john' handled 'bird/eagle'", handled.events[1], "john should be still able to handle the delayed events even though he was subscribed after the event has been dispatched")
	})
}

type handler struct {
	events []string
}

func (h *handler) withName(id string, name string) subscriber {
	return h.withNameAndPriority(id, name, 0)
}

func (h *handler) withNameAndPriority(id string, name string, order float64) subscriber {
	return subscriber{
		prio: order,
		handler: func(event IEvent) IEvent {
			switch e := event.(type) {
			case Event:
				if e.Id == id {
					oEvent := event.(Event)
					e := fmt.Sprintf("'%s' handled '%s/%s'", name, oEvent.Id, oEvent.Payload)
					h.events = append(h.events, e)
				}
			}
			return event
		},
	}
}

type _ext struct{}

var extension = _ext{}

func (e _ext) withName(id string, name string) subscriber {
	return e.withNameAndPriority(id, name, 0)
}

func (_ext) withNameAndPriority(id string, name string, order float64) subscriber {
	return subscriber{
		prio: order,
		handler: func(originalEvent IEvent) (updatedEvent IEvent) {
			switch e := originalEvent.(type) {
			case Event:
				if e.Id == id {
					return Event{
						Id:      e.Id,
						Payload: fmt.Sprintf("%s/%s", name, originalEvent.(Event).Payload),
					}
				}
			}
			return originalEvent
		},
	}
}
