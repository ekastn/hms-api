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
	patientService := service.NewPatientService(
		patientRepo,
		appointmentRepo,
		medicalRecordRepo,
	)
	docService := service.NewDoctorService(
		docRepo,
		appointmentRepo,
		patientRepo,
	)
	appointmentService := service.NewAppointmentService(
		appointmentRepo,
		patientRepo,
		medicalRecordRepo,
	)
	medicalRecordService := service.NewMedicalRecordService(medicalRecordRepo)
	dashboardService := service.NewDashboardService(
		patientRepo,
		docRepo,
		appointmentRepo,
		medicalRecordRepo,
	)

	// Initialize handlers
	patientHandler := handlers.NewPatientHandler(patientService)
	docHandler := handlers.NewDoctorHandler(docService)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)
	medicalRecordHandler := handlers.NewMedicalRecordHandler(medicalRecordService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	api := a.f.Group("/api")

	// Dashboard routes
	dashboard := api.Group("/dashboard")
	dashboard.Get("/", dashboardHandler.GetDashboardData)

	patients := api.Group("/patients")
	patients.Get("/", patientHandler.GetAll)
	patients.Get("/:id", patientHandler.GetByID)
	patients.Get("/:id/detail", patientHandler.GetPatientDetail) // New endpoint for detailed patient info
	patients.Post("/", patientHandler.Create)
	patients.Put("/:id", patientHandler.Update)
	patients.Delete("/:id", patientHandler.Delete)

	doctors := api.Group("/doctors")
	doctors.Get("/", docHandler.GetAll)
	doctors.Get("/:id", docHandler.GetByID)
	doctors.Get("/:id/detail", docHandler.GetDoctorDetail) // New endpoint for detailed doctor info
	doctors.Post("/", docHandler.Create)
	doctors.Put("/:id", docHandler.Update)
	doctors.Delete("/:id", docHandler.Delete)

	appointments := api.Group("/appointments")
	appointments.Get("/", appointmentHandler.GetAll)
	appointments.Get("/:id", appointmentHandler.GetByID)
	appointments.Get("/:id/detail", appointmentHandler.GetAppointmentDetail) // New endpoint for detailed appointment info
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
