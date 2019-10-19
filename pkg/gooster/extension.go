package gooster

type Extension interface {
	// Config returns a basic module config
	// which specifies how ans where the module is displayed in the app.
	Config() ExtensionConfig

	// Init initializes the extension based on the target module
	// and the provided AppContext. It's like a constructor,
	// and it's called once, when the module is registered in the app.
	Init(Module, *AppContext) error
}

type ExtensionConfig struct {
	Enabled bool `json:"enabled"`
}
