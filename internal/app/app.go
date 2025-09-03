package app

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/darielgaizta/realtime-leaderboard/internal/config"
	db "github.com/darielgaizta/realtime-leaderboard/internal/db/generated"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	Config *config.Config
	DB     *db.Queries
	Server *fiber.App
}

func NewApplication(appConfig *config.Config) (*App, error) {
	// Establish DB connection.
	dbConnection, err := sql.Open("postgres", appConfig.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	if err = dbConnection.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping to db: %w", err)
	}
	log.Printf("Connected to database %s:%s", appConfig.DBHost, appConfig.DBPort)

	// Initialize application.
	app := &App{
		Config: appConfig,
		DB:     db.New(dbConnection),
		Server: fiber.New(),
	}
	return app, nil
}

func (app *App) Start() {
	log.Fatal(app.Server.Listen(fmt.Sprintf(":%s", app.Config.Port)))
}
