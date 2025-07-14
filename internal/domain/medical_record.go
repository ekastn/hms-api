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

type MedicalRecordEntity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PatientID   primitive.ObjectID `bson:"patientId"`
	DoctorID    primitive.ObjectID `bson:"doctorId"`
	Date        time.Time          `bson:"date"`
	RecordType  MedicalRecordType  `bson:"recordType"`
	Description string             `bson:"description"`
	Diagnosis   string             `bson:"diagnosis"`
	Treatment   string             `bson:"treatment"`
	Notes       string             `bson:"notes,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

type MedicalRecordDTO struct {
	ID          string    `json:"id"`
	PatientID   string    `json:"patientId"`
	DoctorID    string    `json:"doctorId"`
	Date        string    `json:"date"`
	RecordType  string    `json:"recordType"`
	Description string    `json:"description"`
	Diagnosis   string    `json:"diagnosis"`
	Treatment   string    `json:"treatment"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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

type CreateMedicalRecordRequest struct {
	PatientID   string `json:"patientId" validate:"required"`
	DoctorID    string `json:"doctorId" validate:"required"`
	RecordType  string `json:"recordType" validate:"required"`
	Description string `json:"description" validate:"required"`
	Diagnosis   string `json:"diagnosis" validate:"required"`
	Treatment   string `json:"treatment" validate:"required"`
	Notes       string `json:"notes,omitempty"`
}

type UpdateMedicalRecordRequest struct {
	RecordType  string `json:"recordType,omitempty"`
	Description string `json:"description,omitempty"`
	Diagnosis   string `json:"diagnosis,omitempty"`
	Treatment   string `json:"treatment,omitempty"`
	Notes       string `json:"notes,omitempty"`
}
