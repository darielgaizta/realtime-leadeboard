package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func NewApplication() *fiber.App {
	app := fiber.New()
	return app
}

func Start(app *fiber.App) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
