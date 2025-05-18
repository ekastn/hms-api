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
	// Validate required fields
	if doctor.Name == "" || doctor.Specialty == "" || doctor.Phone == "" || doctor.Email == "" {
		return "", fmt.Errorf("all fields are required")
	}

	// Check if name already exists
	existingByName, err := s.docRepo.GetByName(ctx, doctor.Name)
	if err != nil {
		return "", fmt.Errorf("error checking name: %w", err)
	}
	if existingByName != nil {
		return "", fmt.Errorf("doctor with this name already exists")
	}

	// Check if email already exists
	existingByEmail, err := s.docRepo.GetByEmail(ctx, doctor.Email)
	if err != nil {
		return "", fmt.Errorf("error checking email: %w", err)
	}
	if existingByEmail != nil {
		return "", fmt.Errorf("email already exists")
	}

	// Check if phone already exists
	existingByPhone, err := s.docRepo.GetByPhone(ctx, doctor.Phone)
	if err != nil {
		return "", fmt.Errorf("error checking phone: %w", err)
	}
	if existingByPhone != nil {
		return "", fmt.Errorf("phone number already exists")
	}

	// Set timestamps
	now := time.Now()
	doctor.CreatedAt = now
	doctor.UpdatedAt = now

	// Create doctor
	id, err := s.docRepo.Create(ctx, doctor)
	if err != nil {
		return "", fmt.Errorf("failed to create doctor: %w", err)
	}

	return id.Hex(), nil
}

func (s *DoctorService) Update(ctx context.Context, id string, doctor *domain.DoctorEntity) error {
	if id == "" {
		return fmt.Errorf("doctor ID is required")
	}

	if doctor.Name == "" || doctor.Specialty == "" || doctor.Phone == "" || doctor.Email == "" {
		return fmt.Errorf("all fields are required")
	}

	// Convert ID to ObjectID
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Get existing doctor
	existing, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return fmt.Errorf("failed to get doctor: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("doctor not found")
	}

	// Check if name is being changed and already exists
	if existing.Name != doctor.Name {
		existingByName, err := s.docRepo.GetByName(ctx, doctor.Name)
		if err != nil {
			return fmt.Errorf("error checking name: %w", err)
		}
		if existingByName != nil {
			return fmt.Errorf("doctor with this name already exists")
		}
	}

	// Check if email is being changed and already exists
	if existing.Email != doctor.Email {
		existingByEmail, err := s.docRepo.GetByEmail(ctx, doctor.Email)
		if err != nil {
			return fmt.Errorf("error checking email: %w", err)
		}
		if existingByEmail != nil {
			return fmt.Errorf("email already exists")
		}
	}

	// Check if phone is being changed and already exists
	if existing.Phone != doctor.Phone {
		existingByPhone, err := s.docRepo.GetByPhone(ctx, doctor.Phone)
		if err != nil {
			return fmt.Errorf("error checking phone: %w", err)
		}
		if existingByPhone != nil {
			return fmt.Errorf("phone number already exists")
		}
	}

	// Preserve created_at and update updated_at
	doctor.CreatedAt = existing.CreatedAt
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
