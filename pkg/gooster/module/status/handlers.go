package status

import (
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	"github.com/pkg/errors"
	"os/user"
	"path/filepath"
	"strings"
)

func (m *Module) handleEventChangeDir(event workdir.EventChangeDir) {
	path, err := filepath.Abs(event.Path)
	if err != nil {
		m.Log().Error(errors.WithMessage(err, "could not obtain working directory"))
	} else {
		usr, err := user.Current()
		if err != nil {
			m.Log().Error(errors.WithMessage(err, "could not obtain user directory"))
		} else {
			path = strings.Replace(path, usr.HomeDir, "~", 1)
		}

		m.wd.SetText(path)
	}
}
