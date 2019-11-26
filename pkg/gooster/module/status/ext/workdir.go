package ext

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/status"
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"os/user"
	"path/filepath"
	"strings"
)

type WorkDirConfig struct {
	Col    int                 `json:"col"`
	Align  int                 `json:"align"` // tview.Align* constants
	Colors WorkDirColorsConfig `json:"colors"`
}

type WorkDirColorsConfig struct {
	Text config.Color `json:"text"`
}

type WorkDir struct {
	gooster.Context
	cfg WorkDirConfig
}

func NewWorkDir() gooster.Extension {
	return &WorkDir{cfg: WorkDirConfig{
		Col:   0,
		Align: tview.AlignLeft,
		Colors: WorkDirColorsConfig{
			Text: config.Color(tcell.ColorGold),
		},
	}}
}

func (ext *WorkDir) Name() string {
	return "workdir"
}

func (ext *WorkDir) Init(_ gooster.Module, ctx gooster.Context) error {
	ext.Context = ctx
	if err := ctx.LoadConfig(&ext.cfg); err != nil {
		return err
	}

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

	if ext.cfg.Colors.Text.Origin() != tcell.ColorDefault {
		currPath = fmt.Sprintf("[#%06x]%s[-]", ext.cfg.Colors.Text.Origin().Hex(), currPath)
	}

	ext.Events().Dispatch(status.EventShowInStatus{
		Value: currPath,
		Col:   ext.cfg.Col,
		Align: ext.cfg.Align,
	})
}
