package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Patient struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Age       int                `bson:"age" json:"age"`
	Gender    string             `bson:"gender" json:"gender"`
	Phone     string             `bson:"phone" json:"phone"`
	Address   string             `bson:"address" json:"address"`
	LastVisit time.Time          `bson:"lastVisit" json:"lastVisit"`
}

type Doctor struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name         string               `bson:"name" json:"name"`
	Specialty    string               `bson:"specialty" json:"specialty"`
	Availability []string             `bson:"AvailableDay" json:"availability"` // e.g. ["Monday", "Wednesday"]
	Email        string               `bson:"email" json:"email"`
	Phone        string               `bson:"phone" json:"phone"`
	Patients     []primitive.ObjectID `bson:"patients" json:"patients"`
}

type Appointment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PatientID      primitive.ObjectID `bson:"patientId" json:"patientId"`
	DoctorID       primitive.ObjectID `bson:"doctorId" json:"doctorId"`
	Date           time.Time          `bson:"date" json:"date"`
	Time           time.Time          `bson:"time" json:"time"`
	Status         string             `bson:"status" json:"status"` // scheduled, completed, cancelled
	Type           string             `bson:"type" json:"type"`
	Location       string             `bson:"location" json:"location"`
	Notes          string             `bson:"notes" json:"notes"`
	PatientHistory []string           `bson:"patientHistory" json:"patientHistory"`
}

type MedicalRecord struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PatientID    primitive.ObjectID `bson:"patient_id" json:"patient_id"`
	DoctorID     primitive.ObjectID `bson:"doctor_id" json:"doctor_id"`
	VisitDate    time.Time          `bson:"visit_date" json:"visit_date"`
	Diagnosis    string             `bson:"diagnosis" json:"diagnosis"`
	Prescription string             `bson:"prescription" json:"prescription"`
	Notes        string             `bson:"notes" json:"notes"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
