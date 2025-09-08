package tools

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	Issuer        string
	Secret        string
	AccessExpire  int64
	RefreshExpire int64
}

type JWTClaims struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func (h *JWT) GeneratAccessToken(userID int32, username, email string) (string, error) {
	tokenID, err := generateTokenID()
	if err != nil {
		return "", nil
	}

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    h.Issuer,
			Subject:   "access-token",
			ID:        tokenID,
			Audience:  []string{"your-app-audience"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.Secret)
}

func (h *JWT) GenerateRefreshToken() (string, error) {
	tokenID, err := generateTokenID()
	if err != nil {
		return "", nil
	}

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    h.Issuer,
		Subject:   "refresh-token",
		ID:        tokenID,
		Audience:  []string{"yout-app-audience"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.Secret)
}

func (h *JWT) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return h.Secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, isOK := token.Claims.(*JWTClaims); isOK && token.Valid {
		return claims, nil
	}

	return nil, fiber.ErrUnauthorized
}

func generateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
