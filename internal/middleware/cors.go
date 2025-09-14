package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Reference: https://docs.gofiber.io/api/middleware/cors/
func (m *Middleware) CORSProtection() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     m.Config.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-CSRF-Token",
	})
}
