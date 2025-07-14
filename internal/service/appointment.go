package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentService struct {
	appRepo     *repository.AppointmentRepository
	patientRepo *repository.PatientRepository
	recordRepo  *repository.MedicalRecordRepository
	activityService *ActivityService
	mongoClient *mongo.Client
}

func NewAppointmentService(
	repo *repository.AppointmentRepository,
	patientRepo *repository.PatientRepository,
	recordRepo *repository.MedicalRecordRepository,
	activityService *ActivityService,
	mongoClient *mongo.Client,
) *AppointmentService {
	return &AppointmentService{
		appRepo:     repo,
		patientRepo: patientRepo,
		recordRepo:  recordRepo,
		activityService: activityService,
		mongoClient: mongoClient,
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

func (s *AppointmentService) Create(ctx context.Context, appointment *domain.AppointmentEntity) (string, error) {
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

		// Validate required fields
		if appointment.PatientID == primitive.NilObjectID {
			return errors.New("patient ID is required")
		}

		if appointment.DoctorID == primitive.NilObjectID {
			return errors.New("doctor ID is required")
		}

		if appointment.Type == "" {
			return errors.New("appointment type is required")
		}

		// Validate appointment type
		if !appointment.Type.IsValid() {
			return fmt.Errorf("invalid appointment type: %s", appointment.Type)
		}

		if appointment.DateTime.IsZero() {
			return errors.New("appointment date and time is required")
		}

		if appointment.Duration <= 0 {
			return errors.New("appointment duration must be greater than 0")
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
			return fmt.Errorf("invalid appointment status: %s", appointment.Status)
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

		// Create the appointment
		id, err := s.appRepo.Create(sessionContext, appointment)
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}
		newAppointmentID = id.Hex()

		// Log activity
		err = s.activityService.CreateActivity(sessionContext, domain.ActivityTypeAppointment, "New Appointment Scheduled", fmt.Sprintf("Appointment for patient %s with doctor %s on %s has been scheduled.", appointment.PatientID.Hex(), appointment.DoctorID.Hex(), appointment.DateTime.Format(time.RFC3339)))
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

func (s *AppointmentService) Update(ctx context.Context, id string, appointment *domain.AppointmentEntity) error {
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

		// Validate status transition
		if appointment.Status != "" && !appointment.Status.IsValid() {
			return errors.New("invalid appointment status")
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
				sessionContext,
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

		if err := s.appRepo.Update(sessionContext, appointmentID, appointment); err != nil {
			return fmt.Errorf("failed to update appointment: %w", err)
		}

		err = s.activityService.CreateActivity(sessionContext, domain.ActivityTypeAppointment, "Appointment Updated", fmt.Sprintf("Appointment %s has been updated. New status: %s.", id, appointment.Status))
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

	err = s.activityService.CreateActivity(ctx, domain.ActivityTypeAppointment, "Appointment Cancelled", fmt.Sprintf("Appointment %s has been cancelled.", id))
	if err != nil {
		// Log the error but don't block the appointment cancellation
		fmt.Printf("Warning: failed to log activity for appointment cancellation: %v\n", err)
	}

	return nil
}
