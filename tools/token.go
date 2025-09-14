package tools

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

type IssuedToken struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func NewJWT(issuer, secret string, accessExpire, refreshExpire int64) *JWT {
	return &JWT{
		Issuer:        issuer,
		Secret:        secret,
		AccessExpire:  accessExpire,
		RefreshExpire: refreshExpire,
	}
}

func (h *JWT) generateAccessToken(userID int32, username, email string) (string, error) {
	tokenID, err := h.generateTokenID()
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
	return token.SignedString([]byte(h.Secret))
}

func (h *JWT) generateRefreshToken() (string, error) {
	tokenID, err := h.generateTokenID()
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
	return token.SignedString([]byte(h.Secret))
}

func (h *JWT) generateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (h *JWT) IssueToken(userID int32, username, email string) (*IssuedToken, error) {
	accessToken, err := h.generateAccessToken(userID, username, email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := h.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &IssuedToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *JWT) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(h.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, isOK := token.Claims.(*JWTClaims); isOK && token.Valid {
		return claims, nil
	}

	return nil, fiber.ErrUnauthorized
}

func (h *JWT) HashToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(hash[:])
}
