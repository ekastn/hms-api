package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Description	Patient object
// @swagger:model
type PatientEntity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"60d0fe4f53115a001f000001"`
	Name      string             `bson:"name" json:"name" example:"Jane Doe"`
	Age       int                `bson:"age" json:"age" example:30`
	Gender    string             `bson:"gender" json:"gender" example:"Female"`
	Phone     string             `bson:"phone" json:"phone" example:"1234567890"`
	Email     string             `bson:"email" json:"email" example:"jane.doe@example.com"`
	Address   string             `bson:"address" json:"address" example:"123 Main St"`
	LastVisit time.Time          `bson:"lastVisit" json:"lastVisit" example:"2025-07-17T10:00:00Z"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy,omitempty"`
	IsDeleted bool               `bson:"isDeleted" json:"isDeleted" example:false`
	DeletedAt *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// @Description	Patient data transfer object
// @swagger:model
type PatientDTO struct {
	ID        string    `json:"id" example:"60d0fe4f53115a001f000001"`
	Name      string    `json:"name" example:"Jane Doe"`
	Age       int       `json:"age" example:30`
	Gender    string    `json:"gender" example:"Female"`
	Phone     string    `json:"phone" example:"1234567890"`
	Email     string    `json:"email" example:"jane.doe@example.com"`
	Address   string    `json:"address" example:"123 Main St"`
	LastVisit time.Time `json:"lastVisit" example:"2025-07-17T10:00:00Z"`
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

// @Description	Request body for creating a new patient
// @swagger:model
type CreatePatientRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=100" example:"John Doe"`
	Age     int    `json:"age" validate:"required,gt=0,lte=120" example:45`
	Gender  string `json:"gender" validate:"required,oneof=Male Female Other" example:"Male"`
	Phone   string `json:"phone" validate:"required,e164" example:"+1234567890"`
	Email   string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Address string `json:"address" validate:"required,min=10,max=200" example:"123 Main St, Anytown, USA"`
}

// @Description	Request body for updating an existing patient
// @swagger:model
type UpdatePatientRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=100" example:"John Doe Jr."`
	Age     int    `json:"age" validate:"required,gt=0,lte=120" example:46`
	Gender  string `json:"gender" validate:"required,oneof=Male Female Other" example:"Male"`
	Phone   string `json:"phone" validate:"required,e164" example:"+1234567890"`
	Email   string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Address string `json:"address" validate:"required,min=10,max=200" example:"456 Oak Ave, Somewhere, USA"`
}

// @Description	Detailed patient information including recent appointments and medical history
// @swagger:model
type PatientDetailResponse struct {
	Patient            PatientDTO         `json:"patient"`
	RecentAppointments []AppointmentDTO   `json:"recentAppointments"`
	MedicalHistory     []MedicalRecordDTO `json:"medicalHistory"`
}
