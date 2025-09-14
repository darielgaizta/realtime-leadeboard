package router

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterCSRFRoutes(router fiber.Router, handler *handler.CSRFHandler) {
	r := router.Group("/csrf")
	r.Get("/token", handler.GetCSRFToken)
}
