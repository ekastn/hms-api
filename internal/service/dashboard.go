package service

import (
	"context"
	

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
)

type DashboardService struct {
	patientRepo *repository.PatientRepository
	docRepo     *repository.DoctorRepository
	apptRepo    *repository.AppointmentRepository
	recordRepo  *repository.MedicalRecordRepository
	activityRepo *repository.ActivityRepository
}

func NewDashboardService(
	patientRepo *repository.PatientRepository,
	docRepo *repository.DoctorRepository,
	apptRepo *repository.AppointmentRepository,
	recordRepo *repository.MedicalRecordRepository,
	activityRepo *repository.ActivityRepository,
) *DashboardService {
	return &DashboardService{
		patientRepo: patientRepo,
		docRepo:     docRepo,
		apptRepo:    apptRepo,
		recordRepo:  recordRepo,
		activityRepo: activityRepo,
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
	recentActivities, err := s.activityRepo.GetRecent(ctx, 10)
	if err != nil {
		return nil, err
	}

	var recentActivitiesDTOs []domain.Activity
	for _, activity := range recentActivities {
		if activity != nil {
			recentActivitiesDTOs = append(recentActivitiesDTOs, domain.Activity{
				ID:          activity.ID.Hex(),
				Type:        string(activity.Type),
				Title:       activity.Title,
				Description: activity.Description,
				Timestamp:   activity.Timestamp,
			})
		}
	}

	return &domain.DashboardResponse{
		Stats: domain.DashboardStats{
			PatientsCount:       patientsCount,
			DoctorsCount:        doctorsCount,
			AppointmentsCount:   apptsCount,
			MedicalRecordsCount: recordsCount,
		},
		RecentActivities:     recentActivitiesDTOs,
		UpcomingAppointments: upcomingAppointments,
	}, nil
}
