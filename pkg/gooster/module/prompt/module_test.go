package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/config"
	"github.com/jumale/gooster/pkg/filesys/fstub"
	tools "github.com/jumale/gooster/pkg/gooster/test_tools"
	"testing"
	"time"
)

func TestModule(t *testing.T) {
	promptLabel := ">> "
	withLabel := func(v string) string { return promptLabel + v }
	colors := ColorsConfig{
		Bg:      config.Color(tcell.ColorDefault),
		Label:   config.Color(tcell.ColorDefault),
		Text:    config.Color(tcell.ColorDefault),
		Divider: config.Color(tcell.ColorDefault),
		Command: config.Color(tcell.ColorDefault),
	}

	cfg := Config{
		Label:      promptLabel,
		FieldWidth: 0,
		Colors:     colors,
	}

	init := func(t *testing.T, cfg Config) *tools.ModuleTester {
		m := tools.NewModuleTester(t, NewModule(), cfg)
		m.SetSize(10, 1)
		m.AssertInited()
		return m
	}

	t.Run("should display the sent prompt", func(t *testing.T) {
		module := init(t, cfg)

		module.SendEvent(EventSetPrompt{Input: "foo bar"})
		module.AssertView(withLabel("foo bar"))

		t.Run("should clear the previous prompt", func(t *testing.T) {
			module.SendEvent(EventClearPrompt{})
			module.AssertView(promptLabel)
		})
	})

	t.Run("should print the command and divider if configured", func(t *testing.T) {
		cfg := Config{
			Label:        promptLabel,
			PrintCommand: true,
			PrintDivider: true,
			Colors: ColorsConfig{
				Divider: config.Color(tcell.ColorRed),
				Command: config.Color(tcell.ColorBlue),
			},
		}
		module := init(t, cfg)

		module.SendEvent(EventExecCommand{Cmd: `echo "foo"`})
		time.Sleep(100 * time.Millisecond)

		module.AssertOutputHasLines(
			"[blue]"+promptLabel+"echo \"foo\"[-]", // printed command
			"foo",              // the command output itself
			"[red]--------[-]", // printed divider
		)
	})

	t.Run("should run command with user input", func(t *testing.T) {
		module := init(t, cfg)
		module.SendEvent(EventExecCommand{Cmd: "bash ./testdata/prompt.sh"})

		time.Sleep(100 * time.Millisecond)
		module.SendEvent(EventSendUserInput{Input: "John"})

		time.Sleep(100 * time.Millisecond)
		module.SendEvent(EventSendUserInput{Input: "Doe"})

		time.Sleep(100 * time.Millisecond)
		module.AssertOutputHasLines(
			"First name: John",
			"Last name: Doe",
			"Hello John Doe!",
		)
		module.AssertHasLog("Starting command `bash ./testdata/prompt.sh`")
		module.AssertHasLog("Command finished `bash ./testdata/prompt.sh`")
	})

	cfgWithHistory := Config{
		Label:       promptLabel,
		Colors:      colors,
		FieldWidth:  10,
		HistoryFile: "/history",
		Keys: KeysConfig{
			HistoryPrev: config.NewKey(tcell.KeyUp),
			HistoryNext: config.NewKey(tcell.KeyDown),
		},
	}

	t.Run("should navigate history", func(t *testing.T) {
		module := tools.NewModuleTester(t, NewModule(), cfgWithHistory)
		module.SetSize(10, 1)
		module.Fs.Root().Add("/history", fstub.NewFile("foo", "bar", "baz"))
		module.AssertInited()

		module.Draw()
		module.SendEvent(EventSetPrompt{Input: "init"})
		module.AssertView(withLabel("init"))

		module.PressKey(cfgWithHistory.Keys.HistoryPrev.Type)
		module.AssertView(withLabel("baz"))

		module.PressKey(cfgWithHistory.Keys.HistoryPrev.Type)
		module.AssertView(withLabel("bar"))

		module.PressKey(cfgWithHistory.Keys.HistoryNext.Type)
		module.AssertView(withLabel("baz"))

		module.PressKey(cfgWithHistory.Keys.HistoryNext.Type)
		module.AssertView(withLabel("init"))
	})

	t.Run("should use configured colors", func(t *testing.T) {
		cfg := Config{
			Label:      promptLabel,
			FieldWidth: 10,
			Colors: ColorsConfig{
				Bg:    config.Color(tcell.ColorRed),
				Label: config.Color(tcell.ColorGreen),
				Text:  config.Color(tcell.ColorBlue),
			},
		}

		module := init(t, cfg)

		module.SendEvent(EventSetPrompt{Input: "foo bar"})
		module.AssertView("[green:red]>> [blue]foo bar")
	})
}
