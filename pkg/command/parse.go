package command

import (
	"github.com/pkg/errors"
)

const (
	space     rune = 32
	dQuote         = 34
	quote          = 39
	semicolon      = 59
	escape         = 92
	pipe           = 124
)

var QuoteErr = errors.New("Non closed quotes")

func ParseCommands(input string) (commands []Definition, err error) {
	// target values
	cmd := Definition{}
	var value []rune
	var quoteSign rune

	// flags
	var quoted, escaped, isCmd, isArg bool

	flushVal := func() {
		// remove quotes
		ln := len(value)
		if ln > 0 && ((value[0] == quote && value[ln-1] == quote) || (value[0] == dQuote && value[ln-1] == dQuote)) {
			value = value[1 : ln-1]
		}

		if quoted || escaped {
			err = QuoteErr
		}

		if isCmd {
			isCmd = false
			isArg = true
			cmd.Command = string(value)
		} else if len(value) > 0 {
			cmd.Args = append(cmd.Args, string(value))
		}

		// reset value holder and flags
		value = nil
		quoted = false
		escaped = false
	}

	flushCmd := func() {
		flushVal()
		commands = append(commands, cmd)

		cmd = Definition{}
		isCmd = false
		isArg = false
	}

	for _, r := range input {
		if (r == pipe || r == semicolon) && !quoted && !escaped {
			flushCmd()
			continue
		}
		if r == space && !isCmd && !isArg {
			continue
		}
		if !isCmd && !isArg {
			isCmd = true
		}
		if r == space && !quoted && !escaped {
			flushVal()
		} else {
			if r == quoteSign && quoted && !escaped {
				quoteSign = 0
				quoted = false
			} else if (r == quote || r == dQuote) && !quoted && !escaped {
				quoteSign = r
				quoted = true
			}
			escaped = r == escape
			value = append(value, r)
		}
	}

	flushCmd()

	return commands, err
}
