
package handlers

import (
	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related requests.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// Login handles the user login request.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	resp, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Login successful", resp)
}
