package tools

import "github.com/gofiber/fiber/v2"

func RespondWith400(c *fiber.Ctx, err error) error {
	return c.Status(400).JSON(fiber.Map{
		"message": err.Error(),
	})
}
