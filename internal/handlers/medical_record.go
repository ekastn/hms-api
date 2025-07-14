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

func (h *MedicalRecordHandler) GetAll(c *fiber.Ctx) error {
	records, err := h.recordService.GetAll(c.Context())
	if err != nil {
		log.Printf("Error getting medical records: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical records", nil)
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

func (h *MedicalRecordHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Record ID is required", nil)
	}

	record, err := h.recordService.GetByID(c.Context(), id)
	if err != nil {
		log.Printf("Error getting medical record: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical record", nil)
	}

	if record == nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Medical record not found", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Medical record details", record.ToDTO())
}

func (h *MedicalRecordHandler) GetByPatientID(c *fiber.Ctx) error {
	patientID := c.Params("patientId")
	if patientID == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Patient ID is required", nil)
	}

	records, err := h.recordService.GetByPatientID(c.Context(), patientID)
	if err != nil {
		log.Printf("Error getting medical records: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical records", nil)
	}

	var recordDTOs []domain.MedicalRecordDTO
	for _, record := range records {
		recordDTOs = append(recordDTOs, record.ToDTO())
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Patient medical records", recordDTOs)
}

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
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical records", nil)
	}

	var recordDTOs []domain.MedicalRecordDTO
	for _, record := range records {
		recordDTOs = append(recordDTOs, record.ToDTO())
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Medical records by date range", recordDTOs)
}

func (h *MedicalRecordHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateMedicalRecordRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validationErrors := utils.ValidateStruct(req)
	if validationErrors != nil { // Check if there are any validation errors
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
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

func (h *MedicalRecordHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Record ID is required", nil)
	}

	var req domain.UpdateMedicalRecordRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	validationErrors := utils.ValidateStruct(req)
	if validationErrors != nil { // Check if there are any validation errors
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	existingRecord, err := h.recordService.GetByID(c.Context(), id)
	if err != nil {
		log.Printf("Error getting medical record: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to get medical record", nil)
	}

	if existingRecord == nil {
		return utils.ResponseJSON(c, fiber.StatusNotFound, "Medical record not found", nil)
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

	if err := h.recordService.Update(c.Context(), id, existingRecord); err != nil {
		log.Printf("Error updating medical record: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to update medical record", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusNoContent, "Medical record updated successfully", nil)
}

func (h *MedicalRecordHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ResponseJSON(c, fiber.StatusBadRequest, "Record ID is required", nil)
	}

	if err := h.recordService.Delete(c.Context(), id); err != nil {
		log.Printf("Error deleting medical record: %v", err)
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to delete medical record", nil)
	}

	return utils.ResponseJSON(c, fiber.StatusOK, "Medical record deleted successfully", nil)
}
