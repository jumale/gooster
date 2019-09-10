package prompt

import (
	"github.com/jumale/gooster/pkg/gooster"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Command struct {
	Stdout io.Writer
	Stderr io.Writer
	ctx    *gooster.AppContext
}

var pathRegex = regexp.MustCompile(`^\.{1,2}/`)

func (c *Command) Run(input string) error {
	args := strings.Split(input, " ")

	// If it looks like "cd" command:
	if path := c.detectWorkDirPath(args); path != "" {
		c.ctx.Log.DebugF("Detected cd path '%s' from args %+v", path, args)
		c.ctx.Actions.SetWorkDir(path)
		return nil
	}

	// Otherwise just exec the command:
	cmd := exec.Command("bash", "-c", input)
	cmd.Stderr = c.Stderr
	cmd.Stdout = c.Stdout
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

func (c *Command) detectWorkDirPath(args []string) (path string) {
	if args[0] == "cd" && len(args) >= 2 {
		path = args[1]
	} else if pathRegex.MatchString(args[0]) {
		path = args[0]
	}

	if strings.HasPrefix(path, "~") {
		ud, _ := os.UserHomeDir()
		path = strings.Replace(path, "~", ud, 1)
	}

	if !strings.HasPrefix(path, "/") {
		wd, _ := os.Getwd()
		path = strings.Replace(wd+"/"+path, "//", "/", -1)
	}

	return path
}
