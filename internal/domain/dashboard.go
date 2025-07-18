package domain

import "time"

// @Description	Dashboard statistics
// @swagger:model
type DashboardStats struct {
	PatientsCount       int64 `json:"patientsCount" example:100`
	DoctorsCount        int64 `json:"doctorsCount" example:20`
	AppointmentsCount   int64 `json:"appointmentsCount" example:500`
	MedicalRecordsCount int64 `json:"medicalRecordsCount" example:1200`
}

// @Description	Activity log entry
// @swagger:model
type Activity struct {
	ID          string    `json:"id" bson:"_id" example:"60d0fe4f53115a001f000001"`
	Type        string    `json:"type" example:"APPOINTMENT"` // e.g., "APPOINTMENT", "MEDICAL_RECORD"
	Title       string    `json:"title" example:"New Appointment Scheduled"`
	Description string    `json:"description" example:"Appointment for patient John Doe with doctor Jane Smith on 2025-07-17." `
	Timestamp   time.Time `json:"timestamp" example:"2025-07-17T10:30:00Z"`
}

// @Description	Upcoming appointment details
// @swagger:model
type UpcomingAppointment struct {
	ID          string    `json:"id" bson:"_id" example:"60d0fe4f53115a001f000002"`
	PatientName string    `json:"patientName" example:"John Doe"`
	DoctorName  string    `json:"doctorName" example:"Jane Smith"`
	Date        time.Time `json:"date" example:"2025-07-18T14:00:00Z"`
	Status      string    `json:"status" example:"Scheduled"`
}

// @Description	Dashboard response containing statistics, recent activities, and upcoming appointments
// @swagger:model
type DashboardResponse struct {
	Stats                DashboardStats        `json:"stats"`
	RecentActivities     []Activity            `json:"recentActivities"`
	UpcomingAppointments []UpcomingAppointment `json:"upcomingAppointments"`
}
