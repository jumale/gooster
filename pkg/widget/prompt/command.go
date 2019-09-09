package prompt

import (
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Stdout io.Writer
	Stderr io.Writer
	ctx    *gooster.AppContext
}

var errNoPath = errors.New("path required")

func (c *Command) Run(input string) error {
	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return errNoPath
		}
		//c.ctx.Actions.SetWorkDir(args[1])
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
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
