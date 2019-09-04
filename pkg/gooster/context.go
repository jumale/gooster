package gooster

type AppContext struct {
	EventManager *EventManager
	Logger       Logger
}

func NewAppContext(cfg AppConfig) (ctx *AppContext, err error) {
	ctx = &AppContext{}

	ctx.EventManager, err = NewEventManager(cfg.EventsLogPath)
	if err != nil {
		return nil, err
	}

	ctx.Logger = NewSelfLogger(cfg.LogLevel, ctx.EventManager)

	return ctx, nil
}
