package prompt

import (
	"context"
	"github.com/jumale/gooster/pkg/gooster"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Command struct {
	Cmd   string
	Async bool
	Ctx   context.Context
}

type CmdRunner struct {
	Stdout io.Writer
	Stderr io.Writer
	ctx    *gooster.AppContext
}

func (c *CmdRunner) Run(cmd Command) error {
	if cmd.Async {
		go func() {
			err := c.run(cmd.Cmd, cmd.Ctx)
			if err != nil {
				c.ctx.Log().Error(err)
			}
		}()
		return nil

	} else {
		return c.run(cmd.Cmd, cmd.Ctx)
	}
}

func (c *CmdRunner) run(input string, ctx context.Context) error {
	c.ctx.Log().DebugF("Executing command `%s`", input)

	// If it's exit command
	if input == "exit" {
		go func() {
			c.ctx.Log().Debug("Executing exit command")
			c.ctx.Actions().Exit()
		}()
		return nil
	}

	// If it looks like "cd" command:
	if path := detectWorkDirPath(input); path != "" {
		c.ctx.Log().DebugF("Detected cd path '%s' from command '%s'", path, input)
		c.ctx.Actions().SetWorkDir(path)
		return nil
	}

	// Otherwise just exec the command:
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, "bash", "-l", "-c", input)
	cmd.Stderr = c.Stderr
	cmd.Stdout = c.Stdout
	c.ctx.Log().DebugF("Starting command `%s`", input)
	err := cmd.Run()
	// Most commands would return errors like "exit status 1" (e.g. `echo "foo" | grep bar`).
	// We're not interested in those errors and don't want to flood our log with them,
	// so let's filter them out.
	if err != nil && !strings.HasPrefix(err.Error(), "exit status") {
		return err
	} else {
		return nil
	}
}

var (
	userHomeDir = os.UserHomeDir
	getWd       = os.Getwd
)

var pathRegex = regexp.MustCompile(`^(?:(?:\.{1,2}/)|(?:(?:/[^/]+)+))`)

func detectWorkDirPath(command string) (path string) {
	args := strings.Split(command, " ")

	if args[0] == "cd" && len(args) >= 2 {
		path = args[1]
	} else if pathRegex.MatchString(args[0]) {
		path = args[0]
	}

	if path == "" {
		return ""
	}

	if strings.HasPrefix(path, "~") {
		ud, _ := userHomeDir()
		path = strings.Replace(path, "~", ud, 1)
	}

	if strings.HasPrefix(path, "./") {
		path = strings.Replace(path, "./", "", 1)
	}

	if !strings.HasPrefix(path, "/") {
		wd, _ := getWd()
		path = strings.Replace(wd+"/"+path, "//", "/", -1)
	}

	return path
}
