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

func (as AppointmentStatus) IsValid() bool {
	switch as {
	case AppointmentStatusScheduled, AppointmentStatusConfirmed, AppointmentStatusCompleted, AppointmentStatusCancelled:
		return true
	}
	return false
}

type AppointmentType string

const (
	AppointmentTypeCheckUp      AppointmentType = "check-up"
	AppointmentTypeFollowUp     AppointmentType = "follow-up"
	AppointmentTypeConsultation AppointmentType = "consultation"
	AppointmentTypeProcedure    AppointmentType = "procedure"
	AppointmentTypeEmergency    AppointmentType = "emergency"
)

func (at AppointmentType) IsValid() bool {
	switch at {
	case AppointmentTypeCheckUp, AppointmentTypeFollowUp, AppointmentTypeConsultation, AppointmentTypeProcedure, AppointmentTypeEmergency:
		return true
	}
	return false
}

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
	CreatedBy      primitive.ObjectID `bson:"createdBy"`
	UpdatedBy      primitive.ObjectID `bson:"updatedBy"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}

type AppointmentDTO struct {
	ID             string            `json:"id"`
	PatientID      string            `json:"patientId" validate:"required,mongodb"`
	DoctorID       string            `json:"doctorId" validate:"required,mongodb"`
	Type           AppointmentType   `json:"type" validate:"required,oneof=check-up follow-up consultation procedure emergency"`
	DateTime       time.Time         `json:"dateTime" validate:"required,datetime"`
	Duration       int               `json:"duration" validate:"required,gt=0"`
	Status         AppointmentStatus `json:"status" validate:"required,oneof=Scheduled Confirmed Completed Cancelled"`
	Location       string            `json:"location" validate:"required,min=3,max=100"`
	Notes          string            `json:"notes,omitempty" validate:"max=500"`
	PatientHistory string            `json:"patientHistory,omitempty" validate:"max=1000"`
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
