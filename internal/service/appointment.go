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
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentService struct {
	appRepo         *repository.AppointmentRepository
	patientRepo     *repository.PatientRepository
	recordRepo      *repository.MedicalRecordRepository
	activityService *ActivityService
	mongoClient     *mongo.Client
}

func NewAppointmentService(
	repo *repository.AppointmentRepository,
	patientRepo *repository.PatientRepository,
	recordRepo *repository.MedicalRecordRepository,
	activityService *ActivityService,
	mongoClient *mongo.Client,
) *AppointmentService {
	return &AppointmentService{
		appRepo:         repo,
		patientRepo:     patientRepo,
		recordRepo:      recordRepo,
		activityService: activityService,
		mongoClient:     mongoClient,
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

func (s *AppointmentService) GetAppointmentDetail(ctx context.Context, id string) (*domain.AppointmentDetailResponse, error) {
	// Get the appointment
	appointment, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}
	if appointment == nil {
		return nil, nil
	}

	// Get patient details
	var patient *domain.PatientEntity
	var lastRecord *domain.MedicalRecordEntity

	if appointment.PatientID != primitive.NilObjectID {
		// Get patient
		patient, err = s.patientRepo.GetByID(ctx, appointment.PatientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get patient: %w", err)
		}

		// Get the last medical record for the patient
		records, err := s.recordRepo.GetByPatientID(ctx, appointment.PatientID)
		if err != nil {
			// Log the error but don't fail the request
			fmt.Printf("Warning: failed to get medical records: %v\n", err)
		} else if len(records) > 0 {
			// Get the most recent record (assuming records are ordered by date descending)
			lastRecord = records[0]
		}
	}

	// Convert to detail DTO
	return appointment.ToDetailDTO(patient, lastRecord), nil
}

func (s *AppointmentService) Create(ctx context.Context, req *domain.CreateAppointmentRequest, creatorID primitive.ObjectID) (string, error) {
	// Start a session for transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return "", fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var newAppointmentID string

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}

		patientID, err := primitive.ObjectIDFromHex(req.PatientID)
		if err != nil {
			return fmt.Errorf("invalid patient ID format: %w", err)
		}
		doctorID, err := primitive.ObjectIDFromHex(req.DoctorID)
		if err != nil {
			return fmt.Errorf("invalid doctor ID format: %w", err)
		}

		appointment := domain.AppointmentEntity{
			PatientID:      patientID,
			DoctorID:       doctorID,
			Type:           req.Type,
			DateTime:       req.DateTime,
			Duration:       req.Duration,
			Location:       req.Location,
			Notes:          req.Notes,
			PatientHistory: req.PatientHistory,
			Status:         domain.AppointmentStatusScheduled, // Default status for new appointments
		}

		// Check for existing appointment at the same time
		existing, err := s.appRepo.GetByDoctorAndDateRange(
			sessionContext,
			appointment.DoctorID,
			appointment.DateTime,
			appointment.DateTime.Add(time.Duration(appointment.Duration)*time.Minute),
		)
		if err != nil {
			return fmt.Errorf("failed to check for existing appointments: %w", err)
		}

		// Filter out cancelled appointments from the conflict check
		var activeAppointments []*domain.AppointmentEntity
		for _, a := range existing {
			if a.Status != domain.AppointmentStatusCancelled {
				activeAppointments = append(activeAppointments, a)
			}
		}

		if len(activeAppointments) > 0 {
			return errors.New("doctor is not available at the requested time")
		}

		// Set audit fields
		appointment.CreatedBy = creatorID
		appointment.UpdatedBy = creatorID

		// Create the appointment
		id, err := s.appRepo.Create(sessionContext, &appointment)
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}
		newAppointmentID = id.Hex()

		// Log activity
		err = s.activityService.CreateActivity(sessionContext, domain.ActivityTypeAppointment, "New Appointment Scheduled", fmt.Sprintf("Appointment for patient %s with doctor %s on %s has been scheduled.", req.PatientID, req.DoctorID, req.DateTime.Format(time.RFC3339)))
		if err != nil {
			// Log the error but don't block the appointment creation
			fmt.Printf("Warning: failed to log activity for new appointment: %v\n", err)
		}

		return session.CommitTransaction(sessionContext)
	})
	if err != nil {
		session.AbortTransaction(ctx)
		return "", err
	}

	return newAppointmentID, nil
}

