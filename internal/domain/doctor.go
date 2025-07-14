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
	DayOfWeek int       `bson:"dayOfWeek" json:"dayOfWeek"` // 0-6 (Sunday-Saturday)
	StartTime time.Time `bson:"startTime" json:"startTime"`
	EndTime   time.Time `bson:"endTime" json:"endTime"`
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
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Specialty string `json:"specialty" validate:"required,min=3,max=100"`
	Phone     string `json:"phone" validate:"required,e164"`
	Email     string `json:"email" validate:"required,email"`
}

type UpdateDoctorRequet struct {
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Specialty string `json:"specialty" validate:"required,min=3,max=100"`
	Phone     string `json:"phone" validate:"required,e164"`
	Email     string `json:"email" validate:"required,email"`
}

type DoctorDetailResponse struct {
	Doctor         *DoctorDTO   `json:"doctor"`
	RecentPatients []PatientDTO `json:"recentPatients,omitempty"`
}

func (d *DoctorEntity) ToDetailDTO(recentPatients []*PatientEntity) *DoctorDetailResponse {
	dto := d.ToDTO()
	detail := &DoctorDetailResponse{
		Doctor: &dto,
	}

	if len(recentPatients) > 0 {
		detail.RecentPatients = make([]PatientDTO, 0, len(recentPatients))
		for _, p := range recentPatients {
			if p != nil {
				detail.RecentPatients = append(detail.RecentPatients, p.ToDTO())
			}
		}
	}

	return detail
}
