package helper

import (
	"github.com/jumale/gooster/pkg/events"
)

func (m *Module) handleSetCompletion(event events.Event) {
	m.Log().DebugF("Set completion: %+v", event.Payload)
}
