package app

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dialog"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/complete"
	completeExt "github.com/jumale/gooster/pkg/gooster/module/complete/ext"
	"github.com/jumale/gooster/pkg/gooster/module/output"
	"github.com/jumale/gooster/pkg/gooster/module/prompt"
	"github.com/jumale/gooster/pkg/gooster/module/status"
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	workdirExt "github.com/jumale/gooster/pkg/gooster/module/workdir/ext"
	"github.com/jumale/gooster/pkg/log"
	"os"
	"strings"
	"time"
)

func Run() {
	args := struct {
		ShowHelp bool
		Debug    bool
		LogLevel log.Level
	}{
		ShowHelp: hasArg("-h", "--help"),
		Debug:    hasArg("-d", "--debug"),
		LogLevel: log.LevelFromString(getArg("-l", "--log")),
	}

	grid := gooster.GridConfig{
		Cols: []int{20, -1},
		Rows: []int{1, -1, 1, 5},
	}

	shell, err := gooster.NewApp(gooster.AppConfig{
		LogLevel: args.LogLevel,
		Grid:     grid,
		Keys: gooster.KeysConfig{
			Exit: tcell.KeyF12,
		},
		Dialog: dialog.Config{
			Colors: dialog.ColorsConfig{
				Bg:  tcell.ColorCornflowerBlue,
				Btn: tcell.ColorCornflowerBlue,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	shell.RegisterModule(workdir.NewModule(workdir.Config{
		//InitDir:  "/Users/yurii.maltsev/Dev/src/source.xing.com/growth/assistant",
		ModuleConfig: gooster.ModuleConfig{
			Position: gooster.Position{
				Col: 0, Row: 1,
				Width: 1, Height: 3,
			},
			Enabled:  true,
			Focused:  false,
			FocusKey: tcell.KeyCtrlW,
		},
		Colors: workdir.ColorsConfig{
			Bg:       tcell.NewHexColor(0x405454),
			Graphics: tcell.ColorLightSeaGreen,
			Folder:   tcell.ColorLightGreen,
			File:     tcell.ColorLightSteelBlue,
		},
		Keys: workdir.KeysConfig{
			NewFile: tcell.KeyF2,
			View:    tcell.KeyF3,
			NewDir:  tcell.KeyF7,
			Delete:  tcell.KeyF8,
			Open:    tcell.KeyEnter,
		},
	}), workdirExt.NewSortTree(workdirExt.SortTreeConfig{
		ExtensionConfig: gooster.ExtensionConfig{
			Enabled: true,
		},
		Mode: workdirExt.SortByType,
	}), workdirExt.NewTypingSearch(workdirExt.TypingSearchConfig{
		ExtensionConfig: gooster.ExtensionConfig{
			Enabled: true,
		},
		KeyPressInterval: 400 * time.Millisecond,
	}))

	shell.RegisterModule(output.NewModule(output.Config{
		ModuleConfig: gooster.ModuleConfig{
			Position: gooster.Position{
				Col: 1, Row: 1,
				Width: 1, Height: 1,
			},
			Enabled: true,
			Focused: false,
		},
		Colors: output.ColorsConfig{
			Bg:   tcell.NewHexColor(0x222222),
			Text: tcell.ColorDefault,
		},
	}))

	shell.RegisterModule(prompt.NewModule(prompt.Config{
		ModuleConfig: gooster.ModuleConfig{
			Position: gooster.Position{
				Col: 1, Row: 2,
				Width: 1, Height: 1,
			},
			Enabled:  true,
			Focused:  true,
			FocusKey: tcell.KeyCtrlF,
		},
		PrintDivider: true,
		PrintCommand: true,
		HistoryFile:  "~/.bash_history",
		Colors: prompt.ColorsConfig{
			Bg:      tcell.NewHexColor(0x555555),
			Label:   tcell.ColorLime,
			Text:    tcell.ColorLightGray,
			Divider: tcell.ColorLightGreen,
			Command: tcell.ColorLightSkyBlue,
		},
		Keys: prompt.KeysConfig{
			HistoryNext: tcell.KeyDown,
			HistoryPrev: tcell.KeyUp,
		},
	}), completeExt.NewBashCompletion(completeExt.BashCompletionConfig{
		ExtensionConfig: gooster.ExtensionConfig{
			Enabled: true,
		},
	}))

	shell.RegisterModule(status.NewModule(status.Config{
		ModuleConfig: gooster.ModuleConfig{
			Position: gooster.Position{
				Col: 0, Row: 0,
				Width: 2, Height: 1,
			},
			Enabled: true,
			Focused: false,
		},
		Colors: status.ColorsConfig{
			Bg:         tcell.ColorDimGray,
			WorkDir:    tcell.ColorGold,
			Branch:     tcell.ColorLimeGreen,
			K8sContext: tcell.ColorSkyblue,
		},
	}))

	shell.RegisterModule(complete.NewModule(complete.Config{
		ModuleConfig: gooster.ModuleConfig{
			Position: gooster.Position{
				Col: 1, Row: 3,
				Width: 1, Height: 1,
			},
			Enabled: true,
			Focused: false,
		},
		Colors: complete.ColorsConfig{
			Bg: tcell.NewHexColor(0x333333),
		},
		Keys: complete.KeysConfig{
			NextItem: tcell.KeyTab,
			Select:   tcell.KeyEnter,
		},
	}))

	shell.Run()
}

func hasArg(names ...string) bool {
	for _, arg := range names {
		for _, val := range os.Args {
			if val == arg {
				return true
			}
		}
	}
	return false
}

func getArg(names ...string) string {
	for _, arg := range names {
		for _, val := range os.Args {
			if strings.HasPrefix(val, arg) {
				return strings.Split(val, "=")[1]
			}
		}
	}
	return ""
}
