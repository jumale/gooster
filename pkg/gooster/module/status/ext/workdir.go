package ext

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/status"
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	"github.com/pkg/errors"
	"os/user"
	"path/filepath"
	"strings"
)

type WorkDirConfig struct {
	gooster.ExtensionConfig `json:",inline"`
	Col                     int                 `json:"col"`
	Align                   int                 `json:"align"` // tview.Align* constants
	Colors                  WorkDirColorsConfig `json:"colors"`
}

type WorkDirColorsConfig struct {
	Text tcell.Color `json:"text"`
}

type WorkDir struct {
	cfg WorkDirConfig
	*gooster.AppContext
}

func NewWorkDir(cfg WorkDirConfig) gooster.Extension {
	return &WorkDir{cfg: cfg}
}

func (ext *WorkDir) Config() gooster.ExtensionConfig {
	return ext.cfg.ExtensionConfig
}

func (ext *WorkDir) Init(m gooster.Module, ctx *gooster.AppContext) error {
	ext.AppContext = ctx

	ctx.Events().Subscribe(events.HandleWithPrio(events.AfterAllOtherChanges, func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case workdir.EventChangeDir:
			ext.handleEventChangeDir(event)
		}
		return e
	}))
	return nil
}

func (ext *WorkDir) handleEventChangeDir(event workdir.EventChangeDir) {
	currPath, err := filepath.Abs(event.Path)
	if err != nil {
		ext.Log().Error(errors.WithMessage(err, "could not obtain working directory"))
		return
	}

	usr, err := user.Current()
	if err != nil {
		ext.Log().Error(errors.WithMessage(err, "could not obtain user directory"))
	} else {
		currPath = strings.Replace(currPath, usr.HomeDir, "~", 1)
	}

	if ext.cfg.Colors.Text != tcell.ColorDefault {
		currPath = fmt.Sprintf("[#%06x]%s[-]", ext.cfg.Colors.Text.Hex(), currPath)
	}

	ext.Events().Dispatch(status.EventShowInStatus{
		Value: currPath,
		Col:   ext.cfg.Col,
		Align: ext.cfg.Align,
	})
}
