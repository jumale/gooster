package status

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
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
	*gooster.BaseModule
	cfg Config
}

func (m *Module) Init(ctx *gooster.AppContext) error {
	view := tview.NewTable()
	m.BaseModule = gooster.NewBaseModule(m.cfg.ModuleConfig, ctx, view, view.Box)

	view.SetBorder(false)
	view.SetBorders(false)
	view.SetBackgroundColor(m.cfg.Colors.Bg)

	wd := tview.NewTableCell("")
	wd.SetTextColor(m.cfg.Colors.WorkDir)
	wd.SetExpansion(2)
	wd.SetAlign(tview.AlignLeft)

	view.SetCell(0, 0, wd)

	m.Events().Subscribe(
		events.Subscriber{Id: "workdir:change_dir", Fn: func(event events.Event) { // @todo use Actions
			path, err := filepath.Abs(event.Payload.(string))
			if err != nil {
				m.Log().Error(errors.WithMessage(err, "could not obtain working directory"))
			} else {
				usr, err := user.Current()
				if err != nil {
					m.Log().Error(errors.WithMessage(err, "could not obtain user directory"))
				} else {
					path = strings.Replace(path, usr.HomeDir, "~", 1)
				}

				wd.SetText(path)
			}
		}},
	)

	branch := tview.NewTableCell("master")
	branch.SetTextColor(m.cfg.Colors.Branch)
	branch.SetExpansion(1)
	branch.SetAlign(tview.AlignCenter)
	view.SetCell(0, 1, branch)

	k8sCtx := tview.NewTableCell("some.long-context.preview.ams1.example.com")
	k8sCtx.SetTextColor(m.cfg.Colors.K8sContext)
	k8sCtx.SetExpansion(2)
	k8sCtx.SetAlign(tview.AlignRight)
	view.SetCell(0, 2, k8sCtx)

	return nil
}
