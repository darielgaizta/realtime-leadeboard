package main

import (
	"github.com/darielgaizta/realtime-leaderboard/bootstrap"
	"github.com/darielgaizta/realtime-leaderboard/internal/config"
)

func main() {
	config.LoadEnv()
	app := bootstrap.NewApplication()
	bootstrap.Start(app)
}
