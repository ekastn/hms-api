
package app

import (
	"strings"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponseJSON(c, fiber.StatusUnauthorized, "Missing or malformed JWT", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.ErrorResponseJSON(c, fiber.StatusUnauthorized, "Missing or malformed JWT", nil)
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return utils.ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid or expired JWT", nil)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid JWT claims", nil)
		}

		c.Locals("userID", claims["sub"]) // subject is the user ID
		c.Locals("userRole", claims["role"])

		return c.Next()
	}
}

func RBACMiddleware(allowedRoles ...domain.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok {
			return utils.ErrorResponseJSON(c, fiber.StatusForbidden, "Access denied", nil)
		}

		for _, allowedRole := range allowedRoles {
			if domain.Role(role) == allowedRole {
				return c.Next()
			}
		}

		return utils.ErrorResponseJSON(c, fiber.StatusForbidden, "Access denied", nil)
	}
}
