package main

import (
	"fmt"
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
	}{
		ShowHelp: len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h"),
	}
	fmt.Println(os.Args, args)

	grid := gooster.GridConfig{
		Cols: []int{20, -1},
		Rows: []int{1, -1, 1, 5},
	}

	shell, err := gooster.NewApp(gooster.AppConfig{
		InitDir:  getWd(),
		LogLevel: log.Debug,
		Grid:     grid,
		//EventsLogPath: "/tmp/gooster-events.log",
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
		Colors: prompt.ColorsConfig{
			Bg:    tcell.NewHexColor(0x555555),
			Label: tcell.ColorLime,
			Text:  tcell.ColorLightGray,
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

func getWd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}
