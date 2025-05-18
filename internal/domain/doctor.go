package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorEntity struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Specialty    string             `bson:"specialty"`
	Phone        string             `bson:"phone"`
	Email        string             `bson:"email"`
	Availability []TimeSlot         `bson:"availability"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}

type TimeSlot struct {
	DayOfWeek int       `bson:"dayOfWeek"` // 0-6 (Sunday-Saturday)
	StartTime time.Time `bson:"startTime"`
	EndTime   time.Time `bson:"endTime"`
}

type DoctorDTO struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Specialty    string     `json:"specialty"`
	Phone        string     `json:"phone"`
	Email        string     `json:"email"`
	Availability []TimeSlot `json:"availability"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

func (d DoctorDTO) ToEntity() DoctorEntity {
	id, _ := primitive.ObjectIDFromHex(d.ID)
	return DoctorEntity{
		ID:           id,
		Name:         d.Name,
		Specialty:    d.Specialty,
		Phone:        d.Phone,
		Email:        d.Email,
		Availability: d.Availability,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func (d DoctorEntity) ToDTO() DoctorDTO {
	return DoctorDTO{
		ID:           d.ID.Hex(),
		Name:         d.Name,
		Specialty:    d.Specialty,
		Phone:        d.Phone,
		Email:        d.Email,
		Availability: d.Availability,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

type CreateDoctorRequet struct {
	Name      string `json:"name" validate:"required"`
	Specialty string `json:"specialty" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Email     string `json:"email" validate:"required"`
}

type UpdateDoctorRequet struct {
	Name      string `json:"name" validate:"required"`
	Specialty string `json:"specialty" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Email     string `json:"email" validate:"required"`
}
