package workdir

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	InitDir              string       `json:"init_dir"`
	Colors               ColorsConfig `json:"colors"`
	Keys                 KeysConfig   `json:"keys"`
}

type ColorsConfig struct {
	Bg       tcell.Color `json:"bg"`
	Graphics tcell.Color `json:"lines"`
	Folder   tcell.Color `json:"folder"`
	File     tcell.Color `json:"file"`
}

type KeysConfig struct {
	NewFile tcell.Key `json:"new_file"`
	NewDir  tcell.Key `json:"new_dir"`
	View    tcell.Key `json:"view"`
	Delete  tcell.Key `json:"delete"`
	Open    tcell.Key `json:"open"`
}
