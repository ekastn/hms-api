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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve appointments", nil)
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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve appointment", nil)
	}

	if appointment == nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Appointment not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Appointment details", appointment.ToDTO())
}

func (h *AppointmentHandler) Create(c *fiber.Ctx) error {
	var body domain.AppointmentDTO
	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Convert DTO to Entity
	appointment, err := body.ToEntity()
	if err != nil {
		log.Printf("Error converting DTO to entity: %v", err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	// Create the appointment
	id, err := h.appointmentService.Create(c.Context(), &appointment)
	if err != nil {
		log.Printf("Error creating appointment: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	createdAppointment, err := h.appointmentService.GetByID(c.Context(), id)
	if err != nil {
		log.Printf("Error fetching created appointment: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Appointment created but failed to fetch details", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Appointment created successfully", createdAppointment.ToDTO())
}

func (h *AppointmentHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Appointment ID is required", nil)
	}

	var body domain.AppointmentDTO
	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	appointment, err := body.ToEntity()
	if err != nil {
		log.Printf("Error converting DTO to entity: %v", err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	err = h.appointmentService.Update(c.Context(), id, &appointment)
	if err != nil {
		log.Printf("Error updating appointment: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	updatedAppointment, err := h.appointmentService.GetByID(c.Context(), id)
	if err != nil {
		log.Printf("Error fetching updated appointment: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Appointment updated but failed to fetch details", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Appointment updated successfully", updatedAppointment.ToDTO())
}

func (h *AppointmentHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Appointment ID is required", nil)
	}

	err := h.appointmentService.Delete(c.Context(), id)
	if err != nil {
		log.Printf("Error deleting appointment %s: %v", id, err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Appointment cancelled successfully", nil)
}
