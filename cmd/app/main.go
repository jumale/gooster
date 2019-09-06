package main

import (
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/log"
	"github.com/jumale/gooster/pkg/widget/output"
	"github.com/jumale/gooster/pkg/widget/status"
	"github.com/jumale/gooster/pkg/widget/wdtree"
	"os"
)

func main() {
	shell, err := gooster.NewApp(gooster.AppConfig{
		InitDir:       getWd(),
		LogLevel:      log.Debug,
		EventsLogPath: "/tmp/gooster-events.log",

		Grid: gooster.GridConfig{
			Cols: []int{20, -1},
			Rows: []int{1, -1, 5, 5},
		},
	})
	if err != nil {
		panic(err)
	}

	shell.AddWidget(wdtree.NewWidget(wdtree.Config{
		WidgetConfig: gooster.WidgetConfig{
			Position: gooster.Position{
				Col: 0, Row: 1,
				Width: 1, Height: 3,
			},
			Focused: true,
		},
	}))

	shell.AddWidget(output.NewWidget(output.Config{
		WidgetConfig: gooster.WidgetConfig{
			Position: gooster.Position{
				Col: 1, Row: 1,
				Width: 1, Height: 1,
			},
			Focused: false,
		},
	}))

	shell.AddWidget(status.NewWidget(status.Config{
		WidgetConfig: gooster.WidgetConfig{
			Position: gooster.Position{
				Col: 0, Row: 0,
				Width: 2, Height: 1,
			},
			Focused: false,
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
