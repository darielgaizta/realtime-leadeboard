package router

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(router fiber.Router, handler *handler.AuthHandler) {
	r := router.Group("/auth")
	r.Post("/login", handler.Login)
	r.Post("/refresh", handler.RefreshToken)
}
