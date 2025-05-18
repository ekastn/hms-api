package seed

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Seeder struct {
	db *mongo.Database
}

func NewSeeder(db *mongo.Database) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	log.Println("Starting database seeding...")

	// Clear existing data
	if err := s.clearCollections(ctx); err != nil {
		return fmt.Errorf("failed to clear collections: %w", err)
	}

	// Seed data
	if err := s.SeedDoctors(ctx); err != nil {
		return fmt.Errorf("failed to seed doctors: %w", err)
	}

	if err := s.SeedPatients(ctx); err != nil {
		return fmt.Errorf("failed to seed patients: %w", err)
	}

	if err := s.SeedAppointments(ctx); err != nil {
		return fmt.Errorf("failed to seed appointments: %w", err)
	}

	if err := s.SeedMedicalRecords(ctx); err != nil {
		return fmt.Errorf("failed to seed medical records: %w", err)
	}

	log.Println("Database seeding completed successfully")
	return nil
}

func (s *Seeder) clearCollections(ctx context.Context) error {
	collections := []string{"doctors", "patients", "appointments", "medical_records"}
	for _, coll := range collections {
		if err := s.db.Collection(coll).Drop(ctx); err != nil {
			return fmt.Errorf("failed to drop collection %s: %w", coll, err)
		}
	}
	return nil
}

