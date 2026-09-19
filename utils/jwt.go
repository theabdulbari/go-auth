package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	return []byte(secret)
}

func accessExpiry() time.Duration {
	mins, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_MINUTES"))
	if mins <= 0 {
		mins = 15 // default 15 minutes
	}
	return time.Duration(mins) * time.Minute
}

func refreshExpiry() time.Duration {
	hours, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_HOURS"))
	if hours <= 0 {
		hours = 24 * 7 // default 7 days
	}
	return time.Duration(hours) * time.Hour
}

// --- Access Tokens (JWT) ---

func GenerateAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(accessExpiry()).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func ParseAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Guard against refresh tokens being used as access tokens
	if t, _ := claims["type"].(string); t != "access" {
		return nil, errors.New("wrong token type")
	}
	return claims, nil
}

// --- Refresh Tokens (opaque random) ---

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func RefreshTokenExpiry() time.Duration {
	return refreshExpiry()
}