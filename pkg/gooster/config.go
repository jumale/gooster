package gooster

import (
	"github.com/gdamore/tcell"
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
	Exit tcell.Key `json:"exit"`
}
