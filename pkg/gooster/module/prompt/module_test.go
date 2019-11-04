package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/jumale/gooster/pkg/filesys/fstub"
	tools "github.com/jumale/gooster/pkg/gooster/test_tools"
	"testing"
	"time"
)

func TestModule(t *testing.T) {
	promptLabel := ">> "
	withLabel := func(v string) string { return promptLabel + v }
	colors := ColorsConfig{
		Bg:      tcell.ColorDefault,
		Label:   tcell.ColorDefault,
		Text:    tcell.ColorDefault,
		Divider: tcell.ColorDefault,
		Command: tcell.ColorDefault,
	}

	cfg := Config{
		Label:      promptLabel,
		FieldWidth: 0,
		Colors:     colors,
	}
	fsProps := fstub.Config{
		WorkDir: "/wd",
		HomeDir: "/hd",
	}
	fs := fstub.New(fsProps)

	init := func(t *testing.T, cfg Config, fs filesys.FileSys) *tools.ModuleTester {
		m := tools.TestableModule(t, newModule(cfg, fs))
		m.SetSize(10, 1)
		return m
	}

	t.Run("should display the sent prompt", func(t *testing.T) {
		module := init(t, cfg, fs)

		module.SendEvent(EventSetPrompt{Input: "foo bar"})
		module.AssertView(withLabel("foo bar"))

		t.Run("should clear the previous prompt", func(t *testing.T) {
			module.SendEvent(EventClearPrompt{})
			module.AssertView(promptLabel)
		})
	})

	t.Run("should use a default label if not configured (or empty)", func(t *testing.T) {
		cfg := Config{
			Label:  "",
			Colors: colors,
		}
		module := init(t, cfg, fs)

		module.Draw()
		module.AssertView(" >")
	})

	t.Run("should print the command and divider if configured", func(t *testing.T) {
		cfg := Config{
			Label:        promptLabel,
			PrintCommand: true,
			PrintDivider: true,
			Colors: ColorsConfig{
				Divider: tcell.ColorRed,
				Command: tcell.ColorBlue,
			},
		}
		module := init(t, cfg, fs)

		module.SendEvent(EventExecCommand{Cmd: `echo "foo"`})
		time.Sleep(100 * time.Millisecond)

		module.AssertOutputHasLines(
			"[blue]"+promptLabel+"echo \"foo\"[-]", // printed command
			"foo",              // the command output itself
			"[red]--------[-]", // printed divider
		)
	})

	t.Run("should run command with user input", func(t *testing.T) {
		module := init(t, cfg, fs)
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

	t.Run("should navigate history", func(t *testing.T) {
		cfg := Config{
			Label:       promptLabel,
			Colors:      colors,
			FieldWidth:  10,
			HistoryFile: "/history",
			Keys: KeysConfig{
				HistoryPrev: tcell.KeyUp,
				HistoryNext: tcell.KeyDown,
			},
		}
		fs := fstub.New(fsProps)
		fs.Root().Add("/history", fstub.NewFile(
			"foo",
			"bar",
			"baz",
		))
		module := init(t, cfg, fs)

		module.Draw()
		module.AssertView(promptLabel)

		module.PressKey(cfg.Keys.HistoryPrev)
		module.AssertView(withLabel("baz"))

		module.PressKey(cfg.Keys.HistoryPrev)
		module.AssertView(withLabel("bar"))

		module.PressKey(cfg.Keys.HistoryNext)
		module.AssertView(withLabel("baz"))

		module.PressKey(cfg.Keys.HistoryNext)
		module.AssertView(promptLabel)
	})

	t.Run("should use configured colors", func(t *testing.T) {
		cfg := Config{
			Label:      promptLabel,
			FieldWidth: 10,
			Colors: ColorsConfig{
				Bg:    tcell.ColorRed,
				Label: tcell.ColorGreen,
				Text:  tcell.ColorBlue,
			},
		}

		module := init(t, cfg, fs)

		module.SendEvent(EventSetPrompt{Input: "foo bar"})
		module.AssertView("[green:red]>> [blue]foo bar")
	})
}
