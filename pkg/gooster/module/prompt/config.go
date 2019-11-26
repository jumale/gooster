package prompt

import (
	"github.com/jumale/gooster/pkg/config"
)

type Config struct {
	Label        string       `json:"label"`
	Colors       ColorsConfig `json:"colors"`
	PrintDivider bool         `json:"print_divider"`
	PrintCommand bool         `json:"print_command"`
	HistoryFile  string       `json:"history_file"`
	FieldWidth   int          `json:"field_width"`
	Keys         KeysConfig   `json:"keys"`
}

type ColorsConfig struct {
	Bg      config.Color `json:"bg"`
	Label   config.Color `json:"label"`
	Text    config.Color `json:"text"`
	Divider config.Color `json:"divider"`
	Command config.Color `json:"command"`
}

type KeysConfig struct {
	HistoryNext config.Key `json:"history_next"`
	HistoryPrev config.Key `json:"history_prev"`
}
