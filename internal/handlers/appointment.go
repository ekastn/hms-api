package handlers

import (
	"log"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

// GetAll handles the request to get all appointments.
//
//	@Summary		Get all appointments
//	@Description	Retrieve a list of all appointments.
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	utils.SuccessResponse{data=[]domain.AppointmentDTO}	"List of appointments"
//	@Failure		500	{object}	utils.ErrorResponse									"Failed to retrieve appointments"
//	@Router			/appointments [get]
func (h *AppointmentHandler) GetAll(c *fiber.Ctx) error {
	appointments, err := h.appointmentService.GetAll(c.Context())
	if err != nil {
		log.Printf("Error getting appointments: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve appointments", err.Error())
	}

	var appointmentDTOs []domain.AppointmentDTO
	for _, appt := range appointments {
		appointmentDTOs = append(appointmentDTOs, appt.ToDTO())
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "List of appointments", appointmentDTOs)
}

// GetByID handles the request to get an appointment by ID.
//
//	@Summary		Get appointment by ID
//	@Description	Retrieve a single appointment by its ID.
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string												true	"Appointment ID"
//	@Success		200	{object}	utils.SuccessResponse{data=domain.AppointmentDTO}	"Appointment retrieved successfully"
//	@Failure		400	{object}	utils.ErrorResponse									"Appointment ID is required"
//	@Failure		404	{object}	utils.ErrorResponse									"Appointment not found"
//	@Failure		500	{object}	utils.ErrorResponse									"Failed to retrieve appointment"
//	@Router			/appointments/{id} [get]
func (h *AppointmentHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Appointment ID is required", nil)
	}

	appointment, err := h.appointmentService.GetByID(c.Context(), id)
	if err != nil {
		log.Printf("Error getting appointment %s: %v", id, err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve appointment", err.Error())
	}

	if appointment == nil {
		return utils.ErrorResponseJSON(c, fiber.StatusNotFound, "Appointment not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Appointment retrieved successfully", appointment.ToDTO())
}

// GetAppointmentDetail returns detailed appointment information including patient and medical history
//
//	@Summary		Get detailed appointment information
//	@Description	Retrieve detailed information for a single appointment, including patient and medical history.
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string															true	"Appointment ID"
//	@Success		200	{object}	utils.SuccessResponse{data=domain.AppointmentDetailResponse}	"Appointment details retrieved successfully"
//	@Failure		400	{object}	utils.ErrorResponse												"Appointment ID is required"
//	@Failure		404	{object}	utils.ErrorResponse												"Appointment not found"
//	@Failure		500	{object}	utils.ErrorResponse												"Failed to retrieve appointment details"
//	@Router			/appointments/{id}/detail [get]
func (h *AppointmentHandler) GetAppointmentDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Appointment ID is required", nil)
	}

	detail, err := h.appointmentService.GetAppointmentDetail(c.Context(), id)
	if err != nil {
		log.Printf("Error getting appointment detail %s: %v", id, err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve appointment details", err.Error())
	}

	if detail == nil {
		return utils.ErrorResponseJSON(c, fiber.StatusNotFound, "Appointment not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Appointment details retrieved successfully", detail)
}

// Create handles the request to create a new appointment.
//
//	@Summary		Create a new appointment
//	@Description	Create a new appointment.
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			appointment	body		domain.AppointmentDTO							true	"Appointment object to be created"
//	@Success		201			{object}	utils.SuccessResponse{data=object{id=string}}	"Appointment created successfully"
//	@Failure		400			{object}	utils.ErrorResponse								"Invalid request body or validation failed"
//	@Failure		500			{object}	utils.ErrorResponse								"Failed to create appointment"
//	@Router			/appointments [post]
func (h *AppointmentHandler) Create(c *fiber.Ctx) error {
	var body domain.AppointmentDTO
	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	validationErrors := utils.ValidateStruct(body)
	if validationErrors != nil { // Check if there are any validation errors
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	creatorID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Invalid user ID", nil)
	}

	// Convert DTO to Entity
	appointment, err := body.ToEntity()
	if err != nil {
		log.Printf("Error converting DTO to entity: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	// Create the appointment
	id, err := h.appointmentService.Create(c.Context(), &appointment, creatorID)
	if err != nil {
		log.Printf("Error creating appointment: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Appointment created successfully", fiber.Map{"id": id})
}

// Update handles the request to update an appointment.
//
//	@Summary		Update an existing appointment
//	@Description	Update details of an existing appointment.
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id			path		string					true	"Appointment ID"
//	@Param			appointment	body		domain.AppointmentDTO	true	"Appointment object with updated fields"
//	@Success		204			{object}	utils.SuccessResponse	"Appointment updated successfully"
//	@Failure		400			{object}	utils.ErrorResponse		"Invalid request body or validation failed"
//	@Failure		500			{object}	utils.ErrorResponse		"Failed to update appointment"
//	@Router			/appointments/{id} [put]
func (h *AppointmentHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Appointment ID is required", nil)
	}

	var body domain.AppointmentDTO
	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	validationErrors := utils.ValidateStruct(body)
	if validationErrors != nil { // Check if there are any validation errors
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	updaterID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Invalid user ID", nil)
	}

	appointment, err := body.ToEntity()
	if err != nil {
		log.Printf("Error converting DTO to entity: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	err = h.appointmentService.Update(c.Context(), id, &appointment, updaterID)
	if err != nil {
		log.Printf("Error updating appointment: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Appointment updated successfully", nil)
}

// Delete handles the request to cancel an appointment.
//
//	@Summary		Cancel an appointment
//	@Description	Cancel an appointment (soft delete by changing status to 'Cancelled').
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string					true	"Appointment ID"
//	@Success		204	{object}	utils.SuccessResponse	"Appointment cancelled successfully"
//	@Failure		400	{object}	utils.ErrorResponse		"Appointment ID is required"
//	@Failure		500	{object}	utils.ErrorResponse		"Failed to cancel appointment"
//	@Router			/appointments/{id} [delete]
func (h *AppointmentHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Appointment ID is required", nil)
	}

	updaterID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Invalid user ID", nil)
	}

	err = h.appointmentService.Delete(c.Context(), id, updaterID)
	if err != nil {
		log.Printf("Error deleting appointment %s: %v", id, err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Appointment cancelled successfully", nil)
}
