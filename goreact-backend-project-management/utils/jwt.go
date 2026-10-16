package utils

import (
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int64, role, email string, publicID uuid.UUID) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	duration, _ := time.ParseDuration(os.Getenv("JWT_EXPIRED"))

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"pub_id":  publicID,
		"email":   email,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID int64) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	duration, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED"))

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