func (s *AppointmentService) Update(ctx context.Context, id string, req *domain.UpdateAppointmentRequest, updaterID primitive.ObjectID) error {
	appointmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Start a session for transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}

		// Get existing appointment to preserve created_at and other fields
		existingAppointment, err := s.appRepo.GetByID(sessionContext, appointmentID)
		if err != nil {
			return fmt.Errorf("failed to get appointment: %w", err)
		}

		if existingAppointment == nil {
			return errors.New("appointment not found")
		}

		// Apply updates from request to existing appointment
		updatedFields := req.ApplyUpdates(existingAppointment)

		// Update UpdatedAt and UpdatedBy if any fields were changed
		if updatedFields {
			existingAppointment.UpdatedAt = time.Now()
			existingAppointment.UpdatedBy = updaterID
		}

		// If date/time, duration, or status is being changed, re-check for conflicts
		// This logic needs to use the potentially updated values from existingAppointment
		// The condition should check if any of these fields were actually part of the request
		// or if the status changed to something that requires re-checking.
		// For simplicity, I'll keep the existing conflict check but ensure it uses the updated existingAppointment.
		// A more robust solution might involve checking if req.DateTime != nil || req.Duration != nil || req.Status != nil
		// and then performing the conflict check.
		// However, the current logic already uses existingAppointment.DateTime, existingAppointment.Duration, existingAppointment.Status
		// which will now reflect the applied updates.

		// Only check for conflicts if the status is not cancelled
		if existingAppointment.Status != domain.AppointmentStatusCancelled {
			endTime := existingAppointment.DateTime.Add(time.Duration(existingAppointment.Duration) * time.Minute)
			conflictingAppointments, err := s.appRepo.GetByDoctorAndDateRange(
				sessionContext,
				existingAppointment.DoctorID,
				existingAppointment.DateTime,
				endTime,
			)
			if err != nil {
				return fmt.Errorf("failed to check for existing appointments: %w", err)
			}

			// Filter out the current appointment and cancelled appointments
			var activeAppointments []*domain.AppointmentEntity
			for _, a := range conflictingAppointments {
				if a.ID != appointmentID && a.Status != domain.AppointmentStatusCancelled {
					activeAppointments = append(activeAppointments, a)
				}
			}

			if len(activeAppointments) > 0 {
				return errors.New("doctor is not available at the requested time")
			}
		}

		if err := s.appRepo.Update(sessionContext, appointmentID, existingAppointment); err != nil {
			return fmt.Errorf("failed to update appointment: %w", err)
		}

		err = s.activityService.CreateActivity(sessionContext, domain.ActivityTypeAppointment, "Appointment Updated", fmt.Sprintf("Appointment %s has been updated. New status: %s.", id, existingAppointment.Status))
		if err != nil {
			// Log the error but don't block the appointment update
			fmt.Printf("Warning: failed to log activity for appointment update: %v\n", err)
		}

		return session.CommitTransaction(sessionContext)
	})
	if err != nil {
		session.AbortTransaction(ctx)
		return err
	}

	return nil
}

func (s *AppointmentService) Delete(ctx context.Context, id string, updaterID primitive.ObjectID) error {
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
	existingAppointment.UpdatedBy = updaterID
	err = s.appRepo.Update(ctx, appointmentID, existingAppointment)
	if err != nil {
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}

	err = s.activityService.CreateActivity(ctx, domain.ActivityTypeAppointment, "Appointment Cancelled", fmt.Sprintf("Appointment %s has been cancelled.", id))
	if err != nil {
		// Log the error but don't block the appointment cancellation
		log.Printf("Warning: failed to log activity for appointment cancellation: %v", err)
	}

	return nil
}

// UpdateStatus updates the status of an appointment.
func (s *AppointmentService) UpdateStatus(ctx context.Context, id string, status domain.AppointmentStatus, updaterID primitive.ObjectID) error {
	appointmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Start a session for transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}

		// Get existing appointment
		existingAppointment, err := s.appRepo.GetByID(sessionContext, appointmentID)
		if err != nil {
			return fmt.Errorf("failed to get appointment: %w", err)
		}

		if existingAppointment == nil {
			return errors.New("appointment not found")
		}

		// Update only the status
		existingAppointment.Status = status
		existingAppointment.UpdatedAt = time.Now()
		existingAppointment.UpdatedBy = updaterID

		if err := s.appRepo.Update(sessionContext, appointmentID, existingAppointment); err != nil {
			return fmt.Errorf("failed to update appointment status: %w", err)
		}

		err = s.activityService.CreateActivity(sessionContext, domain.ActivityTypeAppointment, "Appointment Status Updated", fmt.Sprintf("Appointment %s status changed to %s.", id, status))
		if err != nil {
			log.Printf("Warning: failed to log activity for appointment status update: %v", err)
		}

		return session.CommitTransaction(sessionContext)
	})
	if err != nil {
		session.AbortTransaction(ctx)
		return err
	}

	return nil
}
