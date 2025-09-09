package bootstrap

import (
	"log"

	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	"github.com/darielgaizta/realtime-leaderboard/internal/config"
	"github.com/darielgaizta/realtime-leaderboard/internal/router"
	"github.com/darielgaizta/realtime-leaderboard/tools"
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

	jwt := tools.NewJWT(
		configuration.Name,
		configuration.JWTSecret,
		configuration.JWTAccessExpire,
		configuration.JWTRefreshExpire,
	)

	r := router.NewRouter(application, jwt)
	r.Install(application.Server)

	return application
}
