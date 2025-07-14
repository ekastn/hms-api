package handlers

import (
	"log"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

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

	// Convert DTO to Entity
	appointment, err := body.ToEntity()
	if err != nil {
		log.Printf("Error converting DTO to entity: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	// Create the appointment
	id, err := h.appointmentService.Create(c.Context(), &appointment)
	if err != nil {
		log.Printf("Error creating appointment: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Appointment created successfully", fiber.Map{"id": id})
}

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

	appointment, err := body.ToEntity()
	if err != nil {
		log.Printf("Error converting DTO to entity: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	err = h.appointmentService.Update(c.Context(), id, &appointment)
	if err != nil {
		log.Printf("Error updating appointment: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Appointment updated successfully", nil)
}

func (h *AppointmentHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Appointment ID is required", nil)
	}

	err := h.appointmentService.Delete(c.Context(), id)
	if err != nil {
		log.Printf("Error deleting appointment %s: %v", id, err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Appointment cancelled successfully", nil)
}
