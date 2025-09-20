package handler

import (
	"fmt"

	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	db "github.com/darielgaizta/realtime-leaderboard/internal/db/generated"
	"github.com/darielgaizta/realtime-leaderboard/internal/dto"
	"github.com/darielgaizta/realtime-leaderboard/tools"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	App *app.App
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{App: app}
}

func (h *UserHandler) Hello(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)
	username := c.Locals("username").(string)
	email := c.Locals("email").(string)
	return c.Status(200).JSON(fiber.Map{
		"message": fmt.Sprintf("[%d-%s] Hello %s!", int(userID), email, username),
	})
}

func (h *UserHandler) GetUserScore(c *fiber.Ctx) error {
	var request dto.UserScoreRequest
	if err := c.BodyParser(&request); err != nil {
		return tools.RespondWith400(c, "Invalid body request")
	}

	userID := c.Locals("user_id").(float64)

	userScore, err := h.App.DB.GetUserScoreByGame(c.Context(), db.GetUserScoreByGameParams{
		UserID: int32(userID),
		GameID: request.GameID,
	})
	if err != nil {
		return tools.RespondWith404(c, "User has no score in this game")
	}

	return c.Status(200).JSON(dto.UserScoreResponse{
		UserID: userScore.UserID,
		GameID: userScore.GameID,
		Score:  userScore.Value,
	})
}
