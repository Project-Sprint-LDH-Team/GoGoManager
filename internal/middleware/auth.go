package middleware

import (
	"github.com/Project-Sprint-LDH-Team/GoGoManager/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type AuthMiddleware struct {
	jwtMaker jwt.Maker
}

func NewAuthMiddleware(jwtMaker jwt.Maker) *AuthMiddleware {
	return &AuthMiddleware{
		jwtMaker: jwtMaker,
	}
}

func (m *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Check bearer scheme dengan case insensitive
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || !strings.EqualFold(headerParts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		// Get token
		token := headerParts[1]
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing token",
			})
		}

		// Verify token
		claims, err := m.jwtMaker.VerifyToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Set user ID to context
		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}
