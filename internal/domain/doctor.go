package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Description	Doctor object
// @swagger:model
type DoctorEntity struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"60d0fe4f53115a001f000001"`
	Name         string             `bson:"name" json:"name" example:"Dr. John Doe"`
	Specialty    string             `bson:"specialty" json:"specialty" example:"Cardiology"`
	Phone        string             `bson:"phone" json:"phone" example:"1234567890"`
	Email        string             `bson:"email" json:"email" example:"john.doe@example.com"`
	Availability []TimeSlot         `bson:"availability" json:"availability"`
	CreatedBy    primitive.ObjectID `bson:"createdBy" json:"createdBy,omitempty"`
	UpdatedBy    primitive.ObjectID `bson:"updatedBy" json:"updatedBy,omitempty"`
	IsDeleted    bool               `bson:"isDeleted" json:"isDeleted" example:false`
	DeletedAt    *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type TimeSlot struct {
	DayOfWeek int       `bson:"dayOfWeek" json:"dayOfWeek"` // 0-6 (Sunday-Saturday)
	StartTime time.Time `bson:"startTime" json:"startTime"`
	EndTime   time.Time `bson:"endTime" json:"endTime"`
}

// @Description	Doctor data transfer object
// @swagger:model
type DoctorDTO struct {
	ID           string     `json:"id" example:"60d0fe4f53115a001f000001"`
	Name         string     `json:"name" example:"Dr. John Doe"`
	Specialty    string     `json:"specialty" example:"Cardiology"`
	Phone        string     `json:"phone" example:"1234567890"`
	Email        string     `json:"email" example:"john.doe@example.com"`
	Availability []TimeSlot `json:"availability"`
	CreatedAt    time.Time  `json:"createdAt" example:"2025-07-17T09:00:00Z"`
	UpdatedAt    time.Time  `json:"updatedAt" example:"2025-07-17T09:00:00Z"`
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

// @Description	Request body for creating a new doctor
// @swagger:model
type CreateDoctorRequet struct {
	Name      string `json:"name" validate:"required,min=3,max=100" example:"Dr. Jane Smith"`
	Specialty string `json:"specialty" validate:"required,min=3,max=100" example:"Pediatrics"`
	Phone     string `json:"phone" validate:"required,e164" example:"+1987654321"`
	Email     string `json:"email" validate:"required,email" example:"jane.smith@example.com"`
}

// @Description	Request body for updating an existing doctor
// @swagger:model
type UpdateDoctorRequet struct {
	Name      string `json:"name" validate:"required,min=3,max=100" example:"Dr. Jane Smith-Doe"`
	Specialty string `json:"specialty" validate:"required,min=3,max=100" example:"Pediatrics"`
	Phone     string `json:"phone" validate:"required,e164" example:"+1987654321"`
	Email     string `json:"email" validate:"required,email" example:"jane.smith@example.com"`
}

// @Description	Detailed doctor information including recent patients
// @swagger:model
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
