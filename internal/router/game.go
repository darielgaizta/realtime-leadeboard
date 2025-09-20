package router

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterGameRoutes(router fiber.Router, handler *handler.GameHandler) {
	r := router.Group("/game")
	r.Get("", handler.GetGames)
	r.Post("", handler.CreateGame)
	r.Get("/:game_id/score", handler.GetUserScore)
	r.Post("/:game_id/score", handler.CreateUserScore)
}
