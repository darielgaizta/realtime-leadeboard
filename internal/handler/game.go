package handler

import (
	"fmt"
	"strconv"

	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	db "github.com/darielgaizta/realtime-leaderboard/internal/db/generated"
	"github.com/darielgaizta/realtime-leaderboard/internal/dto"
	"github.com/darielgaizta/realtime-leaderboard/tools"
	"github.com/gofiber/fiber/v2"
)

type GameHandler struct {
	App *app.App
}

func NewGameHandler(app *app.App) *GameHandler {
	return &GameHandler{
		App: app,
	}
}

func (h *GameHandler) CreateGame(c *fiber.Ctx) error {
	var request dto.GameRequest
	if err := c.BodyParser(&request); err != nil {
		return tools.RespondWith400(c, "Invalid body request")
	}

	game, err := h.App.DB.CreateGame(c.Context(), request.Name)
	if err != nil {
		return tools.RespondWith400(c, fmt.Sprintf("Game with name %s already exists", request.Name))
	}

	return c.Status(201).JSON(dto.GameResponse{
		ID:   game.ID,
		Name: game.Name,
	})
}

func (h *GameHandler) GetGames(c *fiber.Ctx) error {
	games, err := h.App.DB.GetGames(c.Context())
	if err != nil {
		return tools.RespondWith404(c, "There is no game yet")
	}

	return c.Status(200).JSON(dto.ToGameResponses(games))
}

func (h *GameHandler) GetUserScore(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)

	// Parse URL Parameter
	gameID, err := strconv.Atoi(c.Params("game_id"))
	if err != nil {
		return tools.RespondWith400(c, "Invalid game ID")
	}

	userScore, err := h.App.DB.GetUserScoreByGame(c.Context(), db.GetUserScoreByGameParams{
		UserID: int32(userID),
		GameID: int32(gameID),
	})
	if err != nil {
		return tools.RespondWith404(c, "User has no score in this game")
	}

	return c.Status(200).JSON(dto.UserGameScoreResponse{
		UserID: userScore.UserID,
		GameID: userScore.GameID,
		Score:  userScore.Value,
	})
}

func (h *GameHandler) CreateUserScore(c *fiber.Ctx) error {
	var request dto.UserGameScoreRequest
	if err := c.BodyParser(&request); err != nil {
		return tools.RespondWith400(c, "Invalid body request")
	}

	userID := c.Locals("user_id").(float64)

	// Parse URL Parameter
	gameID, err := strconv.Atoi(c.Params("game_id"))
	if err != nil {
		return tools.RespondWith400(c, "Invalid game ID")
	}

	userScore, err := h.App.DB.CreateUserScore(c.Context(), db.CreateUserScoreParams{
		UserID: int32(userID),
		GameID: int32(gameID),
		Value:  request.Score,
	})
	if err != nil {
		return tools.RespondWith400(c, "Either the user or the game is invalid")
	}

	return c.Status(201).JSON(dto.UserGameScoreResponse{
		UserID: userScore.UserID,
		GameID: userScore.GameID,
		Score:  userScore.Value,
	})
}
