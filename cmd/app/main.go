package main

import (
	"github.com/darielgaizta/realtime-leaderboard/bootstrap"
)

func main() {
	app := bootstrap.SetupApplication()
	app.Start()
}
