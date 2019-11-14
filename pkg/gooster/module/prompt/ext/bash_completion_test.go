package ext

import (
	"github.com/jumale/gooster/pkg/command"
	"github.com/jumale/gooster/pkg/gooster"
	tools "github.com/jumale/gooster/pkg/gooster/test_tools"
	_assert "github.com/stretchr/testify/assert"
	"testing"
)

const completionDir = "./testdata/completion_context"

func TestExtension(t *testing.T) {
	t.Run("should add completion to the event", func(t *testing.T) {
		ext := tools.NewExtensionTester(t, NewBashCompletion(BashCompletionConfig{}), nil)

		ext.SendEvent(gooster.EventSetCompletion{
			Commands: []command.Definition{
				{Command: "cd", Args: []string{completionDir + "/fo"}},
			},
		})

		ext.AssertFinalEvent(gooster.EventSetCompletion{
			Commands: []command.Definition{
				{Command: "cd", Args: []string{completionDir + "/fo"}},
			},
			Completion: []string{completionDir + "/foo"},
		})
	})

	t.Run("should not add/modify completion if the event already has non-empty completion", func(t *testing.T) {
		ext := tools.NewExtensionTester(t, NewBashCompletion(BashCompletionConfig{}), nil)
		originalEvent := gooster.EventSetCompletion{
			Commands: []command.Definition{
				{Command: "cd", Args: []string{completionDir + "/fo"}},
			},
			Completion: []string{"bar", "baz"},
		}

		ext.SendEvent(originalEvent)
		ext.AssertFinalEvent(originalEvent)
	})

}

func TestGetBashCompletion(t *testing.T) {
	assert := _assert.New(t)

	// this test uses a real directory for path auto-completions
	path := func(p string) string {
		return completionDir + p
	}

	cmd := func(cmd string, args ...string) command.Definition {
		return command.Definition{Command: cmd, Args: args}
	}

	complete := func(commands ...command.Definition) Completions {
		c, err := getBashCompletion(commands)
		assert.NoError(err)
		return c
	}

	comp := func(val ...string) Completions {
		return val
	}

	t.Run("should complete command, when there are no arguments", func(t *testing.T) {
		actual := complete(cmd("expor"))
		assert.Equal(comp("export"), actual)
	})

	t.Run("should complete arguments", func(t *testing.T) {
		t.Run("as a directory for dir-specific commands (like cd)", func(t *testing.T) {
			actual := complete(cmd("cd", path("/fo")))
			assert.Equal(comp(path("/foo")), actual)
		})

		t.Run("otherwise complete it as a file or directory", func(t *testing.T) {
			actual := complete(cmd("foo", path("/fo")))
			assert.Equal(comp(path("/foo_1.txt"), path("/foo_2.log"), path("/foo")), actual)
		})
	})

	t.Run("should always complete the latest command", func(t *testing.T) {
		actual := complete(cmd("ech"), cmd("alia"), cmd("expor"))
		assert.Equal("export", actual[0])
	})

	//t.Run("should remove duplications", func(t *testing.T) {
	//	actual := complete(cmd("ec"))
	//	assert.Equal(comp(path("/asdf")), actual)
	//})
}

func TestShiftArg(t *testing.T) {
	assert := _assert.New(t)

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
