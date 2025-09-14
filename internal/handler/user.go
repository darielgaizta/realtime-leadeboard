package handler

import (
	"fmt"

	"github.com/darielgaizta/realtime-leaderboard/internal/app"
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
