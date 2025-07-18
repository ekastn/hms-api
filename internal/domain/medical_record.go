package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MedicalRecordType string

const (
	RecordTypeCheckUp   MedicalRecordType = "checkup"
	RecordTypeFollowUp  MedicalRecordType = "followup"
	RecordTypeProcedure MedicalRecordType = "procedure"
	RecordTypeEmergency MedicalRecordType = "emergency"
)

func (mrt MedicalRecordType) IsValid() bool {
	switch mrt {
	case RecordTypeCheckUp, RecordTypeFollowUp, RecordTypeProcedure, RecordTypeEmergency:
		return true
	}
	return false
}

// @Description	Medical record object
// @swagger:model
type MedicalRecordEntity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"60d0fe4f53115a001f000001"`
	PatientID   primitive.ObjectID `bson:"patientId" json:"patientId" example:"60d0fe4f53115a001f000002"`
	DoctorID    primitive.ObjectID `bson:"doctorId" json:"doctorId" example:"60d0fe4f53115a001f000003"`
	Date        time.Time          `bson:"date" json:"date" example:"2025-07-17T10:00:00Z"`
	RecordType  MedicalRecordType  `bson:"recordType" json:"recordType" example:"checkup"`
	Description string             `bson:"description" json:"description" example:"Patient presented with flu-like symptoms." `
	Diagnosis   string             `bson:"diagnosis" json:"diagnosis" example:"Influenza A"`
	Treatment   string             `bson:"treatment" json:"treatment" example:"Prescribed Tamiflu and rest." `
	Notes       string             `bson:"notes,omitempty" json:"notes,omitempty" example:"Advised patient to stay hydrated." `
	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy,omitempty"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy" json:"updatedBy,omitempty"`
	IsDeleted   bool               `bson:"isDeleted" json:"isDeleted" example:false`
	DeletedAt   *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// @Description	Medical record data transfer object
// @swagger:model
type MedicalRecordDTO struct {
	ID          string    `json:"id" example:"60d0fe4f53115a001f000001"`
	PatientID   string    `json:"patientId" example:"60d0fe4f53115a001f000002"`
	DoctorID    string    `json:"doctorId" example:"60d0fe4f53115a001f000003"`
	Date        string    `json:"date" example:"2025-07-17T10:00:00Z"`
	RecordType  string    `json:"recordType" example:"checkup"`
	Description string    `json:"description" example:"Patient presented with flu-like symptoms." `
	Diagnosis   string    `json:"diagnosis" example:"Influenza A"`
	Treatment   string    `json:"treatment" example:"Prescribed Tamiflu and rest." `
	Notes       string    `json:"notes,omitempty" example:"Advised patient to stay hydrated." `
	CreatedAt   time.Time `json:"createdAt" example:"2025-07-17T09:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2025-07-17T09:00:00Z"`
}

func (m MedicalRecordDTO) ToEntity() (MedicalRecordEntity, error) {
	var entity MedicalRecordEntity
	var err error

	entity.ID, _ = primitive.ObjectIDFromHex(m.ID)

	recordDate, err := time.Parse(time.RFC3339, m.Date)
	if err != nil {
		return MedicalRecordEntity{}, err
	}

	entity.PatientID, err = primitive.ObjectIDFromHex(m.PatientID)
	if err != nil {
		return MedicalRecordEntity{}, err
	}

	entity.DoctorID, err = primitive.ObjectIDFromHex(m.DoctorID)
	if err != nil {
		return MedicalRecordEntity{}, err
	}

	entity.Date = recordDate
	entity.RecordType = MedicalRecordType(m.RecordType)
	entity.Description = m.Description
	entity.Diagnosis = m.Diagnosis
	entity.Treatment = m.Treatment
	entity.Notes = m.Notes
	entity.CreatedAt = m.CreatedAt
	entity.UpdatedAt = m.UpdatedAt

	return entity, nil
}

func (m *MedicalRecordEntity) ToDTO() MedicalRecordDTO {
	return MedicalRecordDTO{
		ID:          m.ID.Hex(),
		PatientID:   m.PatientID.Hex(),
		DoctorID:    m.DoctorID.Hex(),
		Date:        m.Date.Format(time.RFC3339),
		RecordType:  string(m.RecordType),
		Description: m.Description,
		Diagnosis:   m.Diagnosis,
		Treatment:   m.Treatment,
		Notes:       m.Notes,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// @Description	Request body for creating a new medical record
// @swagger:model
type CreateMedicalRecordRequest struct {
	PatientID   string `json:"patientId" validate:"required,mongodb" example:"60d0fe4f53115a001f000002"`
	DoctorID    string `json:"doctorId" validate:"required,mongodb" example:"60d0fe4f53115a001f000003"`
	RecordType  string `json:"recordType" validate:"required,oneof=checkup followup procedure emergency" example:"checkup"`
	Description string `json:"description" validate:"required,min=10,max=1000" example:"Patient presented with flu-like symptoms." `
	Diagnosis   string `json:"diagnosis" validate:"required,min=5,max=200" example:"Influenza A"`
	Treatment   string `json:"treatment" validate:"required,min=5,max=1000" example:"Prescribed Tamiflu and rest." `
	Notes       string `json:"notes,omitempty" validate:"max=500" example:"Advised patient to stay hydrated." `
}

// @Description	Request body for updating an existing medical record
// @swagger:model
type UpdateMedicalRecordRequest struct {
	RecordType  string `json:"recordType,omitempty" validate:"oneof=checkup followup procedure emergency" example:"followup"`
	Description string `json:"description,omitempty" validate:"min=10,max=1000" example:"Patient's symptoms have improved." `
	Diagnosis   string `json:"diagnosis,omitempty" validate:"min=5,max=200" example:"Resolved Influenza A"`
	Treatment   string `json:"treatment,omitempty" validate:"min=5,max=1000" example:"Continue current medication." `
	Notes       string `json:"notes,omitempty" validate:"max=500" example:"Patient is recovering well." `
}
