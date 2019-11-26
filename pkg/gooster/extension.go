package gooster

type Extension interface {
	Name() string
	Init(Module, Context) error
}

type ExtensionConfig struct {
	Enabled bool `json:"enabled"`
}

var defaultExtConfig = ExtensionConfig{
	Enabled: true,
}
