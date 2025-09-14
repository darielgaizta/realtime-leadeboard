package middleware

import (
	"time"

	"github.com/darielgaizta/realtime-leaderboard/tools"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

// Reference: https://docs.gofiber.io/api/middleware/csrf/
func (m *Middleware) CSRFProtection() fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		CookieSecure:   m.App.Config.CSRFCookieSecure, // Ensure cookies are only sent over HTTPS, "fals" for development only.
		CookieHTTPOnly: true,                          // Prevent unauthorized client from accessing the cookie (defense against XSS).
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return tools.RespondWith403(c, err.Error())
		},
	})
}
