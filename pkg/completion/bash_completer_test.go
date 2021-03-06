package completion

import (
	"github.com/jumale/gooster/pkg/command"
	"github.com/stretchr/testify/require"
	"testing"
)

const completionDir = "./testdata/completion_context"

func TestBashCompleter(t *testing.T) {
	assert := require.New(t)

	completion := func(cmd string, args ...string) *completionTester {
		return &completionTester{t: t, def: command.Definition{Command: cmd, Args: args}}
	}

	// this test uses a real directory for path auto-completions
	path := func(p string) string {
		return completionDir + p
	}

	t.Run("should complete command, when there are no arguments", func(t *testing.T) {
		completion("expor").ShouldReturn(TypeCommand, "export")
	})

	t.Run("should complete arguments", func(t *testing.T) {
		t.Run("as a directory for dir-specific commands (like cd)", func(t *testing.T) {
			completion("cd", path("/fo")).ShouldReturn(TypeDir, path("/foo"))
		})

		t.Run("otherwise complete it as a file or directory", func(t *testing.T) {
			completion("someCmd", path("/fo")).
				ShouldReturn(TypeFile, path("/foo_1.txt"), path("/foo_2.log"), path("/foo"))
		})
	})

	t.Run("should remove duplications", func(t *testing.T) {
		completion("ec").ShouldReturn(TypeCommand, "echo", "ecpg")
	})

	t.Run("should apply installed completions", func(t *testing.T) {
		installedCompletion := "complete -W 'foo bar baz' someCmd"
		completer := NewBashCompleter(BashCompleterConfig{
			CompleteBin: installedCompletion + "; complete",
			CompgenBin:  installedCompletion + "; compgen",
		})

		actual, err := completer.Get(command.Definition{Command: "someCmd", Args: []string{"ba"}})
		assert.NoError(err)
		assert.Equal(Completion{Type: TypeCustom, Suggested: []string{"bar", "baz"}}, actual)
	})
}

func TestShiftArg(t *testing.T) {
	assert := require.New(t)

	t.Run("should keep shifting last non-empty arguments until list is empty", func(t *testing.T) {
		var arg string
		input := []string{"foo", "bar", "baz"}

		arg, input = shiftArg(input)
		assert.Equal("baz", arg)
		assert.Equal([]string{"foo", "bar"}, input)

		arg, input = shiftArg(input)
		assert.Equal("bar", arg)
		assert.Equal([]string{"foo"}, input)

		arg, input = shiftArg(input)
		assert.Equal("foo", arg)
		assert.Empty(input)

		arg, input = shiftArg(input)
		assert.Empty(arg)
		assert.Empty(input)
	})

	t.Run("should skip empty values", func(t *testing.T) {
		var arg string
		input := []string{"foo", "", "bar", ""}

		arg, input = shiftArg(input)
		assert.Equal("bar", arg)
		assert.Equal([]string{"foo", ""}, input)

		arg, input = shiftArg(input)
		assert.Equal("foo", arg)
		assert.Empty(input)
	})
}

type completionTester struct {
	t   *testing.T
	def command.Definition
}

func (c completionTester) ShouldReturn(t Type, completions ...string) {
	completer := NewBashCompleter(BashCompleterConfig{})
	actual, err := completer.Get(c.def)
	require.NoError(c.t, err)
	require.Equal(c.t, Completion{Type: t, Suggested: completions}, actual)
}
