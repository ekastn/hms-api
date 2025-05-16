package app

import "github.com/gofiber/fiber/v2"

func (a *App) SetupRoutes() {
	api := a.f.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
