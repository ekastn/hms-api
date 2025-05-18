package service

import (
	"context"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
)

type DashboardService struct {
	patientRepo *repository.PatientRepository
	docRepo     *repository.DoctorRepository
	apptRepo    *repository.AppointmentRepository
	recordRepo  *repository.MedicalRecordRepository
}

func NewDashboardService(
	patientRepo *repository.PatientRepository,
	docRepo *repository.DoctorRepository,
	apptRepo *repository.AppointmentRepository,
	recordRepo *repository.MedicalRecordRepository,
) *DashboardService {
	return &DashboardService{
		patientRepo: patientRepo,
		docRepo:     docRepo,
		apptRepo:    apptRepo,
		recordRepo:  recordRepo,
	}
}

func (s *DashboardService) GetDashboardData(ctx context.Context) (*domain.DashboardResponse, error) {
	// Get counts in parallel
	patientsCh := make(chan int64)
	doctorsCh := make(chan int64)
	apptsCh := make(chan int64)
	recordsCh := make(chan int64)
	errCh := make(chan error, 4)

	// Get patients count
	go func() {
		count, err := s.patientRepo.Count(ctx)
		if err != nil {
			errCh <- err
			return
		}
		patientsCh <- count
	}()

	// Get doctors count
	go func() {
		count, err := s.docRepo.Count(ctx)
		if err != nil {
			errCh <- err
			return
		}
		doctorsCh <- count
	}()

	// Get appointments count
	go func() {
		count, err := s.apptRepo.GetAppointmentsCount(ctx)
		if err != nil {
			errCh <- err
			return
		}
		apptsCh <- count
	}()

	// Get medical records count
	go func() {
		count, err := s.recordRepo.Count(ctx)
		if err != nil {
			errCh <- err
			return
		}
		recordsCh <- count
	}()

	// Wait for all counts
	var patientsCount, doctorsCount, apptsCount, recordsCount int64
	for i := 0; i < 4; i++ {
		select {
		case count := <-patientsCh:
			patientsCount = count
		case count := <-doctorsCh:
			doctorsCount = count
		case count := <-apptsCh:
			apptsCount = count
		case count := <-recordsCh:
			recordsCount = count
		case err := <-errCh:
			return nil, err
		}
	}

	// Get upcoming appointments
	upcomingAppts, err := s.apptRepo.GetUpcomingAppointments(ctx, 5) // Get next 5 upcoming appointments
	if err != nil {
		return nil, err
	}

	// Convert []*UpcomingAppointment to []UpcomingAppointment
	var upcomingAppointments []domain.UpcomingAppointment
	for _, appt := range upcomingAppts {
		if appt != nil {
			upcomingAppointments = append(upcomingAppointments, *appt)
		}
	}

	// Get recent activities (last 10 activities)
	recentActivities, err := s.getRecentActivities(ctx, 10)
	if err != nil {
		return nil, err
	}

	return &domain.DashboardResponse{
		Stats: domain.DashboardStats{
			PatientsCount:       patientsCount,
			DoctorsCount:        doctorsCount,
			AppointmentsCount:   apptsCount,
			MedicalRecordsCount: recordsCount,
		},
		RecentActivities:     recentActivities,
		UpcomingAppointments: upcomingAppointments,
	}, nil
}

func (s *DashboardService) getRecentActivities(ctx context.Context, limit int) ([]domain.Activity, error) {
	// TODO: This is a simplified implementation
	// 1. Create a dedicated activities collection
	// 2. Log activities when important events happen (appointment created, patient registered, etc.)
	// 3. Query the activities collection with proper sorting and limiting

	return []domain.Activity{
		{
			ID:          "1",
			Type:        "APPOINTMENT",
			Title:       "New appointment scheduled",
			Description: "Dr. Smith has a new appointment with John Doe",
			Timestamp:   time.Now(),
		},
		{
			ID:          "2",
			Type:        "MEDICAL_RECORD",
			Title:       "Medical record updated",
			Description: "Updated medical record for patient Jane Smith",
			Timestamp:   time.Now().Add(-1 * time.Hour),
		},
	}, nil
}
