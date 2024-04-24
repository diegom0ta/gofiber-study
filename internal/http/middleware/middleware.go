package middleware

import (
	"fmt"
	"strings"

	"github.com/diegom0ta/gofiber-study/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuth() fiber.Handler {
	secretKey := config.LoadJwtSecret()

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Invalid Authorization header format", "data": nil})
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
		}

		c.Locals("user", token)
		return c.Next()
	}
}
