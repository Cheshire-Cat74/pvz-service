package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"pvz-service/internal/handler/models"
	"strings"
)

func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Missing authorization header"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := jwt.MapClaims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) && !tkn.Valid {
				return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Invalid token expiration"})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Invalid token"})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

func CheckRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(jwt.MapClaims)

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Message: "Invalid role in token"})
		}

		for _, allowedRole := range roles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Message: "Insufficient role"})
	}
}