func (s *Seeder) SeedDoctors(ctx context.Context) error {
	repo := repository.NewDoctorRepository(s.db.Collection("doctors"))

	doctors := []*domain.DoctorEntity{
		{
			Name:      "Dr. John Smith",
			Specialty: "Cardiology",
			Phone:     "+1234567890",
			Email:     "john.smith@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 1, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
				{DayOfWeek: 3, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Sarah Johnson",
			Specialty: "Pediatrics",
			Phone:     "+1987654321",
			Email:     "sarah.j@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 2, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
				{DayOfWeek: 4, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Michael Brown",
			Specialty: "Neurology",
			Phone:     "+1555123456",
			Email:     "michael.b@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 1, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
				{DayOfWeek: 4, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Emily Davis",
			Specialty: "Dermatology",
			Phone:     "+1555234567",
			Email:     "emily.d@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 2, StartTime: time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 30, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 30, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. James Wilson",
			Specialty: "Orthopedics",
			Phone:     "+1555345678",
			Email:     "james.w@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 1, StartTime: time.Date(0, 1, 1, 8, 30, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 30, 0, 0, time.UTC)},
				{DayOfWeek: 3, StartTime: time.Date(0, 1, 1, 8, 30, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 30, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Jennifer Lee",
			Specialty: "Ophthalmology",
			Phone:     "+1555456789",
			Email:     "jennifer.l@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 2, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
				{DayOfWeek: 4, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Robert Taylor",
			Specialty: "Cardiology",
			Phone:     "+1555567890",
			Email:     "robert.t@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 3, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Lisa Garcia",
			Specialty: "Pediatrics",
			Phone:     "+1555678901",
			Email:     "lisa.g@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 1, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
				{DayOfWeek: 4, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. David Martinez",
			Specialty: "Neurology",
			Phone:     "+1555789012",
			Email:     "david.m@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 2, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Karen White",
			Specialty: "Dermatology",
			Phone:     "+1555890123",
			Email:     "karen.w@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 3, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Thomas Anderson",
			Specialty: "Orthopedics",
			Phone:     "+1555901234",
			Email:     "thomas.a@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 1, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
				{DayOfWeek: 4, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Patricia Clark",
			Specialty: "Ophthalmology",
			Phone:     "+1555012345",
			Email:     "patricia.c@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 2, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Richard Lewis",
			Specialty: "Cardiology",
			Phone:     "+1555123457",
			Email:     "richard.l@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 1, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
				{DayOfWeek: 3, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Susan Walker",
			Specialty: "Pediatrics",
			Phone:     "+1555234568",
			Email:     "susan.w@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 2, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
				{DayOfWeek: 4, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Charles Young",
			Specialty: "Neurology",
			Phone:     "+1555345679",
			Email:     "charles.y@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 1, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Jessica Hall",
			Specialty: "Dermatology",
			Phone:     "+1555456790",
			Email:     "jessica.h@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 3, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Daniel Allen",
			Specialty: "Orthopedics",
			Phone:     "+1555568901",
			Email:     "daniel.a@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 2, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
				{DayOfWeek: 4, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Nancy Scott",
			Specialty: "Ophthalmology",
			Phone:     "+1555679012",
			Email:     "nancy.s@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 1, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
				{DayOfWeek: 3, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Kevin King",
			Specialty: "Cardiology",
			Phone:     "+1555780123",
			Email:     "kevin.k@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 2, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Michelle Wright",
			Specialty: "Pediatrics",
			Phone:     "+1555891234",
			Email:     "michelle.w@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 1, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
				{DayOfWeek: 4, StartTime: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Steven Lopez",
			Specialty: "Neurology",
			Phone:     "+1555902345",
			Email:     "steven.l@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 2, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)},
			},
		},
		{
			Name:      "Dr. Laura Hill",
			Specialty: "Dermatology",
			Phone:     "+1555013456",
			Email:     "laura.h@example.com",
			Availability: []domain.TimeSlot{
				{DayOfWeek: 3, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
				{DayOfWeek: 5, StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), EndTime: time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)},
			},
		},
	}

	for _, doc := range doctors {
		doc.CreatedAt = time.Now()
		doc.UpdatedAt = time.Now()
		_, err := repo.Create(ctx, doc)
		if err != nil {
			return fmt.Errorf("failed to create doctor %s: %w", doc.Name, err)
		}
	}

	log.Printf("Seeded %d doctors\n", len(doctors))
	return nil
}

func (s *Seeder) SeedPatients(ctx context.Context) error {
	repo := repository.NewPatientRepository(s.db.Collection("patients"))

	// Generate sample patients
	patients := []*domain.PatientEntity{
		{
			Name:      "Alice Johnson",
			Age:       38,
			Gender:    "Female",
			Phone:     "+1555123456",
			Email:     "alice@example.com",
			Address:   "123 Main St, Anytown, USA",
			LastVisit: time.Now().Add(-7 * 24 * time.Hour), // 1 week ago
		},
		{
			Name:      "Bob Williams",
			Age:       45,
			Gender:    "Male",
			Phone:     "+1555654321",
			Email:     "bob@example.com",
			Address:   "456 Oak Ave, Somewhere, USA",
			LastVisit: time.Now().Add(-14 * 24 * time.Hour), // 2 weeks ago
		},
		{
			Name:      "Carol Davis",
			Age:       29,
			Gender:    "Female",
			Phone:     "+1555123789",
			Email:     "carol.d@example.com",
			Address:   "789 Pine St, Anytown, USA",
			LastVisit: time.Now().Add(-3 * 24 * time.Hour), // 3 days ago
		},
		{
			Name:      "David Miller",
			Age:       52,
			Gender:    "Male",
			Phone:     "+1555234567",
			Email:     "david.m@example.com",
			Address:   "321 Elm St, Somewhere, USA",
			LastVisit: time.Now().Add(-30 * 24 * time.Hour), // 1 month ago
		},
		{
			Name:      "Eve Wilson",
			Age:       34,
			Gender:    "Female",
			Phone:     "+1555345678",
			Email:     "eve.w@example.com",
			Address:   "654 Maple Ave, Anytown, USA",
			LastVisit: time.Now().Add(-2 * 24 * time.Hour), // 2 days ago
		},
		{
			Name:      "Frank Moore",
			Age:       61,
			Gender:    "Male",
			Phone:     "+1555456789",
			Email:     "frank.m@example.com",
			Address:   "987 Cedar Ln, Somewhere, USA",
			LastVisit: time.Now().Add(-45 * 24 * time.Hour), // 45 days ago
		},
		{
			Name:      "Grace Taylor",
			Age:       27,
			Gender:    "Female",
			Phone:     "+1555567890",
			Email:     "grace.t@example.com",
			Address:   "159 Oak Dr, Anytown, USA",
			LastVisit: time.Now().Add(-5 * 24 * time.Hour), // 5 days ago
		},
		{
			Name:      "Henry Brown",
			Age:       42,
			Gender:    "Male",
			Phone:     "+1555678901",
			Email:     "henry.b@example.com",
			Address:   "753 Pine Ave, Somewhere, USA",
			LastVisit: time.Now().Add(-21 * 24 * time.Hour), // 3 weeks ago
		},
		{
			Name:      "Ivy Garcia",
			Age:       31,
			Gender:    "Female",
			Phone:     "+1555789012",
			Email:     "ivy.g@example.com",
			Address:   "246 Elm St, Anytown, USA",
			LastVisit: time.Now().Add(-10 * 24 * time.Hour), // 10 days ago
		},
		{
			Name:      "Jack Martinez",
			Age:       55,
			Gender:    "Male",
			Phone:     "+1555890123",
			Email:     "jack.m@example.com",
			Address:   "864 Maple Dr, Somewhere, USA",
			LastVisit: time.Now().Add(-60 * 24 * time.Hour), // 2 months ago
		},
		{
			Name:      "Karen Anderson",
			Age:       39,
			Gender:    "Female",
			Phone:     "+1555901234",
			Email:     "karen.a@example.com",
			Address:   "975 Oak Ln, Anytown, USA",
			LastVisit: time.Now().Add(-1 * 24 * time.Hour), // 1 day ago
		},
		{
			Name:      "Liam Thomas",
			Age:       48,
			Gender:    "Male",
			Phone:     "+1555012345",
			Email:     "liam.t@example.com",
			Address:   "357 Pine St, Somewhere, USA",
			LastVisit: time.Now().Add(-28 * 24 * time.Hour), // 4 weeks ago
		},
		{
			Name:      "Mia Jackson",
			Age:       33,
			Gender:    "Female",
			Phone:     "+1555123457",
			Email:     "mia.j@example.com",
			Address:   "159 Elm Ave, Anytown, USA",
			LastVisit: time.Now().Add(-7 * 24 * time.Hour), // 1 week ago
		},
		{
			Name:      "Noah White",
			Age:       50,
			Gender:    "Male",
			Phone:     "+1555234568",
			Email:     "noah.w@example.com",
			Address:   "753 Maple St, Somewhere, USA",
			LastVisit: time.Now().Add(-90 * 24 * time.Hour), // 3 months ago
		},
		{
			Name:      "Olivia Harris",
			Age:       36,
			Gender:    "Female",
			Phone:     "+1555345679",
			Email:     "olivia.h@example.com",
			Address:   "246 Oak Dr, Anytown, USA",
			LastVisit: time.Now().Add(-2 * 24 * time.Hour), // 2 days ago
		},
		{
			Name:      "Peter Martin",
			Age:       41,
			Gender:    "Male",
			Phone:     "+1555456790",
			Email:     "peter.m@example.com",
			Address:   "864 Pine Ln, Somewhere, USA",
			LastVisit: time.Now().Add(-14 * 24 * time.Hour), // 2 weeks ago
		},
		{
			Name:      "Quinn Thompson",
			Age:       30,
			Gender:    "Non-binary",
			Phone:     "+1555568901",
			Email:     "quinn.t@example.com",
			Address:   "975 Elm St, Anytown, USA",
			LastVisit: time.Now().Add(-5 * 24 * time.Hour), // 5 days ago
		},
		{
			Name:      "Rachel Garcia",
			Age:       47,
			Gender:    "Female",
			Phone:     "+1555679012",
			Email:     "rachel.g@example.com",
			Address:   "357 Maple Ave, Somewhere, USA",
			LastVisit: time.Now().Add(-21 * 24 * time.Hour), // 3 weeks ago
		},
		{
			Name:      "Samuel Rodriguez",
			Age:       32,
			Gender:    "Male",
			Phone:     "+1555780123",
			Email:     "samuel.r@example.com",
			Address:   "159 Pine Dr, Anytown, USA",
			LastVisit: time.Now().Add(-1 * 24 * time.Hour), // 1 day ago
		},
		{
			Name:      "Tina Lewis",
			Age:       53,
			Gender:    "Female",
			Phone:     "+1555891234",
			Email:     "tina.l@example.com",
			Address:   "753 Oak Ln, Somewhere, USA",
			LastVisit: time.Now().Add(-60 * 24 * time.Hour), // 2 months ago
		},
	}

	for _, p := range patients {
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		_, err := repo.Create(ctx, p)
		if err != nil {
			return fmt.Errorf("failed to create patient %s: %w", p.Name, err)
		}
	}

	log.Printf("Seeded %d patients\n", len(patients))
	return nil
}

func (s *Seeder) SeedAppointments(ctx context.Context) error {
	apptRepo := repository.NewAppointmentRepository(s.db.Collection("appointments"))
	doctorRepo := repository.NewDoctorRepository(s.db.Collection("doctors"))
	patientRepo := repository.NewPatientRepository(s.db.Collection("patients"))

	// Get all doctors and patients
	doctors, err := doctorRepo.GetAll(ctx)
	if err != nil || len(doctors) == 0 {
		return fmt.Errorf("no doctors found: %w", err)
	}

	patients, err := patientRepo.GetAll(ctx)
	if err != nil || len(patients) == 0 {
		return fmt.Errorf("no patients found: %w", err)
	}

	now := time.Now()
	appointmentTypes := []domain.AppointmentType{
		domain.AppointmentTypeCheckUp,
		domain.AppointmentTypeFollowUp,
		domain.AppointmentTypeConsultation,
		domain.AppointmentTypeEmergency,
		domain.AppointmentTypeProcedure,
	}

	statuses := []domain.AppointmentStatus{
		domain.AppointmentStatusScheduled,
		domain.AppointmentStatusConfirmed,
		domain.AppointmentStatusCompleted,
		domain.AppointmentStatusCancelled,
	}

	locations := []string{
		"Room 101", "Room 102", "Room 103", "Room 201", "Room 202", "Room 203",
		"Room 301", "Room 302", "Exam Room A", "Exam Room B", "Exam Room C",
	}

	appointments := make([]*domain.AppointmentEntity, 0, 50) // Increased to 50 appointments

	// Create more meaningful appointment patterns
	for i := 0; i < 50; i++ {
		// Distribute patients and doctors more evenly
		patient := patients[i%len(patients)]
		// Try to match patients with doctors of appropriate specialties
		// Get a doctor, trying to match specialty when possible
		var doctor *domain.DoctorEntity
		if i%3 == 0 {
			// For every 3rd appointment, try to match with a doctor of matching specialty
			for _, d := range doctors {
				if (d.Specialty == "Cardiology" && strings.Contains(patient.Name, "Heart")) ||
					(d.Specialty == "Pediatrics" && patient.Age < 18) ||
					(d.Specialty == "Neurology" && i%7 == 0) ||
					(d.Specialty == "Dermatology" && i%5 == 0) {
					doctor = d
					break
				}
			}
		}
		// If no specialty match or not trying to match, just pick one
		if doctor == nil {
			doctor = doctors[i%len(doctors)]
		}

		// Create appointments with more realistic distribution
		daysInFuture := i % 30
		hour := 8 + (i % 9) // 8 AM to 5 PM
		minute := 0
		switch i % 4 {
		case 0:
			minute = 0
		case 1:
			minute = 15
		case 2:
			minute = 30
		case 3:
			minute = 45
		}

		// Create some recurring appointments
		if i%5 == 0 && i > 0 {
			daysInFuture = (i / 5) % 30 // Group recurring appointments
		}

		// Create some urgent appointments
		isUrgent := i%7 == 0
		if isUrgent {
			daysInFuture = i % 3 // Within next 3 days for urgent
		}

		// Create appointment
		appointment := &domain.AppointmentEntity{
			PatientID: patient.ID,
			DoctorID:  doctor.ID,
			Type:      appointmentTypes[i%len(appointmentTypes)],
			DateTime:  now.Add(time.Duration(daysInFuture)*24*time.Hour + time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute),
			Duration:  30 * (1 + (i % 4)), // 30, 60, 90, or 120 minutes
			Status:    statuses[i%len(statuses)],
			Location:  locations[i%len(locations)],
			Notes: fmt.Sprintf("Appointment for %s regarding %s. %s",
				patient.Name,
				[]string{"routine checkup", "follow-up", "new issue", "consultation"}[i%4],
				map[bool]string{true: "URGENT - " + strings.ToUpper(
					[]string{"Hypertension", "Type 2 Diabetes", "Upper Respiratory Infection", "Migraine"}[i%4]),
					false: ""}[isUrgent]),
			PatientHistory: fmt.Sprintf("Previous history: %s. %s",
				[]string{"Hypertension", "Type 2 Diabetes", "Upper Respiratory Infection", "Migraine"}[i%4],
				[]string{"No known allergies.", "Allergic to penicillin.", "Allergic to NSAIDs.", "No significant history."}[i%4]),
		}

		appointments = append(appointments, appointment)
	}

	for _, appt := range appointments {
		appt.CreatedAt = time.Now()
		appt.UpdatedAt = time.Now()
		_, err := apptRepo.Create(ctx, appt)
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}
	}

	log.Printf("Seeded %d appointments\n", len(appointments))
	return nil
}

func (s *Seeder) SeedMedicalRecords(ctx context.Context) error {
	repo := repository.NewMedicalRecordRepository(s.db.Collection("medical_records"))
	patientRepo := repository.NewPatientRepository(s.db.Collection("patients"))
	doctorRepo := repository.NewDoctorRepository(s.db.Collection("doctors"))

	// Get all patients and doctors
	patients, err := patientRepo.GetAll(ctx)
	if err != nil || len(patients) == 0 {
		return fmt.Errorf("no patients found: %w", err)
	}

	doctors, err := doctorRepo.GetAll(ctx)
	if err != nil || len(doctors) == 0 {
		return fmt.Errorf("no doctors found: %w", err)
	}

	// Define record types for medical records
	recordTypes := []domain.MedicalRecordType{
		domain.RecordTypeCheckUp,
		domain.RecordTypeFollowUp,
		domain.RecordTypeEmergency,
		domain.RecordTypeProcedure,
	}

	// Define diagnoses for medical records
	diagnoses := []string{
		"Hypertension",
		"Type 2 Diabetes",
		"Upper Respiratory Infection",
		"Migraine",
		"Allergic Rhinitis",
		"Gastroenteritis",
		"Acute Bronchitis",
		"Urinary Tract Infection",
		"Back Pain",
		"Anxiety Disorder",
	}

	treatments := []string{
		"Prescribed medication",
		"Referred to specialist",
		"Lab tests ordered",
		"Physical therapy recommended",
		"Lifestyle changes suggested",
		"Follow-up in 2 weeks",
		"No treatment needed",
		"Vaccination administered",
	}

	notes := []string{
		"Patient responded well to treatment",
		"Patient to return if symptoms worsen",
		"Patient advised to monitor symptoms",
		"Patient educated about condition",
		"No immediate concerns noted",
		"Vital signs stable",
		"Patient reports improvement",
		"Plan to reassess at next visit",
	}

	// Define descriptions for medical records (used in the code)
	descriptions := []string{
		"Routine checkup and consultation",
		"Follow-up for previous condition",
		"Annual physical examination",
		"Emergency treatment provided",
		"Review of lab results",
		"Discussion of imaging findings",
		"Specialist consultation",
		"Post-operative follow-up",
	}

	records := make([]*domain.MedicalRecordEntity, 0, 50) // Increased to 50 records

	// Add more comprehensive medical records
	for i := 0; i < 50; i++ {
		patient := patients[i%len(patients)]
		doctor := doctors[i%len(doctors)]

		// Create records from last 2 years with more recent records being more frequent
		daysAgo := 0
		switch {
		case i < 15: // 30% recent (last 30 days)
			daysAgo = i % 30
		case i < 35: // 40% somewhat recent (1-6 months)
			daysAgo = 30 + (i % 150)
		default: // 30% older (6-24 months)
			daysAgo = 180 + (i % 540)
		}

		// Create more detailed and realistic records
		recordDate := time.Now().AddDate(0, 0, -daysAgo)
		recordType := recordTypes[i%len(recordTypes)]
		diagnosis := diagnoses[i%len(diagnoses)]
		treatment := treatments[i%len(treatments)]
		description := descriptions[i%len(descriptions)]

		// Create more detailed notes based on diagnosis and treatment
		detailedNotes := fmt.Sprintf("%s\nVital signs: BP %d/%d, Pulse %d, Temp %.1f°C\nAssessment: %s\nPlan: %s\nDescription: %s",
			notes[i%len(notes)],
			100+(i%40), 60+(i%30), // BP
			60+(i%40),               // Pulse
			36.5+(float64(i%30)/10), // Temp
			diagnosis,
			treatment,
			description,
		)

		record := &domain.MedicalRecordEntity{
			PatientID:   patient.ID,
			DoctorID:    doctor.ID,
			Date:        recordDate,
			RecordType:  recordType,
			Description: fmt.Sprintf("%s - %s", recordType, diagnosis),
			Diagnosis:   diagnosis,
			Treatment:   treatment,
			Notes:       detailedNotes,
			CreatedAt:   recordDate,
			UpdatedAt:   recordDate,
		}

		// Add follow-up records for some patients
		if i%5 == 0 && i > 0 {
			followUpDate := recordDate.AddDate(0, 0, 14+(i%14)) // 2-4 weeks later
			followUp := &domain.MedicalRecordEntity{
				PatientID:   patient.ID,
				DoctorID:    doctor.ID,
				Date:        followUpDate,
				RecordType:  domain.RecordTypeFollowUp,
				Description: fmt.Sprintf("Follow-up for %s", diagnosis),
				Diagnosis:   diagnosis,
				Treatment: fmt.Sprintf("Continue treatment. %s", []string{
					"Patient reports improvement.",
					"Symptoms persisting, adjusted medication.",
					"Referred to specialist for further evaluation.",
				}[i%3]),
				Notes:     fmt.Sprintf("Follow-up visit. %s", detailedNotes),
				CreatedAt: followUpDate,
				UpdatedAt: followUpDate,
			}
			records = append(records, followUp)
		}

		records = append(records, record)
	}

	// Save all records
	for _, record := range records {
		_, err := repo.Create(ctx, record)
		if err != nil {
			return fmt.Errorf("failed to create medical record: %w", err)
		}
	}

	log.Printf("Seeded %d medical records\n", len(records))
	return nil
}
