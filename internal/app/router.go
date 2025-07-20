package app

import (
	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/handlers"
	"github.com/ekastn/hms-api/internal/repository"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func (a *App) setupRoutes() {
	// Initialize repositories
	patientRepo := repository.NewPatientRepository(a.db.Collection("patients"))
	docRepo := repository.NewDoctorRepository(a.db.Collection("doctors"))
	appointmentRepo := repository.NewAppointmentRepository(a.db.Collection("appointments"))
	medicalRecordRepo := repository.NewMedicalRecordRepository(a.db.Collection("medical_records"))
	activityRepo := repository.NewActivityRepository(a.db.Collection("activities"))
	userRepo := repository.NewUserRepository(a.db.Collection("users"))

	// Initialize services
	activityService := service.NewActivityService(activityRepo)
	patientService := service.NewPatientService(
		patientRepo,
		appointmentRepo,
		medicalRecordRepo,
		activityService,
	)
	docService := service.NewDoctorService(
		docRepo,
		appointmentRepo,
		patientRepo,
		activityService,
	)
	appointmentService := service.NewAppointmentService(
		appointmentRepo,
		patientRepo,
		medicalRecordRepo,
		activityService,
		a.db.Client(),
	)
	medicalRecordService := service.NewMedicalRecordService(medicalRecordRepo, activityService)
	dashboardService := service.NewDashboardService(
		patientRepo,
		docRepo,
		appointmentRepo,
		medicalRecordRepo,
		activityRepo,
	)
	authService := service.NewAuthService(userRepo, a.cfg.jwtSecret)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	patientHandler := handlers.NewPatientHandler(patientService)
	docHandler := handlers.NewDoctorHandler(docService)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)
	medicalRecordHandler := handlers.NewMedicalRecordHandler(medicalRecordService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	api := a.f.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// Middleware
	jwt := JWTMiddleware(a.cfg.jwtSecret)

	// User management routes (Admin only)
	users := api.Group("/users", jwt, RBACMiddleware(domain.RoleAdmin))
	users.Get("/", userHandler.HandleGetAllUsers)
	users.Post("/", userHandler.HandleCreateUser)
	users.Get("/:id", userHandler.HandleGetUserByID)
	users.Put("/:id", userHandler.HandleUpdateUser)
	users.Delete("/:id", userHandler.HandleDeactivateUser)
	users.Put("/:id/password", userHandler.HandleChangePassword)

	// Dashboard routes
	dashboard := api.Group("/dashboard", jwt, RBACMiddleware(domain.RoleAdmin, domain.RoleManagement))
	dashboard.Get("/", dashboardHandler.GetDashboardData)

	patients := api.Group("/patients", jwt)
	patients.Get("/", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleReceptionist, domain.RoleManagement), patientHandler.GetAll)
	patients.Get("/:id", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleReceptionist, domain.RoleManagement), patientHandler.GetByID)
	patients.Get("/:id/detail", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleReceptionist, domain.RoleManagement), patientHandler.GetPatientDetail) // New endpoint for detailed patient info
	patients.Post("/", RBACMiddleware(domain.RoleAdmin, domain.RoleReceptionist), patientHandler.Create)
	patients.Put("/:id", RBACMiddleware(domain.RoleAdmin, domain.RoleReceptionist), patientHandler.Update)
	patients.Delete("/:id", RBACMiddleware(domain.RoleAdmin), patientHandler.Delete)

	doctors := api.Group("/doctors", jwt)
	doctors.Get("/", RBACMiddleware(domain.RoleAdmin, domain.RoleManagement), docHandler.GetAll)
	doctors.Get("/:id", RBACMiddleware(domain.RoleAdmin, domain.RoleManagement), docHandler.GetByID)
	doctors.Get("/:id/detail", RBACMiddleware(domain.RoleAdmin, domain.RoleManagement), docHandler.GetDoctorDetail) // New endpoint for detailed doctor info
	doctors.Post("/", RBACMiddleware(domain.RoleAdmin), docHandler.Create)
	doctors.Put("/:id", RBACMiddleware(domain.RoleAdmin), docHandler.Update)
	doctors.Delete("/:id", RBACMiddleware(domain.RoleAdmin), docHandler.Delete)

	appointments := api.Group("/appointments", jwt)
	appointments.Get("/", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleReceptionist, domain.RoleManagement), appointmentHandler.GetAll)
	appointments.Get("/:id", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleReceptionist, domain.RoleManagement), appointmentHandler.GetByID)
	appointments.Get("/:id/detail", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleReceptionist, domain.RoleManagement), appointmentHandler.GetAppointmentDetail) // New endpoint for detailed appointment info
	appointments.Post("/", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleReceptionist), appointmentHandler.Create)
	appointments.Put("/:id", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleReceptionist), appointmentHandler.Update)
	appointments.Delete("/:id", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleReceptionist), appointmentHandler.Delete)

	records := api.Group("/records", jwt)
	records.Get("/", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleManagement), medicalRecordHandler.GetAll)
	records.Get("/:id", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor, domain.RoleNurse, domain.RoleManagement), medicalRecordHandler.GetByID)
	records.Post("/", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor), medicalRecordHandler.Create)
	records.Put("/:id", RBACMiddleware(domain.RoleAdmin, domain.RoleDoctor), medicalRecordHandler.Update)
	records.Delete("/:id", RBACMiddleware(domain.RoleAdmin), medicalRecordHandler.Delete)

	api.Get("/docs/*", swagger.HandlerDefault)

	api.Get("/health", healthCheck)
}

// @Summary		Health check endpoint
// @Description	Checks if the server is healthy
// @Tags			Health
// @Accept			json
// @Produce		json
// @Success		200	{object}	string	"OK"
// @Router			/health [get]
func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
