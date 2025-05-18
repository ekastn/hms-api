package handlers

import (
	"log"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if patient == nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Patient not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Patient details", patient.ToDTO())
}

func (h *PatientHandler) Create(c *fiber.Ctx) error {
	var req domain.CreatePatientRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println(err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// Fetch the created patient to return complete data
	createdPatient, err := h.patientService.GetByID(c.Context(), id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Patient created but failed to fetch data", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Patient created successfully", createdPatient.ToDTO())
}

func (h *PatientHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req domain.UpdatePatientRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
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

	err := h.patientService.Update(c.Context(), id, &patientEntity)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// Fetch the updated patient to return complete data
	updatedPatient, err := h.patientService.GetByID(c.Context(), id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Patient updated but failed to fetch data", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Patient updated successfully", updatedPatient.ToDTO())
}

func (h *PatientHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.patientService.Delete(c.Context(), id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Patient deleted successfully", nil)
}
