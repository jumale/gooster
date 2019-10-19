package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/log"
)

type AppConfig struct {
	Grid     GridConfig
	Keys     KeysConfig
	LogLevel log.Level
	Debug    bool
	Dialog   dialog.Config
}

type GridConfig struct {
	Cols []int
	Rows []int
}

type KeysConfig struct {
	Exit tcell.Key
}
