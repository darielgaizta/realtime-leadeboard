package router

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	"github.com/darielgaizta/realtime-leaderboard/internal/handler"
	"github.com/darielgaizta/realtime-leaderboard/internal/middleware"
	"github.com/darielgaizta/realtime-leaderboard/tools"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	CSRFHandler *handler.CSRFHandler
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
}

func NewRouter(app *app.App, jwt *tools.JWT) *Router {
	return &Router{
		CSRFHandler: handler.NewCSRFHandler(app),
		UserHandler: handler.NewUserHandler(app),
		AuthHandler: handler.NewAuthHandler(app, jwt),
	}
}

func (router *Router) Install(app *fiber.App, m *middleware.Middleware) {
	api := app.Group("/api")

	// Applying middleware
	api.Use(m.CORSProtection())
	api.Use(m.CSRFProtection())

	// API Version 1.0
	v1 := api.Group("/v1")
	RegisterAuthRoutes(v1, router.AuthHandler)
	RegisterCSRFRoutes(v1, router.CSRFHandler)

	// Routes that required authentication.
	protected := v1.Group("").Use(m.AuthRequired())
	RegisterUserRoutes(protected, router.UserHandler)
}
