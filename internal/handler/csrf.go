package handler

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	"github.com/darielgaizta/realtime-leaderboard/tools"
	"github.com/gofiber/fiber/v2"
)

type CSRFHandler struct {
	App *app.App
}

func NewCSRFHandler(app *app.App) *CSRFHandler {
	return &CSRFHandler{
		App: app,
	}
}

func (h *CSRFHandler) GetCSRFToken(c *fiber.Ctx) error {
	token := c.Locals("csrf")
	if token == nil {
		return tools.RespondWith403(c, "CSRF token not available")
	}
	return c.JSON(fiber.Map{
		"csrf_token": token.(string),
	})
}
