package helper

import "github.com/jumale/gooster/pkg/gooster"

func (m *Module) handleSetCompletion(event gooster.EventSetCompletion) {
	m.Log().DebugF("Set completion: %+v", event)
}
