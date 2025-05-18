package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorService struct {
	docRepo *repository.DoctorRepository
}

func NewDoctorService(docRepo *repository.DoctorRepository) *DoctorService {
	return &DoctorService{
		docRepo: docRepo,
	}
}

func (s *DoctorService) GetAll(ctx context.Context) ([]*domain.DoctorEntity, error) {
	return s.docRepo.GetAll(ctx)
}

func (s *DoctorService) GetByID(ctx context.Context, id string) (*domain.DoctorEntity, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	doctor, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}

	return doctor, nil
}

func (s *DoctorService) Create(ctx context.Context, doctor *domain.DoctorEntity) (string, error) {
	now := time.Now()
	doctor.CreatedAt = now
	doctor.UpdatedAt = now

	id, err := s.docRepo.Create(ctx, doctor)
	if err != nil {
		return "", fmt.Errorf("failed to create doctor: %w", err)
	}
	return id.Hex(), nil
}

func (s *DoctorService) Update(ctx context.Context, id string, doctor *domain.DoctorEntity) error {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	existingDoc, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return fmt.Errorf("failed to get doctor: %w", err)
	}

	if existingDoc == nil {
		return fmt.Errorf("doctor not found")
	}

	// Preserve created_at and update updated_at
	doctor.CreatedAt = existingDoc.CreatedAt
	doctor.UpdatedAt = time.Now()

	if err := s.docRepo.Update(ctx, docID, doctor); err != nil {
		return fmt.Errorf("failed to update doctor: %w", err)
	}

	return nil
}

func (s *DoctorService) Delete(ctx context.Context, id string) error {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	existingDoc, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return fmt.Errorf("failed to get doctor: %w", err)
	}

	if existingDoc == nil {
		return fmt.Errorf("doctor not found")
	}

	if err := s.docRepo.Delete(ctx, docID); err != nil {
		return fmt.Errorf("failed to delete doctor: %w", err)
	}

	return nil
}
