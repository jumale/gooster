package app

import (
	"bytes"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/complete"
	completeExt "github.com/jumale/gooster/pkg/gooster/module/complete/ext"
	"github.com/jumale/gooster/pkg/gooster/module/output"
	"github.com/jumale/gooster/pkg/gooster/module/prompt"
	"github.com/jumale/gooster/pkg/gooster/module/status"
	statusExt "github.com/jumale/gooster/pkg/gooster/module/status/ext"
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	workdirExt "github.com/jumale/gooster/pkg/gooster/module/workdir/ext"
	"github.com/pkg/errors"
	"os"
)

func Run(cfgPath string) {
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		panic(errors.WithMessagef(err, "Failed to open config file '%s'", cfgPath))
	}

	shell, err := gooster.NewApp(cfgFile, bytes.NewBufferString(defaultConfig))
	if err != nil {
		panic(err)
	}

	shell.RegisterModule(
		workdir.NewModule(),
		workdirExt.NewSortTree(),
		workdirExt.NewTypingSearch(),
	)
	shell.RegisterModule(
		output.NewModule(),
	)
	shell.RegisterModule(
		prompt.NewModule(),
	)
	shell.RegisterModule(
		status.NewModule(),
		statusExt.NewWorkDir(),
	)
	shell.RegisterModule(
		complete.NewModule(),
		completeExt.NewBashCompletion(),
	)

	shell.Run()
}
