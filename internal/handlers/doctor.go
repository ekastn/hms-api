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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if doc == nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Doctor not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Doctor details", doc.ToDTO())
}

func (h *DoctorHandler) GetDoctorDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get doctor detail with related data
	detail, err := h.docService.GetDoctorDetail(c.Context(), id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if detail == nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Doctor not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Doctor details with recent patients", detail)
}

func (h *DoctorHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateDoctorRequet
	if err := c.BodyParser(&req); err != nil {
		log.Println(err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// Fetch the created doctor to return complete data
	createdDoc, err := h.docService.GetByID(c.Context(), id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Doctor created but failed to fetch data", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Doctor created successfully", createdDoc.ToDTO())
}

func (h *DoctorHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req domain.UpdateDoctorRequet
	if err := c.BodyParser(&req); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	// Get existing doctor to preserve some fields
	existingDoc, err := h.docService.GetByID(c.Context(), id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// Fetch the updated doctor to return complete data
	updatedDoc, err := h.docService.GetByID(c.Context(), id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Doctor updated but failed to fetch updated data", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Doctor updated successfully", updatedDoc.ToDTO())
}

func (h *DoctorHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.docService.Delete(c.Context(), id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Doctor deleted successfully", nil)
}
