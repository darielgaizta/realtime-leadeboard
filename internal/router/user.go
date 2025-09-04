package router

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(router fiber.Router, handler *handler.UserHandler) {
	r := router.Group("/user")
	r.Post("/register", handler.Register)
}
