package middleware

import (
	"github.com/darielgaizta/realtime-leaderboard/tools"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(m.Config.JWTSecret)},
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			c.Locals("user_id", claims["user_id"])
			c.Locals("username", claims["username"])
			c.Locals("email", claims["email"])
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return tools.RespondWith401(c, "Invalid token")
		},
	})
}
