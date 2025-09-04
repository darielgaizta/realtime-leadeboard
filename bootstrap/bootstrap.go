package bootstrap

import (
	"log"

	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	"github.com/darielgaizta/realtime-leaderboard/internal/config"
	"github.com/darielgaizta/realtime-leaderboard/internal/router"
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

	r := router.NewRouter(application)
	r.Install(application.Server)

	return application
}
