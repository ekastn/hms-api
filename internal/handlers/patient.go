package handlers

import (
	"log"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PatientHandler struct {
	patientService *service.PatientService
}

func NewPatientHandler(patientService *service.PatientService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

func (h *PatientHandler) GetAll(c *fiber.Ctx) error {
	patients, err := h.patientService.GetAll(c.Context())
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	var patientDtos []domain.PatientDTO
	for _, patient := range patients {
		patientDtos = append(patientDtos, patient.ToDTO())
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "List of patients", patientDtos)
}

func (h *PatientHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	patient, err := h.patientService.GetByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if patient == nil {
		return utils.ErrorResponseJSON(c, fiber.StatusNotFound, "Patient not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Patient details", patient.ToDTO())
}

// GetPatientDetail returns comprehensive patient information including recent appointments and medical history
func (h *PatientHandler) GetPatientDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get patient detail with related data
	detail, err := h.patientService.GetPatientDetail(c.Context(), id)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if detail == nil {
		return utils.ErrorResponseJSON(c, fiber.StatusNotFound, "Patient not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Patient details with appointments and medical history", detail)
}

func (h *PatientHandler) Create(c *fiber.Ctx) error {
	var req domain.CreatePatientRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validationErrors := utils.ValidateStruct(req)
	if validationErrors != nil { // Check if there are any validation errors
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	// Convert request to entity
	patientEntity := domain.PatientEntity{
		Name:      req.Name,
		Age:       req.Age,
		Gender:    req.Gender,
		Phone:     req.Phone,
		Email:     req.Email,
		Address:   req.Address,
		LastVisit: time.Now(),
	}

	id, err := h.patientService.Create(c.Context(), &patientEntity)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Patient created successfully", fiber.Map{"id": id})
}

func (h *PatientHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req domain.UpdatePatientRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validationErrors := utils.ValidateStruct(req)
	if validationErrors != nil { // Check if there are any validation errors
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	updaterID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Invalid user ID", nil)
	}

	// Convert request to entity
	patientEntity := domain.PatientEntity{
		Name:    req.Name,
		Age:     req.Age,
		Gender:  req.Gender,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}

	err = h.patientService.Update(c.Context(), id, &patientEntity, updaterID)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Patient updated successfully", nil)
}

func (h *PatientHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	updaterID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Invalid user ID", nil)
	}

	err = h.patientService.Delete(c.Context(), id, updaterID)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Patient deleted successfully", nil)
}
