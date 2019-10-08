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

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})

		assert.Len(handled.events, 0)
	})

	t.Run("should let subscriber to handle an event", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(EventBird, handled.withName("bob"))

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})

		assert.Len(handled.events, 1)
		assert.Equal("'bob' handled 'bird/eagle'", handled.events[0])
	})

	t.Run("should let multiple subscribers to handle multiple events", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(EventBird, handled.withName("bob"))
		mng.Subscribe(EventBird, handled.withName("john"))
		mng.Subscribe(EventFish, handled.withName("eric"))

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
		mng.Subscribe(EventBird, handled.withNameAndPriority("bob", -10))
		mng.Subscribe(EventBird, handled.withNameAndPriority("eric", 9999))
		mng.Subscribe(EventBird, handled.withNameAndPriority("john", 0))

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
		mng.Subscribe(EventBird, handled.withName("bob"))
		mng.Subscribe(EventBird, handled.withName("john"))
		mng.Subscribe(EventBird, handled.withNameAndPriority("eric", 9999))

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})

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
		mng.Subscribe(EventBird, handled.withName("bob"))
		mng.Subscribe(EventFish, handled.withName("eric"))

		mng.Extend(EventBird, extension.withName("bird_ext1"))
		mng.Extend(EventBird, extension.withName("bird_ext2"))
		mng.Extend(EventFish, extension.withName("fish_ext1"))

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})
		mng.Dispatch(Event{Id: EventFish, Data: "tuna"})

		assert.Len(handled.events, 2)
		assert.Equal("'bob' handled 'bird/bird_ext2/bird_ext1/eagle'", handled.events[0])
		assert.Equal("'eric' handled 'fish/fish_ext1/tuna'", handled.events[1])
	})

	t.Run("should consider priority while applying extensions", func(t *testing.T) {
		mng, err := NewManager(cfg)
		assert.NoError(err)

		handled := &handler{}
		mng.Subscribe(EventBird, handled.withName("bob"))

		mng.Extend(EventBird, extension.withNameAndPriority("ext1", -10))
		mng.Extend(EventBird, extension.withNameAndPriority("ext2", 9999))
		mng.Extend(EventBird, extension.withNameAndPriority("ext3", 0))

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})

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
		mng.Subscribe(EventBird, handled.withName("bob"))

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
		mng.Subscribe(EventBird, handled.withName("bob"))

		mng.Dispatch(Event{Id: EventBird, Data: "eagle"})
		assert.Len(handled.events, 0, "should not handle events before delayed manager has started")

		// let's add another EventBird handler after the actual event has been dispatched
		mng.Subscribe(EventBird, handled.withName("john"))

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

func (h *handler) withName(name string) Subscriber {
	return h.withNameAndPriority(name, 0)
}

func (h *handler) withNameAndPriority(name string, priority float64) Subscriber {
	return Subscriber{
		Handle: func(event Event) {
			e := fmt.Sprintf("'%s' handled '%s/%s'", name, event.Id, event.Data)
			h.events = append(h.events, e)
		},
		Priority: priority,
	}
}

type _ext struct{}

var extension = _ext{}

func (e _ext) withName(name string) Extension {
	return e.withNameAndPriority(name, 0)
}

func (_ext) withNameAndPriority(name string, priority float64) Extension {
	return Extension{
		Extend: func(data EventPayload) (newData EventPayload) {
			return fmt.Sprintf("%s/%s", name, data)
		},
		Priority: priority,
	}
}
