package router

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	"github.com/darielgaizta/realtime-leaderboard/internal/handler"
	"github.com/darielgaizta/realtime-leaderboard/internal/middleware"
	"github.com/darielgaizta/realtime-leaderboard/tools"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	Middleware  *middleware.Middleware
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
}

func NewRouter(app *app.App, jwt *tools.JWT) *Router {
	return &Router{
		Middleware:  middleware.NewMiddleware(app),
		UserHandler: handler.NewUserHandler(app),
		AuthHandler: handler.NewAuthHandler(app, jwt),
	}
}

func (router *Router) Install(app *fiber.App) {
	api := app.Group("/api")

	// API Version 1.0
	v1 := api.Group("/v1")
	RegisterUserRoutes(v1, router.UserHandler)
	RegisterAuthRoutes(v1, router.AuthHandler)

	v1.Use(router.Middleware.CSRFProtection())
}
