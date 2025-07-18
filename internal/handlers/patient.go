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

// GetAll handles the request to get all patients.
//
//	@Summary		Get all patients
//	@Description	Retrieve a list of all registered patients.
//	@Tags			Patients
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	utils.SuccessResponse{data=[]domain.PatientDTO}	"List of patients"
//	@Failure		500	{object}	utils.ErrorResponse								"Failed to retrieve patients"
//	@Router			/patients [get]
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

// GetByID handles the request to get a patient by ID.
//
//	@Summary		Get patient by ID
//	@Description	Retrieve a single patient by their ID.
//	@Tags			Patients
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string											true	"Patient ID"
//	@Success		200	{object}	utils.SuccessResponse{data=domain.PatientDTO}	"Patient details"
//	@Failure		404	{object}	utils.ErrorResponse								"Patient not found"
//	@Failure		500	{object}	utils.ErrorResponse								"Failed to retrieve patient"
//	@Router			/patients/{id} [get]
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
//
//	@Summary		Get detailed patient information
//	@Description	Retrieve comprehensive information for a single patient, including recent appointments and medical history.
//	@Tags			Patients
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string														true	"Patient ID"
//	@Success		200	{object}	utils.SuccessResponse{data=domain.PatientDetailResponse}	"Patient details with appointments and medical history"
//	@Failure		404	{object}	utils.ErrorResponse											"Patient not found"
//	@Failure		500	{object}	utils.ErrorResponse											"Failed to retrieve patient details"
//	@Router			/patients/{id}/detail [get]
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

// Create handles the request to create a new patient.
//
//	@Summary		Create a new patient
//	@Description	Create a new patient entry.
//	@Tags			Patients
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			patient	body		domain.CreatePatientRequest						true	"Patient object to be created"
//	@Success		201		{object}	utils.SuccessResponse{data=object{id=string}}	"Patient created successfully"
//	@Failure		400		{object}	utils.ErrorResponse								"Invalid request body or validation failed"
//	@Failure		500		{object}	utils.ErrorResponse								"Failed to create patient"
//	@Router			/patients [post]
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

// Update handles the request to update a patient.
//
//	@Summary		Update an existing patient
//	@Description	Update details of an existing patient.
//	@Tags			Patients
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string						true	"Patient ID"
//	@Param			patient	body		domain.UpdatePatientRequest	true	"Patient object with updated fields"
//	@Success		204		{object}	utils.SuccessResponse		"Patient updated successfully"
//	@Failure		400		{object}	utils.ErrorResponse			"Invalid request body or validation failed"
//	@Failure		404		{object}	utils.ErrorResponse			"Patient not found"
//	@Failure		500		{object}	utils.ErrorResponse			"Failed to update patient"
//	@Router			/patients/{id} [put]
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

// Delete handles the request to delete a patient.
//
//	@Summary		Delete a patient
//	@Description	Delete a patient entry.
//	@Tags			Patients
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string					true	"Patient ID"
//	@Success		204	{object}	utils.SuccessResponse	"Patient deleted successfully"
//	@Failure		500	{object}	utils.ErrorResponse		"Failed to delete patient"
//	@Router			/patients/{id} [delete]
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
