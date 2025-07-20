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

// @Description	Appointment object
// @swagger:model
type AppointmentEntity struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"60d0fe4f53115a001f000001"`
	PatientID      primitive.ObjectID `bson:"patientId" json:"patientId" example:"60d0fe4f53115a001f000002"`
	DoctorID       primitive.ObjectID `bson:"doctorId" json:"doctorId" example:"60d0fe4f53115a001f000003"`
	Type           AppointmentType    `bson:"type" json:"type" example:"check-up"`
	DateTime       time.Time          `bson:"dateTime" json:"dateTime" example:"2025-07-17T10:00:00Z"`
	Duration       int                `bson:"duration" json:"duration" example:30` // in minutes
	Status         AppointmentStatus  `bson:"status" json:"status" example:"Scheduled"`
	Location       string             `bson:"location" json:"location" example:"Room 101"`
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty" example:"Patient complained of headache"`
	PatientHistory string             `bson:"patientHistory,omitempty" json:"patientHistory,omitempty" example:"No significant medical history"`
	CreatedBy      primitive.ObjectID `bson:"createdBy" json:"createdBy,omitempty"`
	UpdatedBy      primitive.ObjectID `bson:"updatedBy" json:"updatedBy,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type AppointmentDTO struct {
	ID             string            `json:"id" example:"60d0fe4f53115a001f000001"`
	PatientID      string            `json:"patientId" validate:"required,mongodb" example:"60d0fe4f53115a001f000002"`
	DoctorID       string            `json:"doctorId" validate:"required,mongodb" example:"60d0fe4f53115a001f000003"`
	Type           AppointmentType   `json:"type" validate:"required,oneof=check-up follow-up consultation procedure emergency" example:"check-up"`
	DateTime       time.Time         `json:"dateTime" example:"2025-07-17T10:00:00Z"`
	Duration       int               `json:"duration" validate:"required,gt=0" example:30`
	Status         AppointmentStatus `json:"status" validate:"required,oneof=Scheduled Confirmed Completed Cancelled" example:"Scheduled"`
	Location       string            `json:"location" validate:"required,min=3,max=100" example:"Room 101"`
	Notes          string            `json:"notes,omitempty" validate:"max=500" example:"Patient complained of headache"`
	PatientHistory string            `json:"patientHistory,omitempty" validate:"max=1000" example:"No significant medical history"`
	CreatedAt      time.Time         `json:"createdAt" example:"2025-07-17T09:00:00Z"`
	UpdatedAt      time.Time         `json:"updatedAt" example:"2025-07-17T09:00:00Z"`
}

// @Description	Detailed appointment information
// @swagger:model
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

// @Description	Request body for updating an existing appointment
// @swagger:model
type UpdateAppointmentRequest struct {
	Type           *AppointmentType   `json:"type,omitempty" validate:"oneof=check-up follow-up consultation procedure emergency" example:"check-up"`
	DateTime       *time.Time         `json:"dateTime,omitempty" example:"2025-07-17T10:00:00Z"`
	Duration       *int               `json:"duration,omitempty" validate:"gt=0" example:30`
	Status         *AppointmentStatus `json:"status,omitempty" validate:"oneof=Scheduled Confirmed Completed Cancelled" example:"Scheduled"`
	Location       *string            `json:"location,omitempty" validate:"min=3,max=100" example:"Room 101"`
	Notes          *string            `json:"notes,omitempty" validate:"max=500" example:"Patient complained of headache"`
	PatientHistory *string            `json:"patientHistory,omitempty" example:"No significant medical history"`
}

// ApplyUpdates applies non-nil fields from UpdateAppointmentRequest to an AppointmentEntity.
// It returns true if any field was updated, false otherwise.
func (req *UpdateAppointmentRequest) ApplyUpdates(entity *AppointmentEntity) bool {
	updated := false


	if req.Type != nil {
		entity.Type = *req.Type
		updated = true
	}
	if req.DateTime != nil {
		entity.DateTime = *req.DateTime
		updated = true
	}
	if req.Duration != nil {
		entity.Duration = *req.Duration
		updated = true
	}
	if req.Status != nil {
		entity.Status = *req.Status
		updated = true
	}
	if req.Location != nil {
		entity.Location = *req.Location
		updated = true
	}
	if req.Notes != nil {
		entity.Notes = *req.Notes
		updated = true
	}
	if req.PatientHistory != nil {
		entity.PatientHistory = *req.PatientHistory
		updated = true
	}


	return updated
}

// @Description Request body for updating appointment status
// @swagger:model
type UpdateAppointmentStatusRequest struct {
    Status AppointmentStatus `json:"status" validate:"required,oneof=Scheduled Confirmed Completed Cancelled"`
}

// @Description Request body for creating a new appointment
// @swagger:model
type CreateAppointmentRequest struct {
	PatientID      string          `json:"patientId" validate:"required,mongodb" example:"60d0fe4f53115a001f000002"`
	DoctorID       string          `json:"doctorId" validate:"required,mongodb" example:"60d0fe4f53115a001f000003"`
	Type           AppointmentType `json:"type" validate:"required,oneof=check-up follow-up consultation procedure emergency" example:"check-up"`
	DateTime       time.Time       `json:"dateTime" validate:"-" example:"2025-07-17T10:00:00Z"`
	Duration       int             `json:"duration" validate:"required,gt=0" example:30`
	Location       string          `json:"location" validate:"required,min=3,max=100" example:"Room 101"`
	Notes          string          `json:"notes,omitempty" validate:"max=500" example:"Patient complained of headache"`
	PatientHistory string          `json:"patientHistory,omitempty" validate:"max=1000" example:"No significant medical history"`
}