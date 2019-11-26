package ext

import (
	"github.com/jumale/gooster/pkg/command"
	"github.com/jumale/gooster/pkg/gooster"
	tools "github.com/jumale/gooster/pkg/gooster/test_tools"
	"testing"
)

func TestExtension(t *testing.T) {
	t.Run("should add completion to the event", func(t *testing.T) {
		ext := tools.NewExtensionTester(t, NewBashCompletion(), nil, nil)
		ext.AssertInited()

		commands := []command.Definition{{Command: "comple"}}
		ext.SendEvent(gooster.EventSetCompletion{Commands: commands})

		ext.AssertFinalEvent(gooster.EventSetCompletion{
			Commands:   commands,
			Completion: []string{"complete"},
		})
	})

	t.Run("should not add/modify completion if the event already has non-empty completion", func(t *testing.T) {
		ext := tools.NewExtensionTester(t, NewBashCompletion(), nil, nil)
		ext.AssertInited()

		originalEvent := gooster.EventSetCompletion{
			Commands:   []command.Definition{{Command: "comple"}},
			Completion: []string{"bar", "baz"},
		}

		ext.SendEvent(originalEvent)
		ext.AssertFinalEvent(originalEvent)
	})

	t.Run("should complete only the latest command", func(t *testing.T) {
		ext := tools.NewExtensionTester(t, NewBashCompletion(), nil, nil)
		ext.AssertInited()

		commands := []command.Definition{
			{Command: "compge"},
			{Command: "comple"},
		}
		ext.SendEvent(gooster.EventSetCompletion{Commands: commands})

		ext.AssertFinalEvent(gooster.EventSetCompletion{
			Commands:   commands,
			Completion: []string{"complete"},
		})
	})
}
