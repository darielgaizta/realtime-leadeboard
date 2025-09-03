package bootstrap

import (
	"log"

	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	"github.com/darielgaizta/realtime-leaderboard/internal/config"
)

func init() {
	config.LoadEnv()
}

func SetupApplication() *app.App {
	configuration, err := config.NewConfiguration()
	if err != nil {
		log.Fatalf("Failed to setup configuration: %v", err)
	}

	application, err := app.NewApplication(configuration)
	if err != nil {
		log.Fatalf("Failed to setup application: %v", err)
	}

	return application
}
