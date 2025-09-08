package handler

import (
	"time"

	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	db "github.com/darielgaizta/realtime-leaderboard/internal/db/generated"
	"github.com/darielgaizta/realtime-leaderboard/internal/dto"
	"github.com/darielgaizta/realtime-leaderboard/tools"
	"github.com/gofiber/fiber/v2"
	"github.com/goombaio/namegenerator"
)

type UserHandler struct {
	App *app.App
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{App: app}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var request dto.UserRequest
	if err := c.BodyParser(&request); err != nil {
		tools.RespondWith400(c, "Invalid request body")
	}

	// Hash password
	hashedPassword, err := tools.HashPassword(request.Password)
	if err != nil {
		tools.RespondWith400(c, err.Error())
	}

	// Generate random username
	seed := time.Now().UTC().UnixNano()
	name := namegenerator.NewNameGenerator(seed).Generate()

	user, err := h.App.DB.CreateUser(c.Context(), db.CreateUserParams{
		Username: name,
		Password: hashedPassword,
		Email:    request.Email,
	})
	if err != nil {
		tools.RespondWith400(c, err.Error())
	}

	return c.Status(201).JSON(dto.UserResponse{
		Email:    user.Email,
		Username: user.Username,
	})
}
