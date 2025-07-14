package handlers

import (
	"log"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type DoctorHandler struct {
	docService *service.DoctorService
}

func NewDoctorHandler(docService *service.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		docService: docService,
	}
}

func (h *DoctorHandler) GetAll(c *fiber.Ctx) error {
	doctors, err := h.docService.GetAll(c.Context())
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	var doctorDTOs []domain.DoctorDTO
	for _, doc := range doctors {
		doctorDTOs = append(doctorDTOs, doc.ToDTO())
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "List of doctors", doctorDTOs)
}

func (h *DoctorHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	doc, err := h.docService.GetByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if doc == nil {
		return utils.ErrorResponseJSON(c, fiber.StatusNotFound, "Doctor not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Doctor details", doc.ToDTO())
}

func (h *DoctorHandler) GetDoctorDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get doctor detail with related data
	detail, err := h.docService.GetDoctorDetail(c.Context(), id)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if detail == nil {
		return utils.ErrorResponseJSON(c, fiber.StatusNotFound, "Doctor not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Doctor details with recent patients", detail)
}

func (h *DoctorHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateDoctorRequet
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validationErrors := utils.ValidateStruct(req)
	if validationErrors != nil { // Check if there are any validation errors
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	// Convert request to entity
	docEntity := domain.DoctorEntity{
		Name:         req.Name,
		Specialty:    req.Specialty,
		Phone:        req.Phone,
		Email:        req.Email,
		Availability: []domain.TimeSlot{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	id, err := h.docService.Create(c.Context(), &docEntity)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Doctor created successfully", fiber.Map{"id": id})
}

func (h *DoctorHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req domain.UpdateDoctorRequet
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validationErrors := utils.ValidateStruct(req)
	if validationErrors != nil { // Check if there are any validation errors
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	// Get existing doctor to preserve some fields
	existingDoc, err := h.docService.GetByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// Update only the allowed fields
	docEntity := domain.DoctorEntity{
		Name:         req.Name,
		Specialty:    req.Specialty,
		Phone:        req.Phone,
		Email:        req.Email,
		Availability: existingDoc.Availability, // Preserve existing availability
		CreatedAt:    existingDoc.CreatedAt,    // Preserve created at
		UpdatedAt:    time.Now(),               // Update updated at
	}

	err = h.docService.Update(c.Context(), id, &docEntity)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Doctor updated successfully", nil)
}

func (h *DoctorHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.docService.Delete(c.Context(), id)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Doctor deleted successfully", nil)
}
