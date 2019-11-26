package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/log"
)

type AppConfig struct {
	Grid     GridConfig    `json:"grid"`
	Keys     KeysConfig    `json:"keys"`
	LogLevel log.Level     `json:"log_level"`
	Dialog   dialog.Config `json:"dialog"`
}

type GridConfig struct {
	Cols []int `json:"cols"`
	Rows []int `json:"rows"`
}

type KeysConfig struct {
	Exit config.Key `json:"exit"`
}

var defaultConfig = AppConfig{
	LogLevel: log.Info,
	Grid: GridConfig{
		Cols: []int{20, -1},
		Rows: []int{1, -1, 1, 5},
	},
	Keys: KeysConfig{
		Exit: config.Key(tcell.KeyF12),
	},
	Dialog: dialog.Config{
		Colors: dialog.ColorsConfig{
			Bg:  config.Color(tcell.ColorCornflowerBlue),
			Btn: config.Color(tcell.ColorCornflowerBlue),
		},
	},
}
