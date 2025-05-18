package handlers

import (
	"log"

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
	var body domain.PatientDTO
	if err := c.BodyParser(&body); err != nil {
		log.Println(err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	patientEntity := body.ToEntity()
	id, err := h.patientService.Create(c.Context(), &patientEntity)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	patientDTO := patientEntity.ToDTO()
	patientDTO.ID = id

	return utils.ResponseJSON(c, fiber.StatusCreated, "Patient created successfully", patientDTO)
}

func (h *PatientHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var body domain.PatientDTO
	if err := c.BodyParser(&body); err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	patientEntity := body.ToEntity()
	err := h.patientService.Update(c.Context(), id, &patientEntity)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Patient updated successfully", patientEntity.ToDTO())
}

func (h *PatientHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.patientService.Delete(c.Context(), id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Patient deleted successfully", nil)
}
