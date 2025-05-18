package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PatientService struct {
	docRepo    *repository.PatientRepository
	apptRepo   *repository.AppointmentRepository
	recordRepo *repository.MedicalRecordRepository
}

func NewPatientService(
	repo *repository.PatientRepository,
	apptRepo *repository.AppointmentRepository,
	recordRepo *repository.MedicalRecordRepository,
) *PatientService {
	return &PatientService{
		docRepo:    repo,
		apptRepo:   apptRepo,
		recordRepo: recordRepo,
	}
}

func (s *PatientService) GetAll(ctx context.Context) ([]*domain.PatientEntity, error) {
	return s.docRepo.GetAll(ctx)
}

func (s *PatientService) GetByID(ctx context.Context, id string) (*domain.PatientEntity, error) {
	patientID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	patient, err := s.docRepo.GetByID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	return patient, nil
}

func (s *PatientService) GetPatientDetail(ctx context.Context, id string) (*domain.PatientDetailResponse, error) {
	// Get patient
	patient, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	// Get recent appointments (last 10)
	appointments, err := s.apptRepo.GetByPatientID(ctx, patient.ID)
	if len(appointments) > 10 {
		appointments = appointments[:10] // Limit to 10 most recent
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get patient appointments: %w", err)
	}

	// Get medical history
	medicalRecords, err := s.recordRepo.GetByPatientID(ctx, patient.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient medical records: %w", err)
	}

	// Convert to DTOs
	var appointmentDTOs []domain.AppointmentDTO
	for _, appt := range appointments {
		appointmentDTOs = append(appointmentDTOs, appt.ToDTO())
	}

	var recordDTOs []domain.MedicalRecordDTO
	for _, record := range medicalRecords {
		recordDTOs = append(recordDTOs, record.ToDTO())
	}

	return &domain.PatientDetailResponse{
		Patient:            patient.ToDTO(),
		RecentAppointments: appointmentDTOs,
		MedicalHistory:     recordDTOs,
	}, nil
}

func (s *PatientService) Create(ctx context.Context, patient *domain.PatientEntity) (string, error) {
	// Validate required fields
	if patient.Name == "" || patient.Phone == "" || patient.Email == "" {
		return "", errors.New("name, phone, and email are required")
	}

	// Set timestamps
	now := time.Now()
	patient.CreatedAt = now
	patient.UpdatedAt = now
	patient.LastVisit = now
	patient.ID = primitive.NewObjectID()

	// Create patient in repository
	id, err := s.docRepo.Create(ctx, patient)
	if err != nil {
		return "", fmt.Errorf("failed to create patient: %w", err)
	}

	return id.Hex(), nil
}

func (s *PatientService) Update(ctx context.Context, id string, patient *domain.PatientEntity) error {
	if id == "" {
		return errors.New("patient ID is required")
	}

	// Validate required fields
	if patient.Name == "" || patient.Phone == "" || patient.Email == "" {
		return errors.New("name, phone, and email are required")
	}

	patientID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Get existing patient to preserve timestamps
	existingPatient, err := s.docRepo.GetByID(ctx, patientID)
	if err != nil {
		return fmt.Errorf("failed to get existing patient: %w", err)
	}
	if existingPatient == nil {
		return errors.New("patient not found")
	}

	// Preserve created_at and set updated_at
	patient.CreatedAt = existingPatient.CreatedAt
	patient.UpdatedAt = time.Now()
	patient.ID = patientID

	// Update patient in repository
	if err := s.docRepo.Update(ctx, patientID, patient); err != nil {
		return fmt.Errorf("failed to update patient: %w", err)
	}

	return nil
}

func (s *PatientService) Delete(ctx context.Context, id string) error {
	patientID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Check if patient exists
	_, err = s.docRepo.GetByID(ctx, patientID)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}

	// Delete patient from repository
	if err := s.docRepo.Delete(ctx, patientID); err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	return nil
}
