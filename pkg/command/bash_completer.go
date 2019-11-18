package command

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	CompAlias        rune = 'a'
	CompBuiltin           = 'b'
	CompCmd               = 'c'
	CompDir               = 'd'
	CompExportedVars      = 'e'
	CompFileAndDir        = 'f'
	//CompGroups             = 'g'
	//CompJobs               = 'j'
	//CompReservedWords      = 'k'
	//CompServices           = 's'
	//CompUsers              = 'u'
	//CompShellVars          = 'v'
)

type BashCompleterConfig struct {
	CompleteBin string
	CompgenBin  string
}

type BashCompleter struct {
	cfg BashCompleterConfig
}

func NewBashCompleter(cfg BashCompleterConfig) *BashCompleter {
	if cfg.CompleteBin == "" {
		cfg.CompleteBin = "complete"
	}
	if cfg.CompgenBin == "" {
		cfg.CompgenBin = "compgen"
	}
	return &BashCompleter{cfg: cfg}
}

func (b *BashCompleter) Get(cmd Definition) (Completions, error) {
	if len(cmd.Args) == 0 {
		return b.completeCommand(cmd.Command)
	} else {
		return b.completeArg(cmd.Command, cmd.Args)
	}
}

func (b *BashCompleter) completeCommand(cmd string) (Completions, error) {
	return b.compgen(cmd, CompAlias, CompBuiltin, CompCmd)
}

func (b *BashCompleter) completeArg(cmd string, args []string) (Completions, error) {
	var arg string
	arg, _ = shiftArg(args)

	if strings.HasPrefix(arg, "$") {
		return b.compgen(arg, CompExportedVars)
	} else {
		completer := b.findCustomCompletion(cmd)
		if completer != "" {
			result, err := b.getOutputLines(fmt.Sprintf(`%s "%s"`, completer, arg))
			return b.cleanUp(result), err
		}
	}

	switch cmd {
	case "cd":
		return b.compgen(arg, CompDir)
	default:
		return b.compgen(arg, CompFileAndDir)
	}
}

func (b *BashCompleter) compgen(arg string, flags ...rune) (Completions, error) {
	generate := fmt.Sprintf(`%s -%s "%s"`, b.cfg.CompgenBin, string(flags), arg)
	result, err := b.getOutputLines(generate)
	return b.cleanUp(result), err
}

func (b *BashCompleter) getOutput(cmd string) ([]byte, error) {
	c := exec.Command("bash", "-l", "-c", cmd)

	stdout := bytes.NewBuffer(nil)
	c.Stdout = stdout

	stderr := bytes.NewBuffer(nil)
	c.Stderr = stderr

	err := c.Run()
	if err != nil {
		wd, _ := os.Getwd()
		return nil, errors.WithMessagef(err, "Failed compgen. Work dir: %s, Stderr: %s", wd, stderr.String())
	}
	return bytes.Trim(stdout.Bytes(), "\n"), nil
}

func (b *BashCompleter) getOutputLines(cmd string) ([]string, error) {
	result, err := b.getOutput(cmd)
	return strings.Split(string(result), "\n"), err
}

func (b *BashCompleter) findCustomCompletion(cmd string) string {
	c := exec.Command("bash", "-l", "-c", fmt.Sprintf(`%s -p "%s"`, b.cfg.CompleteBin, cmd))
	stdout := bytes.NewBuffer(nil)
	c.Stdout = stdout

	if err := c.Run(); err != nil {
		return ""
	}

	result := stdout.Bytes()
	if len(result) == 0 {
		return ""
	}

	result = regexp.
		MustCompile(`^complete`).
		ReplaceAll(result, []byte(b.cfg.CompgenBin))

	result = regexp.
		MustCompile(regexp.QuoteMeta(cmd)+`\s*$`).
		ReplaceAll(result, nil)

	return string(result)
}

func (b *BashCompleter) cleanUp(elements []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if encountered[elements[v]] != true && elements[v] != "" {
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
