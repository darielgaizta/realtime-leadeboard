package handler

import (
	"database/sql"
	"log"
	"time"

	"github.com/darielgaizta/realtime-leaderboard/internal/app"
	db "github.com/darielgaizta/realtime-leaderboard/internal/db/generated"
	"github.com/darielgaizta/realtime-leaderboard/internal/dto"
	"github.com/darielgaizta/realtime-leaderboard/tools"
	"github.com/gofiber/fiber/v2"
	"github.com/goombaio/namegenerator"
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

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var request dto.UserRequest
	if err := c.BodyParser(&request); err != nil {
		return tools.RespondWith400(c, "Invalid request body")
	}

	// Hash password
	hashedPassword, err := tools.HashPassword(request.Password)
	if err != nil {
		return tools.RespondWith500(c, "Failed to hash password")
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
		return tools.RespondWith400(c, "Email or password is not available")
	}

	return c.Status(201).JSON(dto.UserResponse{
		Email:    user.Email,
		Username: user.Username,
	})
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

	// Set a new refresh token in cookie.
	tools.SetRefreshTokenCookie(c, issuedToken.RefreshToken)

	return c.Status(201).JSON(dto.TokenResponse{
		AccessToken:  issuedToken.AccessToken,
		RefreshToken: issuedToken.RefreshToken,
		ExpiresIn:    h.JWT.AccessExpire,
		TokenType:    "Bearer",
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshTokenFromCookie := tools.GetRefreshTokenFromCookie(c)
	if refreshTokenFromCookie == "" {
		return tools.RespondWith401(c, "Refresh token not found")
	}

	// Validate refresh token and parse JWT claims.
	claims, err := h.JWT.ValidateToken(refreshTokenFromCookie)
	if err != nil || claims.Subject != "refresh-token" {
		return tools.RespondWith400(c, "Invalid refresh token")
	}

	// Get user data.
	user, err := h.App.DB.GetUserByTokenID(c.Context(), claims.ID)
	if err != nil {
		tools.ClearRefreshTokenCookie(c) // CLear refresh token from cookie it contains invalid token.
		return tools.RespondWith401(c, "Refresh token not found or has expired")
	}

	// "Refresh Token Rotation": Revoke refresh token to issue a new one.
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

	// Set a new refresh token in cookie.
	tools.SetRefreshTokenCookie(c, issuedToken.RefreshToken)

	return c.Status(201).JSON(dto.TokenResponse{
		AccessToken:  issuedToken.AccessToken,
		RefreshToken: issuedToken.RefreshToken,
		ExpiresIn:    h.JWT.AccessExpire,
		TokenType:    "Bearer",
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	refreshTokenFromCookie := tools.GetRefreshTokenFromCookie(c)

	if refreshTokenFromCookie != "" {
		// Validate refresh token and parse JWT claims.
		claims, err := h.JWT.ValidateToken(refreshTokenFromCookie)
		if err != nil || claims.Subject != "refresh-token" {
			log.Println("Invalid refresh token is used for logout")
		} else {
			// Get user data.
			user, err := h.App.DB.GetUserByTokenID(c.Context(), claims.ID)
			if err != nil {
				log.Println("Refresh token does not exist in database")
			} else {
				// Revoke all refresh token by user.
				if err := h.App.DB.RevokeRefreshTokensByUser(c.Context(), user.ID); err != nil {
					log.Printf("Failed to revoke refresh token for user %s", err.Error())
				}
			}
		}
	}

	tools.ClearRefreshTokenCookie(c)

	return c.Status(200).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
