package prompt

import (
	"context"
	"io"
	"os/exec"
)

type Command struct {
	cmd      string
	runner   *exec.Cmd
	cancel   context.CancelFunc
	input    io.WriteCloser
	lastChar byte
}

func NewCommand(cmd string) *Command {
	ctx, cancel := context.WithCancel(context.Background())
	c := exec.CommandContext(ctx, "bash", "-l", "-c", cmd)
	input, _ := c.StdinPipe()

	return &Command{cmd: cmd, runner: c, cancel: cancel, input: input}
}

func (c *Command) Command() string {
	return c.cmd
}

func (c *Command) SetOutput(w io.Writer) *Command {
	w = &writerHook{
		target: w,
		hook: func(p []byte) {
			c.lastChar = p[len(p)-1]
		},
	}
	c.runner.Stdout = w
	c.runner.Stderr = w
	return c
}

func (c *Command) Run() error {
	return c.runner.Run()
}

func (c *Command) Cancel() {
	_ = c.input.Close()
	c.cancel()
}

func (c *Command) Write(p []byte) (n int, err error) {
	return c.input.Write(p)
}

func (c *Command) LastChar() byte {
	return c.lastChar
}

type writerHook struct {
	target io.Writer
	hook   func(p []byte)
}

func (l *writerHook) Write(p []byte) (n int, err error) {
	l.hook(p)
	return l.target.Write(p)
}
