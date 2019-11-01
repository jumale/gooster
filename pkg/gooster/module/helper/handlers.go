package helper

func (m *Module) handleSetCompletion(event EventSetCompletion) {
	m.Log().DebugF("Set completion: %+v", event)
}
