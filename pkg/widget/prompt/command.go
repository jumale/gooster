package prompt

import (
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/pkg/errors"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

type Command struct {
	Stdout io.Writer
	Stderr io.Writer
	ctx    *gooster.AppContext
}

var errNoPath = errors.New("path required")
var pathRegex = regexp.MustCompile(`^\.{1,2}/`)

func (c *Command) Run(input string) error {
	args := strings.Split(input, " ")

	if args[0] == "cd" {
		if len(args) < 2 {
			return errNoPath
		}
		c.ctx.Actions.SetWorkDir(args[1])
		return nil

	} else if pathRegex.MatchString(args[0]) {
		c.ctx.Actions.SetWorkDir(args[0])
		return nil
	}

	// Prepare the command to execute.
	cmd := exec.Command("bash", "-c", input)

	// Set the correct output device.
	cmd.Stderr = c.Stderr
	cmd.Stdout = c.Stdout

	// Execute the command and return the error.
	err := cmd.Run()
	if err != nil && !strings.HasPrefix(err.Error(), "exit status") {
		return err
	} else {
		return nil
	}
}
