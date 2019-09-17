package status

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	gooster.ModuleConfig `json:",inline"`
	Colors               ColorsConfig
}

type ColorsConfig struct {
	Bg         tcell.Color
	WorkDir    tcell.Color
	Branch     tcell.Color
	K8sContext tcell.Color
}

func NewModule(cfg Config) *Module {
	return &Module{cfg: cfg}
}

type Module struct {
	cfg  Config
	view *tview.Table
	*gooster.AppContext
}

func (w *Module) Name() string {
	return "status_bar"
}

func (w *Module) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
	w.AppContext = ctx

	w.view = tview.NewTable()
	w.view.SetBorder(false)
	w.view.SetBorders(false)
	w.view.SetBackgroundColor(w.cfg.Colors.Bg)

	wd := tview.NewTableCell("")
	wd.SetTextColor(w.cfg.Colors.WorkDir)
	wd.SetExpansion(2)
	wd.SetAlign(tview.AlignLeft)
	w.view.SetCell(0, 0, wd)
	w.Actions().OnWorkDirChange(func(newPath string) {
		path, err := filepath.Abs(newPath)
		if err != nil {
			w.Log().Error(errors.WithMessage(err, "could not obtain working directory"))
		} else {
			usr, err := user.Current()
			if err != nil {
				w.Log().Error(errors.WithMessage(err, "could not obtain user directory"))
			} else {
				path = strings.Replace(path, usr.HomeDir, "~", 1)
			}

			wd.SetText(path)
		}
	})

	branch := tview.NewTableCell("master")
	branch.SetTextColor(w.cfg.Colors.Branch)
	branch.SetExpansion(1)
	branch.SetAlign(tview.AlignCenter)
	w.view.SetCell(0, 1, branch)

	k8sCtx := tview.NewTableCell("some.long-context.preview.ams1.example.com")
	k8sCtx.SetTextColor(w.cfg.Colors.K8sContext)
	k8sCtx.SetExpansion(2)
	k8sCtx.SetAlign(tview.AlignRight)
	w.view.SetCell(0, 2, k8sCtx)

	return w.view, w.cfg.ModuleConfig, nil
}
