package ext

import (
	"bytes"
	"fmt"
	"github.com/jumale/gooster/pkg/command"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
)

type Completions = []string

type BashCompletionConfig struct {
	gooster.ExtensionConfig
}

type BashCompletion struct {
	cfg BashCompletionConfig
}

func NewBashCompletion(cfg BashCompletionConfig) gooster.Extension {
	return &BashCompletion{cfg: cfg}
}

func (b *BashCompletion) Config() gooster.ExtensionConfig {
	return b.cfg.ExtensionConfig
}

func (b *BashCompletion) Init(_ gooster.Module, ctx *gooster.AppContext) error {
	ctx.Events().Subscribe(events.HandleWithPrio(10, func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case gooster.EventSetCompletion:
			if len(event.Completion) == 0 { // only if there are no completions yet
				comp, err := getBashCompletion(event.Commands)
				if err != nil {
					ctx.Log().Debug(err)
				}
				return gooster.EventSetCompletion{Commands: event.Commands, Completion: comp}
			}
		}
		return e
	}))

	return nil
}

func getBashCompletion(commands []command.Definition) (Completions, error) {
	if len(commands) == 0 {
		return nil, nil
	}
	target := commands[len(commands)-1]

	if len(target.Args) == 0 {
		return completeCommand(target.Command)
	} else {
		return completeArg(target.Command, target.Args)
	}
}

func completeCommand(cmd string) (Completions, error) {
	return compgen(cmd, CompAlias, CompBuiltin, CompCmd)
}

const (
	CompAlias         rune = 'a'
	CompBuiltin            = 'b'
	CompCmd                = 'c'
	CompDir                = 'd'
	CompExportedVars       = 'e'
	CompFileAndDir         = 'f'
	CompGroups             = 'g'
	CompJobs               = 'j'
	CompReservedWords      = 'k'
	CompServices           = 's'
	CompUsers              = 'u'
	CompShellVars          = 'v'
)

func completeArg(cmd string, args []string) (Completions, error) {
	var arg string
	arg, _ = shiftArg(args)

	if strings.HasPrefix(arg, "$") {
		return compgen(arg, CompExportedVars)
	}

	switch cmd {
	case "cd":
		return compgen(arg, CompDir)
	default:
		return compgen(arg, CompFileAndDir)
	}
}

func compgen(arg string, flags ...rune) (Completions, error) {
	generate := fmt.Sprintf(`compgen -%s "%s"`, string(flags), arg)
	c := exec.Command("bash", "-l", "-c", generate)

	stdout := bytes.NewBuffer(nil)
	c.Stdout = stdout

	stderr := bytes.NewBuffer(nil)
	c.Stderr = stderr

	err := c.Run()
	if err != nil {
		wd, _ := os.Getwd()
		return nil, errors.WithMessagef(err, "Failed compgen. Work dir: %s, Stderr: %s", wd, stderr.String())
	}
	result := strings.Trim(stdout.String(), "\n")

	return removeDuplicates(strings.Split(result, "\n")), nil
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func shiftArg(args []string) (arg string, restArgs []string) {
	for i := len(args) - 1; i >= 0; i = i - 1 {
		if args[i] != "" {
			return args[i], args[:i]
		}
	}
	return "", nil
}
