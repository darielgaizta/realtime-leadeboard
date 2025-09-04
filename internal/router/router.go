package router

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	"github.com/darielgaizta/realtime-leaderboard/internal/handler"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	UserHandler *handler.UserHandler
}

func NewRouter(app *app.App) *Router {
	return &Router{
		UserHandler: handler.NewUserHandler(app),
	}
}

func (router *Router) Install(app *fiber.App) {
	api := app.Group("/api")

	// API Version 1.0
	v1 := api.Group("/v1")
	RegisterUserRoutes(v1, router.UserHandler)
}
