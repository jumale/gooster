package main

import (
	"github.com/jumale/gooster/pkg/gooster/app"
	"github.com/pkg/profile"
	"os"
)

func main() {
	defer profile.Start().Stop()
	homeDir, _ := os.UserHomeDir()
	app.Run(homeDir + "/.gooster.yaml")
}
