package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityType string

const (
	ActivityTypeAppointment ActivityType = "APPOINTMENT"
	ActivityTypeMedicalRecord ActivityType = "MEDICAL_RECORD"
	ActivityTypePatient       ActivityType = "PATIENT"
	ActivityTypeDoctor        ActivityType = "DOCTOR"
)

type ActivityEntity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Type        ActivityType       `bson:"type"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Timestamp   time.Time          `bson:"timestamp"`
}
