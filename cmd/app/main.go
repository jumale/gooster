package main

import (
	"fmt"
	"github.com/jumale/gooster/pkg/gooster/app"
	"os"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	app.Run(homeDir + "/.gooster.yaml")

	foo := "foo"
	bar := false

	if bar == false {
		foo := "test"
	}

	fmt.Print(foo)
	
	strins.tea
}
