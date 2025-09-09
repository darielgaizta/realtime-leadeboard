package handler

import (
	"database/sql"

	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	db "github.com/darielgaizta/realtime-leaderboard/internal/db/generated"
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
		return tools.RespondWith400(c, "Invalid request body")
	}

	// Check credentials: email and password
	user, err := h.App.DB.GetUserByEmail(c.Context(), request.Email)
	if err != nil {
		return tools.RespondWith401(c, "Invalid credentials")
	}
	if err = tools.CompareHashPassword(user.Password, request.Password); err != nil {
		return tools.RespondWith401(c, "Invalid credentials")
	}

	// Issue access token and refresh token
	issuedToken, err := h.JWT.IssueToken(user.ID, user.Username, user.Email)
	if err != nil {
		return tools.RespondWith500(c, "Failed to issue JWT token")
	}

	// Store refresh token.
	refreshTokenClaims, err := h.JWT.ValidateToken(issuedToken.RefreshToken)
	if err != nil {
		return tools.RespondWith500(c, "Invalid refresh token is generated")
	}
	_, err = h.App.DB.CreateRefreshToken(c.Context(), db.CreateRefreshTokenParams{
		TokenID:    refreshTokenClaims.ID,
		UserID:     user.ID,
		TokenHash:  h.JWT.HashToken(issuedToken.RefreshToken),
		ExpiresAt:  refreshTokenClaims.ExpiresAt.Time,
		DeviceInfo: sql.NullString{String: c.Get("User-Agent"), Valid: c.Get("User-Agent") != ""},
		IpAddress:  sql.NullString{String: c.IP(), Valid: c.IP() != ""},
	})
	if err != nil {
		return tools.RespondWith500(c, "Failed to store refresh token")
	}

	return c.Status(201).JSON(dto.TokenResponse{
		AccessToken:  issuedToken.AccessToken,
		RefreshToken: issuedToken.RefreshToken,
		ExpiresIn:    h.JWT.AccessExpire,
		TokenType:    "Bearer",
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var request dto.RefreshTokenRequest
	if err := c.BodyParser(&request); err != nil {
		return tools.RespondWith400(c, "Invalid request body")
	}

	// Validate refresh token and parse JWT claims.
	claims, err := h.JWT.ValidateToken(request.RefreshToken)
	if err != nil || claims.Subject != "refresh-token" {
		return tools.RespondWith400(c, "Invalid refresh token")
	}
	user, err := h.App.DB.GetUserByTokenID(c.Context(), claims.ID)
	if err != nil {
		return tools.RespondWith401(c, "Refresh token not found or has expired")
	}

	// "Refresh Token Rotation"
	if err = h.App.DB.RevokeRefreshTokenByTokenID(c.Context(), claims.ID); err != nil {
		return tools.RespondWith500(c, "Failed to revoke refresh token")
	}

	// Reissue access token and refresh token.
	issuedToken, err := h.JWT.IssueToken(user.ID, user.Username, user.Password)
	if err != nil {
		return tools.RespondWith500(c, "Failed to issue JWT token")
	}

	// Store refresh token.
	claims, err = h.JWT.ValidateToken(issuedToken.RefreshToken)
	if err != nil {
		return tools.RespondWith500(c, "Invalid refresh token is generated")
	}
	_, err = h.App.DB.CreateRefreshToken(c.Context(), db.CreateRefreshTokenParams{
		TokenID:    claims.ID,
		UserID:     user.ID,
		TokenHash:  h.JWT.HashToken(issuedToken.RefreshToken),
		ExpiresAt:  claims.ExpiresAt.Time,
		DeviceInfo: sql.NullString{String: c.Get("User-Agent"), Valid: c.Get("User-Agent") != ""},
		IpAddress:  sql.NullString{String: c.IP(), Valid: c.IP() != ""},
	})
	if err != nil {
		return tools.RespondWith500(c, "Failed to store refresh token")
	}

	return c.Status(201).JSON(dto.TokenResponse{
		AccessToken:  issuedToken.AccessToken,
		RefreshToken: issuedToken.RefreshToken,
		ExpiresIn:    h.JWT.AccessExpire,
		TokenType:    "Bearer",
	})
}
