package handler

import (
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

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var request dto.AuthRequest

	if err := c.BodyParser(&request); err != nil {
		tools.RespondWith400(c, err)
	}

	// Hash password
	hashedPassword, err := tools.HashPassword(request.Password)
	if err != nil {
		tools.RespondWith400(c, err)
	}

	user, err := h.App.DB.CreateUser(c.Context(), db.CreateUserParams{
		Username: request.Email,
		Password: hashedPassword,
		Email:    request.Email,
	})
	if err != nil {
		tools.RespondWith400(c, err)
	}

	return c.Status(201).JSON(dto.AuthResponse{
		Email:    user.Email,
		Username: user.Username,
	})
}
