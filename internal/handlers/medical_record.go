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

type MedicalRecordHandler struct {
	recordService *service.MedicalRecordService
}

// GetAll retrieves all medical records
//
//	@Summary		Get all medical records
//	@Description	Retrieve a list of all medical records.
//	@Tags			Medical Records
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	utils.SuccessResponse{data=[]domain.MedicalRecordDTO}	"Medical records retrieved successfully"
//	@Failure		500	{object}	utils.ErrorResponse										"Failed to get medical records"
//	@Router			/records [get]
func (h *MedicalRecordHandler) GetAll(c *fiber.Ctx) error {
	records, err := h.recordService.GetAll(c.Context())
	if err != nil {
		log.Printf("Error getting medical records: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical records", nil)
	}

	// Convert entities to DTOs
	dtos := make([]domain.MedicalRecordDTO, len(records))
	for i, record := range records {
		dtos[i] = record.ToDTO()
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Medical records retrieved successfully", dtos)
}

func NewMedicalRecordHandler(recordService *service.MedicalRecordService) *MedicalRecordHandler {
	return &MedicalRecordHandler{
		recordService: recordService,
	}
}

// GetByID handles the request to get a medical record by ID.
//
//	@Summary		Get medical record by ID
//	@Description	Retrieve a single medical record by its ID.
//	@Tags			Medical Records
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string												true	"Medical Record ID"
//	@Success		200	{object}	utils.SuccessResponse{data=domain.MedicalRecordDTO}	"Medical record details"
//	@Failure		400	{object}	utils.ErrorResponse									"Record ID is required"
//	@Failure		404	{object}	utils.ErrorResponse									"Medical record not found"
//	@Failure		500	{object}	utils.ErrorResponse									"Failed to get medical record"
//	@Router			/records/{id} [get]
func (h *MedicalRecordHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Record ID is required", nil)
	}

	record, err := h.recordService.GetByID(c.Context(), id)
	if err != nil {
		log.Printf("Error getting medical record: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical record", nil)
	}

	if record == nil {
		return utils.ErrorResponseJSON(c, fiber.StatusNotFound, "Medical record not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Medical record details", record.ToDTO())
}

// GetByPatientID handles the request to get medical records by patient ID.
//
//	@Summary		Get medical records by patient ID
//	@Description	Retrieve a list of medical records for a specific patient.
//	@Tags			Medical Records
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			patientId	path		string													true	"Patient ID"
//	@Success		200			{object}	utils.SuccessResponse{data=[]domain.MedicalRecordDTO}	"Patient medical records"
//	@Failure		400			{object}	utils.ErrorResponse										"Patient ID is required"
//	@Failure		500			{object}	utils.ErrorResponse										"Failed to get medical records"
//	@Router			/records/patient/{patientId} [get]
func (h *MedicalRecordHandler) GetByPatientID(c *fiber.Ctx) error {
	patientID := c.Params("patientId")
	if patientID == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Patient ID is required", nil)
	}

	records, err := h.recordService.GetByPatientID(c.Context(), patientID)
	if err != nil {
		log.Printf("Error getting medical records: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical records", nil)
	}

	var recordDTOs []domain.MedicalRecordDTO
	for _, record := range records {
		recordDTOs = append(recordDTOs, record.ToDTO())
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Patient medical records", recordDTOs)
}

// GetByDateRange handles the request to get medical records by date range.
//
//	@Summary		Get medical records by date range
//	@Description	Retrieve a list of medical records within a specified date range.
//	@Tags			Medical Records
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			start	query		string													true	"Start date (RFC3339 format)"
//	@Param			end		query		string													true	"End date (RFC3339 format)"
//	@Success		200		{object}	utils.SuccessResponse{data=[]domain.MedicalRecordDTO}	"Medical records by date range"
//	@Failure		400		{object}	utils.ErrorResponse										"Both start and end dates are required or invalid format"
//	@Failure		500		{object}	utils.ErrorResponse										"Failed to get medical records"
//	@Router			/records/date-range [get]
func (h *MedicalRecordHandler) GetByDateRange(c *fiber.Ctx) error {
	startDateStr := c.Query("start")
	endDateStr := c.Query("end")

	if startDateStr == "" || endDateStr == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Both start and end dates are required", nil)
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid start date format. Use RFC3339 format", nil)
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid end date format. Use RFC3339 format", nil)
	}

	records, err := h.recordService.GetByDateRange(c.Context(), startDate, endDate)
	if err != nil {
		log.Printf("Error getting medical records: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical records", nil)
	}

	var recordDTOs []domain.MedicalRecordDTO
	for _, record := range records {
		recordDTOs = append(recordDTOs, record.ToDTO())
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Medical records by date range", recordDTOs)
}

// Create handles the request to create a new medical record.
//
//	@Summary		Create a new medical record
//	@Description	Create a new medical record entry.
//	@Tags			Medical Records
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			medicalRecord	body		domain.CreateMedicalRecordRequest				true	"Medical record object to be created"
//	@Success		201				{object}	utils.SuccessResponse{data=object{id=string}}	"Medical record created"
//	@Failure		400				{object}	utils.ErrorResponse								"Invalid request body or validation failed"
//	@Failure		500				{object}	utils.ErrorResponse								"Failed to create medical record"
//	@Router			/records [post]
func (h *MedicalRecordHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateMedicalRecordRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validationErrors := utils.ValidateStruct(req)
	if validationErrors != nil { // Check if there are any validation errors
		return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	// Convert string IDs to ObjectID
	patientID, err := primitive.ObjectIDFromHex(req.PatientID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid patient ID format", nil)
	}

	doctorID, err := primitive.ObjectIDFromHex(req.DoctorID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid doctor ID format", nil)
	}

	record := &domain.MedicalRecordEntity{
		PatientID:   patientID,
		DoctorID:    doctorID,
		Date:        time.Now(),
		RecordType:  domain.MedicalRecordType(req.RecordType),
		Description: req.Description,
		Diagnosis:   req.Diagnosis,
		Treatment:   req.Treatment,
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := h.recordService.Create(c.Context(), record)
	if err != nil {
		log.Printf("Error creating medical record: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to create medical record", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusCreated, "Medical record created", fiber.Map{"id": id})
}

// Update handles the request to update a medical record.
//
//	@Summary		Update an existing medical record
//	@Description	Update details of an existing medical record.
//	@Tags			Medical Records
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id				path		string								true	"Medical Record ID"
//	@Param			medicalRecord	body		domain.UpdateMedicalRecordRequest	true	"Medical record object with updated fields"
//	@Success		204				{object}	utils.SuccessResponse				"Medical record updated successfully"
//	@Failure		400				{object}	utils.ErrorResponse					"Invalid request body or validation failed"
//	@Failure		404				{object}	utils.ErrorResponse					"Medical record not found"
//	@Failure		500				{object}	utils.ErrorResponse					"Failed to update medical record"
//	@Router			/records/{id} [put]
func (h *MedicalRecordHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Record ID is required", nil)
	}

	var req domain.UpdateMedicalRecordRequest
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

	existingRecord, err := h.recordService.GetByID(c.Context(), id)
	if err != nil {
		log.Printf("Error getting medical record: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical record", nil)
	}

	if existingRecord == nil {
		return utils.ErrorResponseJSON(c, fiber.StatusNotFound, "Medical record not found", nil)
	}

	if req.RecordType != "" {
		existingRecord.RecordType = domain.MedicalRecordType(req.RecordType)
	}
	if req.Description != "" {
		existingRecord.Description = req.Description
	}
	if req.Diagnosis != "" {
		existingRecord.Diagnosis = req.Diagnosis
	}
	if req.Treatment != "" {
		existingRecord.Treatment = req.Treatment
	}
	if req.Notes != "" {
		existingRecord.Notes = req.Notes
	}

	existingRecord.UpdatedAt = time.Now()

	if err := h.recordService.Update(c.Context(), id, existingRecord, updaterID); err != nil {
		log.Printf("Error updating medical record: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to update medical record", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Medical record updated successfully", nil)
}

// Delete handles the request to delete a medical record.
//
//	@Summary		Delete a medical record
//	@Description	Delete a medical record entry.
//	@Tags			Medical Records
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string					true	"Medical Record ID"
//	@Success		200	{object}	utils.SuccessResponse	"Medical record deleted successfully"
//	@Failure		400	{object}	utils.ErrorResponse		"Record ID is required"
//	@Failure		500	{object}	utils.ErrorResponse		"Failed to delete medical record"
//	@Router			/records/{id} [delete]
func (h *MedicalRecordHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Record ID is required", nil)
	}

	updaterID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Invalid user ID", nil)
	}

	if err := h.recordService.Delete(c.Context(), id, updaterID); err != nil {
		log.Printf("Error deleting medical record: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to delete medical record", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Medical record deleted successfully", nil)
}
