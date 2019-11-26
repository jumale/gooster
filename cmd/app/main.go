package main

import (
	"github.com/jumale/gooster/pkg/gooster/app"
	"os"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	app.Run(homeDir + "/.gooster.yaml")
}
