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
	patientRepo *repository.PatientRepository
}

func NewPatientService(patientRepo *repository.PatientRepository) *PatientService {
	return &PatientService{
		patientRepo: patientRepo,
	}
}

func (s *PatientService) GetAll(ctx context.Context) ([]*domain.PatientEntity, error) {
	return s.patientRepo.GetAll(ctx)
}

func (s *PatientService) GetByID(ctx context.Context, id string) (*domain.PatientEntity, error) {
	patientID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	patient, err := s.patientRepo.GetByID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	return patient, nil
}

func (s *PatientService) Create(ctx context.Context, patient *domain.PatientEntity) (string, error) {
	// Validate required fields
	if patient.Name == "" || patient.Phone == "" || patient.Email == "" {
		return "", errors.New("name, phone, and email are required")
	}

	// Check if name already exists
	existingByName, err := s.patientRepo.GetByName(ctx, patient.Name)
	if err != nil {
		return "", fmt.Errorf("error checking name: %w", err)
	}
	if existingByName != nil {
		return "", errors.New("patient with this name already exists")
	}

	// Check if email already exists
	existingByEmail, err := s.patientRepo.GetByEmail(ctx, patient.Email)
	if err != nil {
		return "", fmt.Errorf("error checking email: %w", err)
	}
	if existingByEmail != nil {
		return "", errors.New("email already exists")
	}

	// Check if phone already exists
	existingByPhone, err := s.patientRepo.GetByPhone(ctx, patient.Phone)
	if err != nil {
		return "", fmt.Errorf("error checking phone: %w", err)
	}
	if existingByPhone != nil {
		return "", errors.New("phone number already exists")
	}

	// Set last visit timestamp
	patient.LastVisit = time.Now()

	// Insert patient into database
	id, err := s.patientRepo.Create(ctx, patient)
	if err != nil {
		return "", fmt.Errorf("failed to create patient: %w", err)
	}

	return id.Hex(), nil
}

func (s *PatientService) Update(ctx context.Context, id string, patient *domain.PatientEntity) error {
	if id == "" {
		return errors.New("patient ID is required")
	}

	if patient.Name == "" || patient.Phone == "" || patient.Email == "" {
		return errors.New("name, phone, and email are required")
	}

	// Convert ID to ObjectID
	patientID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Get existing patient
	existing, err := s.patientRepo.GetByID(ctx, patientID)
	if err != nil {
		return fmt.Errorf("error getting patient: %w", err)
	}
	if existing == nil {
		return errors.New("patient not found")
	}

	// Check if name is being changed and already exists
	if existing.Name != patient.Name {
		existingByName, err := s.patientRepo.GetByName(ctx, patient.Name)
		if err != nil {
			return fmt.Errorf("error checking name: %w", err)
		}
		if existingByName != nil {
			return errors.New("patient with this name already exists")
		}
	}

	// Check if email is being changed and already exists
	if existing.Email != patient.Email {
		existingByEmail, err := s.patientRepo.GetByEmail(ctx, patient.Email)
		if err != nil {
			return fmt.Errorf("error checking email: %w", err)
		}
		if existingByEmail != nil {
			return errors.New("email already exists")
		}
	}

	// Check if phone is being changed and already exists
	if existing.Phone != patient.Phone {
		existingByPhone, err := s.patientRepo.GetByPhone(ctx, patient.Phone)
		if err != nil {
			return fmt.Errorf("error checking phone: %w", err)
		}
		if existingByPhone != nil {
			return errors.New("phone number already exists")
		}
	}

	// Update last visit timestamp
	patient.LastVisit = time.Now()

	// Update patient in database
	if err := s.patientRepo.Update(ctx, patientID, patient); err != nil {
		return fmt.Errorf("failed to update patient: %w", err)
	}

	return nil
}

func (s *PatientService) Delete(ctx context.Context, id string) error {
	patientID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	patientFound, err := s.patientRepo.GetByID(ctx, patientID)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}

	if patientFound == nil {
		return fmt.Errorf("patient not found")
	}

	if err := s.patientRepo.Delete(ctx, patientID); err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	return nil
}
