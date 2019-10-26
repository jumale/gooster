package main

import (
	"github.com/jumale/gooster/pkg/gooster/app"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start().Stop()
	app.Run()
}
