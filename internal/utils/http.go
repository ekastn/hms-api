package utils

import "github.com/gofiber/fiber/v2"

// ErrorResponse represents the standard error response format.
type ErrorResponse struct {
	Success bool        `json:"success" example:false`
	Message string      `json:"message" example:"Error message"`
	Errors  interface{} `json:"errors,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success" example:true`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseJSON(c *fiber.Ctx, status int, msg string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

// ErrorResponseJSON sends a JSON response for error conditions.
// It includes a 'success' flag (false), a 'message', and an 'errors' field.
func ErrorResponseJSON(c *fiber.Ctx, status int, msg string, errs interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": msg,
		"errors":  errs, // This will contain the detailed error information
	})
}
