package workdir

import (
	"github.com/jumale/gooster/pkg/config"
)

type Config struct {
	InitDir string       `json:"init_dir"`
	Colors  ColorsConfig `json:"colors"`
	Keys    KeysConfig   `json:"keys"`
}

type ColorsConfig struct {
	Bg       config.Color `json:"bg"`
	Graphics config.Color `json:"graphics"`
	Folder   config.Color `json:"folder"`
	File     config.Color `json:"file"`
}

type KeysConfig struct {
	NewFile config.Key `json:"new_file"`
	NewDir  config.Key `json:"new_dir"`
	View    config.Key `json:"view"`
	Delete  config.Key `json:"delete"`
	Open    config.Key `json:"open"`
}
