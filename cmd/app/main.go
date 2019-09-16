package main

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/log"
	"github.com/jumale/gooster/pkg/widget/help"
	"github.com/jumale/gooster/pkg/widget/output"
	"github.com/jumale/gooster/pkg/widget/prompt"
	"github.com/jumale/gooster/pkg/widget/status"
	"github.com/jumale/gooster/pkg/widget/workdir"
	"os"
)

func main() {
	args := struct {
		ShowHelp bool
		Debug    bool
	}{
		ShowHelp: hasArg("-h") || hasArg("--help"),
		Debug:    hasArg("-d") || hasArg("--debug"),
	}

	grid := gooster.GridConfig{
		Cols: []int{20, -1},
		Rows: []int{1, -1, 1, 5},
	}

	shell, err := gooster.NewApp(gooster.AppConfig{
		//InitDir:  "/Users/yurii.maltsev/Dev/src",
		LogLevel: log.Debug,
		Grid:     grid,
		//EventsLogPath: "/tmp/gooster-events.log",
		Debug: args.Debug,
	})
	if err != nil {
		panic(err)
	}

	shell.AddWidget(help.NewWidget(help.Config{
		WidgetConfig: gooster.WidgetConfig{
			Position: gooster.Position{
				Col: 0, Row: 0,
				Width: len(grid.Cols), Height: len(grid.Rows),
			},
			Enabled: args.ShowHelp,
			Focused: false,
		},
	}))

	shell.AddWidget(workdir.NewWidget(workdir.Config{
		WidgetConfig: gooster.WidgetConfig{
			Position: gooster.Position{
				Col: 0, Row: 1,
				Width: 1, Height: 3,
			},
			Enabled:  !args.ShowHelp,
			Focused:  false,
			FocusKey: tcell.KeyCtrlW,
		},
		Colors: workdir.ColorsConfig{
			Bg:     tcell.NewHexColor(0x405454),
			Lines:  tcell.ColorLightSeaGreen,
			Folder: tcell.ColorLightGreen,
			File:   tcell.ColorLightSteelBlue,
		},
		Keys: workdir.KeysConfig{
			ViewFile: tcell.KeyF3,
			Delete:   tcell.KeyBackspace2,
			Open:     tcell.KeyEnter,
		},
	}))

	shell.AddWidget(output.NewWidget(output.Config{
		WidgetConfig: gooster.WidgetConfig{
			Position: gooster.Position{
				Col: 1, Row: 1,
				Width: 1, Height: 1,
			},
			Enabled: !args.ShowHelp,
			Focused: false,
		},
		Colors: output.ColorsConfig{
			Bg:   tcell.NewHexColor(0x222222),
			Text: tcell.ColorDefault,
		},
	}))

	shell.AddWidget(prompt.NewWidget(prompt.Config{
		WidgetConfig: gooster.WidgetConfig{
			Position: gooster.Position{
				Col: 1, Row: 2,
				Width: 1, Height: 1,
			},
			Enabled:  !args.ShowHelp,
			Focused:  true,
			FocusKey: tcell.KeyCtrlF,
		},
		PrintDivider: true,
		PrintCommand: true,
		Colors: prompt.ColorsConfig{
			Bg:      tcell.NewHexColor(0x555555),
			Label:   tcell.ColorLime,
			Text:    tcell.ColorLightGray,
			Divider: tcell.ColorLimeGreen,
			Command: tcell.ColorRoyalBlue,
		},
	}))

	shell.AddWidget(status.NewWidget(status.Config{
		WidgetConfig: gooster.WidgetConfig{
			Position: gooster.Position{
				Col: 0, Row: 0,
				Width: 2, Height: 1,
			},
			Enabled: !args.ShowHelp,
			Focused: false,
		},
		Colors: status.ColorsConfig{
			Bg:         tcell.ColorDimGray,
			WorkDir:    tcell.ColorGold,
			Branch:     tcell.ColorLimeGreen,
			K8sContext: tcell.ColorSkyblue,
		},
	}))

	shell.Run()
}

func hasArg(arg string) bool {
	for _, val := range os.Args {
		if val == arg {
			return true
		}
	}
	return false
}
