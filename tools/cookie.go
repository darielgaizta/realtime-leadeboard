package tools

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetRefreshTokenCookie(c *fiber.Ctx, refreshToken string) {
	cookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
	}
	c.Cookie(cookie)
}

func ClearRefreshTokenCookie(c *fiber.Ctx) {
	cookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
	}
	c.Cookie(cookie)
}

func GetRefreshTokenFromCookie(c *fiber.Ctx) string {
	return c.Cookies("refresh_token")
}
