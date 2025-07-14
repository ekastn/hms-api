package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PatientEntity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Age       int                `bson:"age"`
	Gender    string             `bson:"gender"`
	Phone     string             `bson:"phone"`
	Email     string             `bson:"email"`
	Address   string             `bson:"address"`
	LastVisit time.Time          `bson:"lastVisit"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type PatientDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	LastVisit time.Time `json:"lastVisit"`
}

func (p PatientDTO) ToEntity() PatientEntity {
	id, _ := primitive.ObjectIDFromHex(p.ID)
	return PatientEntity{
		ID:        id,
		Name:      p.Name,
		Age:       p.Age,
		Gender:    p.Gender,
		Phone:     p.Phone,
		Address:   p.Address,
		LastVisit: p.LastVisit,
	}
}

func (p PatientEntity) ToDTO() PatientDTO {
	return PatientDTO{
		ID:        p.ID.Hex(),
		Name:      p.Name,
		Age:       p.Age,
		Gender:    p.Gender,
		Phone:     p.Phone,
		Address:   p.Address,
		LastVisit: p.LastVisit,
	}
}

type CreatePatientRequest struct {
	Name    string `json:"name" validate:"required"`
	Age     int    `json:"age" validate:"required,gt=0"`
	Gender  string `json:"gender" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Address string `json:"address" validate:"required"`
}

type UpdatePatientRequest struct {
	Name    string `json:"name" validate:"required"`
	Age     int    `json:"age" validate:"required,gt=0"`
	Gender  string `json:"gender" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Address string `json:"address" validate:"required"`
}

type PatientDetailResponse struct {
	Patient            PatientDTO         `json:"patient"`
	RecentAppointments []AppointmentDTO   `json:"recentAppointments"`
	MedicalHistory     []MedicalRecordDTO `json:"medicalHistory"`
}
