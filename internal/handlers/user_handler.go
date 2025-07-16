
package handlers

import (
	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user management requests.
type UserHandler struct {
	userService *service.UserService
	authService *service.AuthService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *service.UserService, authService *service.AuthService) *UserHandler {
	return &UserHandler{userService, authService}
}

// HandleGetAllUsers handles the request to get all users.
func (h *UserHandler) HandleGetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers(c.Context())
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve users", err.Error())
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Users retrieved successfully", users)
}

// HandleGetUserByID handles the request to get a user by ID.
func (h *UserHandler) HandleGetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve user", err.Error())
	}
	if user == nil {
		return utils.ErrorResponseJSON(c, fiber.StatusNotFound, "User not found", nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "User retrieved successfully", user)
}

// HandleCreateUser handles the request to create a new user.
func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var user domain.UserEntity
	if err := c.BodyParser(&user); err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validation can be added here

	_, err := h.authService.CreateUser(c.Context(), &user)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to create user", err.Error())
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "User created successfully", nil)
}

// HandleUpdateUser handles the request to update a user.
func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user domain.UserEntity
	if err := c.BodyParser(&user); err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.userService.UpdateUser(c.Context(), id, &user); err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to update user", err.Error())
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "User updated successfully", nil)
}

// HandleDeactivateUser handles the request to deactivate a user.
func (h *UserHandler) HandleDeactivateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userService.DeactivateUser(c.Context(), id); err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to deactivate user", err.Error())
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "User deactivated successfully", nil)
}
