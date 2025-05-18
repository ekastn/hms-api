package utils

import "github.com/gofiber/fiber/v2"

func ResponseJSON(c *fiber.Ctx, status int, msg string, data interface{}) error {
	var success bool
	if status >= 200 && status < 400 {
		success = true
	} else {
		success = false
	}

	return c.Status(status).JSON(fiber.Map{
		"data":    data,
		"success": success,
		"message": msg,
	})
}
