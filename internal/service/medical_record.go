package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MedicalRecordService struct {
	recordRepo *repository.MedicalRecordRepository
	activityService *ActivityService
}

func NewMedicalRecordService(recordRepo *repository.MedicalRecordRepository, activityService *ActivityService) *MedicalRecordService {
	return &MedicalRecordService{
		recordRepo: recordRepo,
		activityService: activityService,
	}
}

// GetAll retrieves all medical records
func (s *MedicalRecordService) GetAll(ctx context.Context) ([]*domain.MedicalRecordEntity, error) {
	records, err := s.recordRepo.FindAll(ctx)
	if err != nil {
		log.Printf("Error getting all medical records: %v", err)
		return nil, fmt.Errorf("failed to get medical records")
	}
	return records, nil
}

func (s *MedicalRecordService) Create(ctx context.Context, record *domain.MedicalRecordEntity) (string, error) {
	if err := validateMedicalRecord(record); err != nil {
		return "", err
	}

	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	id, err := s.recordRepo.Create(ctx, record)
	if err != nil {
		return "", fmt.Errorf("failed to create medical record: %w", err)
	}

	err = s.activityService.CreateActivity(ctx, domain.ActivityTypeMedicalRecord, "New Medical Record Created", fmt.Sprintf("Medical record for patient %s (diagnosis: %s) has been created.", record.PatientID.Hex(), record.Diagnosis))
	if err != nil {
		// Log the error but don't block the medical record creation
		fmt.Printf("Warning: failed to log activity for new medical record: %v\n", err)
	}

	return id.Hex(), nil
}

func (s *MedicalRecordService) GetByID(ctx context.Context, id string) (*domain.MedicalRecordEntity, error) {
	recordID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	record, err := s.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical record: %w", err)
	}

	return record, nil
}

func (s *MedicalRecordService) GetByPatientID(ctx context.Context, patientID string) ([]*domain.MedicalRecordEntity, error) {
	patientObjID, err := primitive.ObjectIDFromHex(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient ID format: %w", err)
	}

	records, err := s.recordRepo.GetByPatientID(ctx, patientObjID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient's medical records: %w", err)
	}

	return records, nil
}

func (s *MedicalRecordService) GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.MedicalRecordEntity, error) {
	records, err := s.recordRepo.GetByDateRange(ctx, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical records: %w", err)
	}

	return records, nil
}

func (s *MedicalRecordService) Update(ctx context.Context, id string, record *domain.MedicalRecordEntity) error {
	recordID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	

	if err := validateMedicalRecord(record); err != nil {
		return err
	}

	record.UpdatedAt = time.Now()

	if err := s.recordRepo.Update(ctx, recordID, record); err != nil {
		return fmt.Errorf("failed to update medical record: %w", err)
	}

	err = s.activityService.CreateActivity(ctx, domain.ActivityTypeMedicalRecord, "Medical Record Updated", fmt.Sprintf("Medical record %s for patient %s has been updated.", id, record.PatientID.Hex()))
	if err != nil {
		// Log the error but don't block the medical record update
		fmt.Printf("Warning: failed to log activity for medical record update: %v\n", err)
	}

	return nil
}

func (s *MedicalRecordService) Delete(ctx context.Context, id string) error {
	recordID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Check if record exists
	existingRecord, err := s.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return fmt.Errorf("failed to get medical record: %w", err)
	}

	if existingRecord == nil {
		return errors.New("medical record not found")
	}

	if err := s.recordRepo.Delete(ctx, recordID); err != nil {
		return fmt.Errorf("failed to delete medical record: %w", err)
	}

	err = s.activityService.CreateActivity(ctx, domain.ActivityTypeMedicalRecord, "Medical Record Deleted", fmt.Sprintf("Medical record %s has been deleted.", id))
	if err != nil {
		// Log the error but don't block the medical record deletion
		fmt.Printf("Warning: failed to log activity for medical record deletion: %v\n", err)
	}

	return nil
}

func validateMedicalRecord(record *domain.MedicalRecordEntity) error {
	if record.PatientID.IsZero() {
		return errors.New("patient ID is required")
	}

	if record.DoctorID.IsZero() {
		return errors.New("doctor ID is required")
	}

	if record.RecordType == "" {
		return errors.New("record type is required")
	}

	if !record.RecordType.IsValid() {
		return fmt.Errorf("invalid record type: %s", record.RecordType)
	}

	if record.Description == "" {
		return errors.New("description is required")
	}

	if record.Diagnosis == "" {
		return errors.New("diagnosis is required")
	}

	if record.Treatment == "" {
		return errors.New("treatment is required")
	}

	return nil
}
