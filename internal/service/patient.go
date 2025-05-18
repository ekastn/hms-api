package service

import (
	"context"
	"fmt"

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
	id, err := s.patientRepo.Create(ctx, patient)
	if err != nil {
		return "", fmt.Errorf("failed to create patient: %w", err)
	}
	return id.Hex(), nil
}

func (s *PatientService) Update(ctx context.Context, id string, patient *domain.PatientEntity) error {
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
