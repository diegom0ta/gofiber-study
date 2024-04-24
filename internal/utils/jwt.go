package utils

import (
	"errors"
	"time"

	"github.com/diegom0ta/gofiber-study/internal/config"
	jwt "github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT token
func GenerateToken(email string) (string, error) {
	if email == "" {
		return "", errors.New("email must be non-empty strings")
	}

	secretKey := config.LoadJwtSecret()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
