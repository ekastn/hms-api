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

type AppointmentService struct {
	appRepo *repository.AppointmentRepository
}

func NewAppointmentService(repo *repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{
		appRepo: repo,
	}
}

func (s *AppointmentService) GetAll(ctx context.Context) ([]*domain.AppointmentEntity, error) {
	return s.appRepo.GetAll(ctx)
}

func (s *AppointmentService) GetByID(ctx context.Context, id string) (*domain.AppointmentEntity, error) {
	appointmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	appointment, err := s.appRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return appointment, nil
}

func (s *AppointmentService) Create(ctx context.Context, appointment *domain.AppointmentEntity) (string, error) {
	// Validate required fields
	if appointment.PatientID == primitive.NilObjectID {
		return "", errors.New("patient ID is required")
	}

	if appointment.DoctorID == primitive.NilObjectID {
		return "", errors.New("doctor ID is required")
	}

	if appointment.Type == "" {
		return "", errors.New("appointment type is required")
	}

	// Validate appointment type
	switch appointment.Type {
	case domain.AppointmentTypeCheckUp, domain.AppointmentTypeFollowUp,
		domain.AppointmentTypeConsultation, domain.AppointmentTypeProcedure,
		domain.AppointmentTypeEmergency:
		// valid type
	default:
		return "", fmt.Errorf("invalid appointment type: %s", appointment.Type)
	}

	if appointment.DateTime.IsZero() {
		return "", errors.New("appointment date and time is required")
	}

	if appointment.Duration <= 0 {
		return "", errors.New("appointment duration must be greater than 0")
	}

	// Set default status if not provided
	if appointment.Status == "" {
		appointment.Status = domain.AppointmentStatusScheduled
	}

	// Validate status
	switch appointment.Status {
	case domain.AppointmentStatusScheduled, domain.AppointmentStatusConfirmed,
		domain.AppointmentStatusCompleted, domain.AppointmentStatusCancelled:
		// valid status
	default:
		return "", fmt.Errorf("invalid appointment status: %s", appointment.Status)
	}

	// Check for existing appointment at the same time
	existing, err := s.appRepo.GetByDoctorAndDateRange(
		ctx,
		appointment.DoctorID,
		appointment.DateTime,
		appointment.DateTime.Add(30*time.Minute),
	)
	if err != nil {
		return "", fmt.Errorf("failed to check for existing appointments: %w", err)
	}

	// Filter out cancelled appointments from the conflict check
	var activeAppointments []*domain.AppointmentEntity
	for _, a := range existing {
		if a.Status != domain.AppointmentStatusCancelled {
			activeAppointments = append(activeAppointments, a)
		}
	}

	if len(activeAppointments) > 0 {
		return "", errors.New("doctor is not available at the requested time")
	}

	// Create the appointment
	id, err := s.appRepo.Create(ctx, appointment)
	if err != nil {
		return "", fmt.Errorf("failed to create appointment: %w", err)
	}

	return id.Hex(), nil
}

func (s *AppointmentService) Update(ctx context.Context, id string, appointment *domain.AppointmentEntity) error {
	appointmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Get existing appointment to preserve created_at and other fields
	existingAppointment, err := s.appRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return fmt.Errorf("failed to get appointment: %w", err)
	}

	if existingAppointment == nil {
		return errors.New("appointment not found")
	}

	// Validate status transition
	if appointment.Status != "" && existingAppointment.Status != appointment.Status {
		switch existingAppointment.Status {
		case domain.AppointmentStatusCompleted:
			return errors.New("cannot modify a completed appointment")
		case domain.AppointmentStatusCancelled:
			if appointment.Status != domain.AppointmentStatusScheduled {
				return errors.New("can only reschedule a cancelled appointment")
			}
		}
	}

	// Preserve immutable fields
	appointment.ID = appointmentID
	appointment.CreatedAt = existingAppointment.CreatedAt
	appointment.UpdatedAt = time.Now()

	// If date/time or doctor is being changed, check for conflicts
	if !appointment.DateTime.Equal(existingAppointment.DateTime) ||
		appointment.DoctorID != existingAppointment.DoctorID ||
		appointment.Duration != existingAppointment.Duration {

		endTime := appointment.DateTime.Add(time.Duration(appointment.Duration) * time.Minute)
		existingAppointments, err := s.appRepo.GetByDoctorAndDateRange(
			ctx,
			appointment.DoctorID,
			appointment.DateTime,
			endTime,
		)
		if err != nil {
			return fmt.Errorf("failed to check for existing appointments: %w", err)
		}

		// Filter out the current appointment and cancelled appointments
		var conflictingAppointments []*domain.AppointmentEntity
		for _, a := range existingAppointments {
			if a.ID != appointmentID && a.Status != domain.AppointmentStatusCancelled {
				conflictingAppointments = append(conflictingAppointments, a)
			}
		}

		if len(conflictingAppointments) > 0 {
			return errors.New("doctor is not available at the requested time")
		}
	}

	if err := s.appRepo.Update(ctx, appointmentID, appointment); err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	return nil
}

func (s *AppointmentService) Delete(ctx context.Context, id string) error {
	appointmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Check if appointment exists
	existingAppointment, err := s.appRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return fmt.Errorf("failed to get appointment: %w", err)
	}

	if existingAppointment == nil {
		return errors.New("appointment not found")
	}

	// Prevent deletion of completed appointments
	if existingAppointment.Status == domain.AppointmentStatusCompleted {
		return errors.New("cannot delete a completed appointment")
	}

	// Instead of deleting, we'll mark it as cancelled
	existingAppointment.Status = domain.AppointmentStatusCancelled
	err = s.appRepo.Update(ctx, appointmentID, existingAppointment)
	if err != nil {
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}

	return nil
}
