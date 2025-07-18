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

type DoctorHandler struct {
	docService *service.DoctorService
}

func NewDoctorHandler(docService *service.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		docService: docService,
	}
}

// GetAll handles the request to get all doctors.
//
//	@Summary		Get all doctors
//	@Description	Retrieve a list of all registered doctors.
//	@Tags			Doctors
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	utils.SuccessResponse{data=[]domain.DoctorDTO}	"List of doctors"
//	@Failure		500	{object}	utils.ErrorResponse								"Failed to retrieve doctors"
//	@Router			/doctors [get]
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

// GetByID handles the request to get a doctor by ID.
//
//	@Summary		Get doctor by ID
//	@Description	Retrieve a single doctor by their ID.
//	@Tags			Doctors
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string											true	"Doctor ID"
//	@Success		200	{object}	utils.SuccessResponse{data=domain.DoctorDTO}	"Doctor details"
//	@Failure		404	{object}	utils.ErrorResponse								"Doctor not found"
//	@Failure		500	{object}	utils.ErrorResponse								"Failed to retrieve doctor"
//	@Router			/doctors/{id} [get]
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

// GetDoctorDetail returns detailed doctor information including recent patients
//
//	@Summary		Get detailed doctor information
//	@Description	Retrieve detailed information for a single doctor, including recent patients.
//	@Tags			Doctors
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string													true	"Doctor ID"
//	@Success		200	{object}	utils.SuccessResponse{data=domain.DoctorDetailResponse}	"Doctor details with recent patients"
//	@Failure		404	{object}	utils.ErrorResponse										"Doctor not found"
//	@Failure		500	{object}	utils.ErrorResponse										"Failed to retrieve doctor details"
//	@Router			/doctors/{id}/detail [get]
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

// Create handles the request to create a new doctor.
//
//	@Summary		Create a new doctor
//	@Description	Create a new doctor entry.
//	@Tags			Doctors
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			doctor	body		domain.CreateDoctorRequet						true	"Doctor object to be created"
//	@Success		201		{object}	utils.SuccessResponse{data=object{id=string}}	"Doctor created successfully"
//	@Failure		400		{object}	utils.ErrorResponse								"Invalid request body or validation failed"
//	@Failure		500		{object}	utils.ErrorResponse								"Failed to create doctor"
//	@Router			/doctors [post]
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

	creatorID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Invalid user ID", nil)
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

	id, err := h.docService.Create(c.Context(), &docEntity, creatorID)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Doctor created successfully", fiber.Map{"id": id})
}

// Update handles the request to update a doctor.
//
//	@Summary		Update an existing doctor
//	@Description	Update details of an existing doctor.
//	@Tags			Doctors
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string						true	"Doctor ID"
//	@Param			doctor	body		domain.UpdateDoctorRequet	true	"Doctor object with updated fields"
//	@Success		204		{object}	utils.SuccessResponse		"Doctor updated successfully"
//	@Failure		400		{object}	utils.ErrorResponse			"Invalid request body or validation failed"
//	@Failure		404		{object}	utils.ErrorResponse			"Doctor not found"
//	@Failure		500		{object}	utils.ErrorResponse			"Failed to update doctor"
//	@Router			/doctors/{id} [put]
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

	updaterID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Invalid user ID", nil)
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

	err = h.docService.Update(c.Context(), id, &docEntity, updaterID)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Doctor updated successfully", nil)
}

// Delete handles the request to delete a doctor.
//
//	@Summary		Delete a doctor
//	@Description	Delete a doctor entry.
//	@Tags			Doctors
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string					true	"Doctor ID"
//	@Success		204	{object}	utils.SuccessResponse	"Doctor deleted successfully"
//	@Failure		500	{object}	utils.ErrorResponse		"Failed to delete doctor"
//	@Router			/doctors/{id} [delete]
func (h *DoctorHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.docService.Delete(c.Context(), id)
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Doctor deleted successfully", nil)
}
