package domain

import "time"

type DashboardStats struct {
	PatientsCount       int64 `json:"patientsCount"`
	DoctorsCount        int64 `json:"doctorsCount"`
	AppointmentsCount   int64 `json:"appointmentsCount"`
	MedicalRecordsCount int64 `json:"medicalRecordsCount"`
}

type Activity struct {
	ID          string    `json:"id" bson:"_id"`
	Type        string    `json:"type"` // e.g., "APPOINTMENT", "MEDICAL_RECORD"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type UpcomingAppointment struct {
	ID          string    `json:"id" bson:"_id"`
	PatientName string    `json:"patientName"`
	DoctorName  string    `json:"doctorName"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
}

type DashboardResponse struct {
	Stats                DashboardStats        `json:"stats"`
	RecentActivities     []Activity            `json:"recentActivities"`
	UpcomingAppointments []UpcomingAppointment `json:"upcomingAppointments"`
}
