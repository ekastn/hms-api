package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "Scheduled"
	AppointmentStatusConfirmed AppointmentStatus = "Confirmed"
	AppointmentStatusCompleted AppointmentStatus = "Completed"
	AppointmentStatusCancelled AppointmentStatus = "Cancelled"
)

type AppointmentType string

const (
	AppointmentTypeCheckUp      AppointmentType = "check-up"
	AppointmentTypeFollowUp     AppointmentType = "follow-up"
	AppointmentTypeConsultation AppointmentType = "consultation"
	AppointmentTypeProcedure    AppointmentType = "procedure"
	AppointmentTypeEmergency    AppointmentType = "emergency"
)

type AppointmentEntity struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	PatientID      primitive.ObjectID `bson:"patientId"`
	DoctorID       primitive.ObjectID `bson:"doctorId"`
	Type           AppointmentType    `bson:"type"`
	DateTime       time.Time          `bson:"dateTime"`
	Duration       int                `bson:"duration"` // in minutes
	Status         AppointmentStatus  `bson:"status"`
	Location       string             `bson:"location"`
	Notes          string             `bson:"notes,omitempty"`
	PatientHistory string             `bson:"patientHistory,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}

type AppointmentDTO struct {
	ID             string            `json:"id"`
	PatientID      string            `json:"patientId"`
	DoctorID       string            `json:"doctorId"`
	Type           AppointmentType   `json:"type"`
	DateTime       time.Time         `json:"dateTime"`
	Duration       int               `json:"duration"` // in minutes
	Status         AppointmentStatus `json:"status"`
	Location       string            `json:"location"`
	Notes          string            `json:"notes,omitempty"`
	PatientHistory string            `json:"patientHistory,omitempty"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type AppointmentDetailResponse struct {
	Appointment AppointmentDTO    `json:"appointment"`
	Patient     *PatientDTO       `json:"patient,omitempty"`
	LastRecord  *MedicalRecordDTO `json:"lastRecord,omitempty"`
}

func (a AppointmentDTO) ToEntity() (AppointmentEntity, error) {
	var entity AppointmentEntity
	var err error

	entity.ID, _ = primitive.ObjectIDFromHex(a.ID)

	entity.PatientID, err = primitive.ObjectIDFromHex(a.PatientID)
	if err != nil {
		return AppointmentEntity{}, err
	}

	entity.DoctorID, err = primitive.ObjectIDFromHex(a.DoctorID)
	if err != nil {
		return AppointmentEntity{}, err
	}

	entity.Type = a.Type
	entity.DateTime = a.DateTime
	entity.Duration = a.Duration
	entity.Status = a.Status
	entity.Location = a.Location
	entity.Notes = a.Notes
	entity.PatientHistory = a.PatientHistory
	entity.CreatedAt = a.CreatedAt
	entity.UpdatedAt = a.UpdatedAt

	return entity, nil
}

func (a *AppointmentEntity) ToDTO() AppointmentDTO {
	return AppointmentDTO{
		ID:             a.ID.Hex(),
		PatientID:      a.PatientID.Hex(),
		DoctorID:       a.DoctorID.Hex(),
		Type:           a.Type,
		DateTime:       a.DateTime,
		Duration:       a.Duration,
		Status:         a.Status,
		Location:       a.Location,
		Notes:          a.Notes,
		PatientHistory: a.PatientHistory,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

func (a *AppointmentEntity) ToDetailDTO(patient *PatientEntity, lastRecord *MedicalRecordEntity) *AppointmentDetailResponse {
	detail := &AppointmentDetailResponse{
		Appointment: a.ToDTO(),
	}

	if patient != nil {
		patientDTO := patient.ToDTO()
		detail.Patient = &patientDTO
	}

	if lastRecord != nil {
		recordDTO := lastRecord.ToDTO()
		detail.LastRecord = &recordDTO
	}

	return detail
}
