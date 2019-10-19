package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	Colors               ColorsConfig `json:"colors"`
	PrintDivider         bool         `json:"print_divider"`
	PrintCommand         bool         `json:"print_command"`
	HistoryFile          string       `json:"history_file"`
	Keys                 KeysConfig   `json:"keys"`
}

type ColorsConfig struct {
	Bg      tcell.Color `json:"bg"`
	Label   tcell.Color `json:"label"`
	Text    tcell.Color `json:"text"`
	Divider tcell.Color `json:"divider"`
	Command tcell.Color `json:"command"`
}

type KeysConfig struct {
	HistoryNext tcell.Key `json:"history_next"`
	HistoryPrev tcell.Key `json:"history_prev"`
}
