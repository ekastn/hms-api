package app

import (
	"github.com/ekastn/hms-api/internal/handlers"
	"github.com/ekastn/hms-api/internal/repository"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func (a *App) setupRoutes() {
	// Initialize repositories
	patientRepo := repository.NewPatientRepository(a.db.Collection("patients"))
	docRepo := repository.NewDoctorRepository(a.db.Collection("doctors"))
	appointmentRepo := repository.NewAppointmentRepository(a.db.Collection("appointments"))
	medicalRecordRepo := repository.NewMedicalRecordRepository(a.db.Collection("medical_records"))

	// Initialize services
	patientService := service.NewPatientService(patientRepo)
	docService := service.NewDoctorService(docRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	medicalRecordService := service.NewMedicalRecordService(medicalRecordRepo)

	// Initialize handlers
	patientHandler := handlers.NewPatientHandler(patientService)
	docHandler := handlers.NewDoctorHandler(docService)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)
	medicalRecordHandler := handlers.NewMedicalRecordHandler(*medicalRecordService)

	api := a.f.Group("/api")

	patients := api.Group("/patients")
	patients.Get("/", patientHandler.GetAll)
	patients.Get("/:id", patientHandler.GetByID)
	patients.Post("/", patientHandler.Create)
	patients.Put("/:id", patientHandler.Update)
	patients.Delete("/:id", patientHandler.Delete)

	doctors := api.Group("/doctors")
	doctors.Get("/", docHandler.GetAll)
	doctors.Get("/:id", docHandler.GetByID)
	doctors.Post("/", docHandler.Create)
	doctors.Put("/:id", docHandler.Update)
	doctors.Delete("/:id", docHandler.Delete)

	appointments := api.Group("/appointments")
	appointments.Get("/", appointmentHandler.GetAll)
	appointments.Get("/:id", appointmentHandler.GetByID)
	appointments.Post("/", appointmentHandler.Create)
	appointments.Put("/:id", appointmentHandler.Update)
	appointments.Delete("/:id", appointmentHandler.Delete)

	records := api.Group("/records")
	records.Get("/", medicalRecordHandler.GetAll)
	records.Get("/:id", medicalRecordHandler.GetByID)
	records.Post("/", medicalRecordHandler.Create)
	records.Put("/:id", medicalRecordHandler.Update)
	records.Delete("/:id", medicalRecordHandler.Delete)

	// Health check route
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
