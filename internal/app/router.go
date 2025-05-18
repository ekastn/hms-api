package app

import (
	"github.com/ekastn/hms-api/internal/handlers"
	"github.com/ekastn/hms-api/internal/repository"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func (a *App) setupRoutes() {
	patientRepo := repository.NewPatientRepository(a.db.Collection("patients"))

	patientService := service.NewPatientService(patientRepo)

	patientHandler := handlers.NewPatientHandler(patientService)

	api := a.f.Group("/api")

	patients := api.Group("/patients")
	patients.Get("/", patientHandler.GetAll)
	patients.Get("/:id", patientHandler.GetByID)
	patients.Post("/", patientHandler.Create)
	patients.Put("/:id", patientHandler.Update)
	patients.Delete("/:id", patientHandler.Delete)

	// Health check route
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
