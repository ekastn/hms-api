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
	doctorRepo      *repository.DoctorRepository
	appointmentRepo *repository.AppointmentRepository
	patientRepo     *repository.PatientRepository
	activityService *ActivityService
}

func NewDoctorService(
	repo *repository.DoctorRepository,
	appointmentRepo *repository.AppointmentRepository,
	patientRepo *repository.PatientRepository,
	activityService *ActivityService,
) *DoctorService {
	return &DoctorService{
		doctorRepo:      repo,
		appointmentRepo: appointmentRepo,
		patientRepo:     patientRepo,
		activityService: activityService,
	}
}

func (s *DoctorService) GetAll(ctx context.Context) ([]*domain.DoctorEntity, error) {
	return s.doctorRepo.GetAll(ctx)
}

func (s *DoctorService) GetByID(ctx context.Context, id string) (*domain.DoctorEntity, error) {
	doctorID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	return s.doctorRepo.GetByID(ctx, doctorID)
}

func (s *DoctorService) GetDoctorDetail(ctx context.Context, id string) (*domain.DoctorDetailResponse, error) {
	// Get the doctor
	doctorID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	doctor, err := s.doctorRepo.GetByID(ctx, doctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}
	if doctor == nil {
		return nil, nil
	}

	// Get recent patient IDs
	patientIDs, err := s.appointmentRepo.GetRecentPatientsByDoctorID(ctx, doctorID, 10) // Last 10 patients
	if err != nil {
		return nil, fmt.Errorf("failed to get recent patients: %w", err)
	}

	// Get patient details
	recentPatients := make([]*domain.PatientEntity, 0, len(patientIDs))
	for _, pid := range patientIDs {
		patient, err := s.patientRepo.GetByID(ctx, pid)
		if err != nil {
			// Log the error but continue with other patients
			continue
		}
		recentPatients = append(recentPatients, patient)
	}

	// Convert to detail DTO
	return doctor.ToDetailDTO(recentPatients), nil
}

func (s *DoctorService) Create(ctx context.Context, doctor *domain.DoctorEntity, creatorID primitive.ObjectID) (string, error) {
	// Validate required fields
	if doctor.Name == "" || doctor.Specialty == "" || doctor.Phone == "" || doctor.Email == "" {
		return "", fmt.Errorf("all fields are required")
	}

	// Check for duplicates
	if err := s.checkDuplicateDoctor(ctx, doctor, primitive.NilObjectID); err != nil {
		return "", err
	}

	// Set timestamps
	now := time.Now()
	doctor.CreatedAt = now
	doctor.UpdatedAt = now
	doctor.CreatedBy = creatorID
	doctor.UpdatedBy = creatorID

	// Create doctor
	id, err := s.doctorRepo.Create(ctx, doctor)
	if err != nil {
		return "", fmt.Errorf("failed to create doctor: %w", err)
	}

	err = s.activityService.CreateActivity(ctx, domain.ActivityTypeDoctor, "New Doctor Added", fmt.Sprintf("Dr. %s (%s) has been added as a %s specialist.", doctor.Name, doctor.Email, doctor.Specialty))
	if err != nil {
		// Log the error but don't block the doctor creation
		fmt.Printf("Warning: failed to log activity for new doctor: %v\n", err)
	}

	return id.Hex(), nil
}

func (s *DoctorService) Update(ctx context.Context, id string, doctor *domain.DoctorEntity, updaterID primitive.ObjectID) error {
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

	// Get existing doctor to preserve created_at
	existing, err := s.doctorRepo.GetByID(ctx, docID)
	if err != nil {
		return fmt.Errorf("failed to get doctor: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("doctor not found")
	}

	// Check for duplicates, excluding the current doctor
	if err := s.checkDuplicateDoctor(ctx, doctor, docID); err != nil {
		return err
	}

	// Preserve created_at and update updated_at
	doctor.CreatedAt = existing.CreatedAt
	doctor.UpdatedAt = time.Now()
	doctor.CreatedBy = existing.CreatedBy
	doctor.UpdatedBy = updaterID

	if err := s.doctorRepo.Update(ctx, docID, doctor); err != nil {
		return fmt.Errorf("failed to update doctor: %w", err)
	}

	err = s.activityService.CreateActivity(ctx, domain.ActivityTypeDoctor, "Doctor Information Updated", fmt.Sprintf("Dr. %s (%s) information has been updated.", doctor.Name, doctor.Email))
	if err != nil {
		// Log the error but don't block the doctor update
		fmt.Printf("Warning: failed to log activity for doctor update: %v\n", err)
	}

	return nil
}

// checkDuplicateDoctor is a helper function to check for existing doctors by name, email, or phone.
// It excludes the doctor with currentDoctorID from the check (useful for update operations).
func (s *DoctorService) checkDuplicateDoctor(ctx context.Context, doctor *domain.DoctorEntity, currentDoctorID primitive.ObjectID) error {
	// Check if name already exists
	existingByName, err := s.doctorRepo.GetByName(ctx, doctor.Name)
	if err != nil {
		return fmt.Errorf("error checking name: %w", err)
	}
	if existingByName != nil && existingByName.ID != currentDoctorID {
		return fmt.Errorf("doctor with this name already exists")
	}

	// Check if email already exists
	existingByEmail, err := s.doctorRepo.GetByEmail(ctx, doctor.Email)
	if err != nil {
		return fmt.Errorf("error checking email: %w", err)
	}
	if existingByEmail != nil && existingByEmail.ID != currentDoctorID {
		return fmt.Errorf("email already exists")
	}

	// Check if phone already exists
	existingByPhone, err := s.doctorRepo.GetByPhone(ctx, doctor.Phone)
	if err != nil {
		return fmt.Errorf("error checking phone: %w", err)
	}
	if existingByPhone != nil && existingByPhone.ID != currentDoctorID {
		return fmt.Errorf("phone number already exists")
	}
	return nil
}

func (s *DoctorService) Delete(ctx context.Context, id string) error {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	existingDoc, err := s.doctorRepo.GetByID(ctx, docID)
	if err != nil {
		return fmt.Errorf("failed to get doctor: %w", err)
	}

	if existingDoc == nil {
		return fmt.Errorf("doctor not found")
	}

	if err := s.doctorRepo.Delete(ctx, docID); err != nil {
		return fmt.Errorf("failed to delete doctor: %w", err)
	}

	err = s.activityService.CreateActivity(ctx, domain.ActivityTypeDoctor, "Doctor Deleted", fmt.Sprintf("Doctor with ID %s has been deleted.", id))
	if err != nil {
		// Log the error but don't block the doctor deletion
		fmt.Printf("Warning: failed to log activity for doctor deletion: %v\n", err)
	}

	return nil
}
