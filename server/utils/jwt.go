package utils

import (
	"fmt"
	"time"

	"github.com/binduni/bun-golang-react-monorepo/server/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenExpiry  = 15 * time.Minute   // 15 minutes
	RefreshTokenExpiry = 7 * 24 * time.Hour // 7 days
)

type JWTClaims struct {
	UserID string           `json:"sub"`
	Email  string           `json:"email"`
	Role   models.UserRole  `json:"role"`
	Type   models.TokenType `json:"type"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a JWT access token
func GenerateAccessToken(user *models.User, secret string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Type:   models.TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken generates a JWT refresh token
func GenerateRefreshToken(user *models.User, secret string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Type:   models.TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractBearerToken extracts the token from "Bearer <token>" format
func ExtractBearerToken(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return authHeader[7:], nil
}
