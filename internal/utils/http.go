package utils

import "github.com/gofiber/fiber/v2"

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
