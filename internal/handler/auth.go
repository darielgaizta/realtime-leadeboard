package handler

import (
	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	"github.com/darielgaizta/realtime-leaderboard/internal/dto"
	"github.com/darielgaizta/realtime-leaderboard/tools"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	App *app.App
	JWT *tools.JWT
}

func NewAuthHandler(app *app.App, jwt *tools.JWT) *AuthHandler {
	return &AuthHandler{
		App: app,
		JWT: jwt,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var request dto.UserRequest
	if err := c.BodyParser(&request); err != nil {
		tools.RespondWith400(c, "Invalid request body")
	}

	// Check credentials: email and password
	user, err := h.App.DB.GetUserByEmail(c.Context(), request.Email)
	if err != nil {
		tools.RespondWith401(c, "Invalid credentials")
	}
	if err = tools.CompareHashPassword(user.Password, request.Password); err != nil {
		tools.RespondWith401(c, "Invalid credentials")
	}

	accessToken, err := h.JWT.GeneratAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		tools.RespondWith500(c, "Failed to generate token")
	}

	refreshToken, err := h.JWT.GenerateRefreshToken()
	if err != nil {
		tools.RespondWith500(c, "Failed to generate refresh token")
	}

	// TODO store refresh token to database

	return c.Status(201).JSON(dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    h.JWT.AccessExpire,
		TokenType:    "Bearer",
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var request dto.RefreshRequest
	if err := c.BodyParser(&request); err != nil {
		tools.RespondWith400(c, "Invalid request body")
	}

	_, err := h.JWT.ValidateToken(request.RefreshToken)
	if err != nil {
		tools.RespondWith401(c, "Invalid refresh token")
	}

	// TODO Check refresh token from database

	return nil
}
