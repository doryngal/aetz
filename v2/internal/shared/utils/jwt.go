package utils

import (
	"binai.net/v2/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateNewAccessToken(user *models.User) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(240 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
