package seed

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Seed struct {
	db                   *mongo.Database
	userService          *service.UserService
	doctorService        *service.DoctorService
	patientService       *service.PatientService
	appointmentService   *service.AppointmentService
	medicalRecordService *service.MedicalRecordService
}

func NewSeeder(
	db *mongo.Database,
	userService *service.UserService,
	doctorService *service.DoctorService,
	patientService *service.PatientService,
	appointmentService *service.AppointmentService,
	medicalRecordService *service.MedicalRecordService,
) *Seed {
	return &Seed{
		db:                   db,
		userService:          userService,
		doctorService:        doctorService,
		patientService:       patientService,
		appointmentService:   appointmentService,
		medicalRecordService: medicalRecordService,
	}
}

func (s *Seed) Seed(ctx context.Context) {
	log.Println("Starting to seed database...")

	users, err := s.seedUsers(ctx)
	if err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}
	log.Printf("Successfully seeded %d users.\n", len(users))

	doctors, err := s.seedDoctors(ctx, users[0].ID)
	if err != nil {
		log.Fatalf("failed to seed doctors: %v", err)
	}
	log.Printf("Successfully seeded %d doctors.\n", len(doctors))

	patients, err := s.seedPatients(ctx, users[0].ID)
	if err != nil {
		log.Fatalf("failed to seed patients: %v", err)
	}
	log.Printf("Successfully seeded %d patients.\n", len(patients))

	appointments, err := s.seedAppointments(ctx, users[0].ID, doctors, patients)
	if err != nil {
		log.Fatalf("failed to seed appointments: %v", err)
	}
	log.Printf("Successfully seeded %d appointments.\n", len(appointments))

	_, err = s.seedMedicalRecords(ctx, users[0].ID, doctors, appointments)
	if err != nil {
		log.Fatalf("failed to seed medical records: %v", err)
	}
	log.Printf("Successfully seeded medical records.\n")

	log.Println("Database seeding completed successfully!")
}

func (s *Seed) seedUsers(ctx context.Context) ([]*domain.UserDTO, error) {
	usersToCreate := []domain.CreateUserRequest{
		{Name: "Admin User", Email: "admin@hms.com", Password: "password123", Role: domain.RoleAdmin},
		{Name: "Dr. Budi Santoso", Email: "budi.santoso@hms.com", Password: "password123", Role: domain.RoleDoctor},
		{Name: "Siti Aminah", Email: "siti.aminah@hms.com", Password: "password123", Role: domain.RoleNurse},
		{Name: "Ayu Lestari", Email: "ayu.lestari@hms.com", Password: "password123", Role: domain.RoleReceptionist},
		{Name: "Rina Wijaya", Email: "rina.wijaya@hms.com", Password: "password123", Role: domain.RoleManagement},
	}

	var createdUsers []*domain.UserDTO
	for _, userReq := range usersToCreate {
		id, err := s.userService.CreateUser(ctx, &userReq)
		if err != nil {
			log.Printf("failed to create user %s: %v, skipping", userReq.Name, err)
			continue
		}
		user, err := s.userService.GetUserByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve created user %s: %w", userReq.Name, err)
		}
		createdUsers = append(createdUsers, user)
	}
	return createdUsers, nil
}

func (s *Seed) seedDoctors(ctx context.Context, creatorID string) ([]*domain.DoctorEntity, error) {
	creatorObjID, _ := primitive.ObjectIDFromHex(creatorID)

	doctorsToCreate := []struct {
		Name      string
		Specialty string
		Phone     string
		Email     string
	}{
		{"Dr. Budi Santoso", "Kardiologi", "+6281234567890", "budi.santoso@hms.com"},
		{"Dr. Citra Lestari", "Neurologi", "+6281234567891", "citra.lestari@hms.com"},
		{"Dr. Dedi Setiawan", "Pediatri", "+6281234567892", "dedi.setiawan@hms.com"},
		{"Dr. Eka Putri", "Ortopedi", "+6281234567893", "eka.putri@hms.com"},
		{"Dr. Fitriani", "Dermatologi", "+6281234567894", "fitriani@hms.com"},
	}

	var createdDoctors []*domain.DoctorEntity
	for _, docData := range doctorsToCreate {
		doctor := &domain.DoctorEntity{
			Name:      docData.Name,
			Specialty: docData.Specialty,
			Phone:     docData.Phone,
			Email:     docData.Email,
		}

		id, err := s.doctorService.Create(ctx, doctor, creatorObjID)
		if err != nil {
			log.Printf("failed to create doctor %s: %v, skipping", doctor.Name, err)
			continue
		}
		newDoctor, err := s.doctorService.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve created doctor %s: %w", doctor.Name, err)
		}
		createdDoctors = append(createdDoctors, newDoctor)
	}
	return createdDoctors, nil
}

func (s *Seed) seedPatients(ctx context.Context, creatorID string) ([]*domain.PatientEntity, error) {
	creatorObjID, _ := primitive.ObjectIDFromHex(creatorID)
	var createdPatients []*domain.PatientEntity

	patientsToCreate := []struct {
		Name    string
		Age     int
		Gender  string
		Phone   string
		Email   string
		Address string
	}{
		{"Andi Pratama", 35, "Male", "+6281111111111", "andi.pratama@example.com", "Jl. Merdeka No. 1, Jakarta"},
		{"Bunga Lestari", 28, "Female", "+6281222222222", "bunga.lestari@example.com", "Jl. Sudirman No. 2, Bandung"},
		{"Cahyo Widodo", 45, "Male", "+6281333333333", "cahyo.widodo@example.com", "Jl. Diponegoro No. 3, Surabaya"},
		{"Dewi Sartika", 32, "Female", "+6281444444444", "dewi.sartika@example.com", "Jl. Gajah Mada No. 4, Yogyakarta"},
		{"Eko Prasetyo", 50, "Male", "+6281555555555", "eko.prasetyo@example.com", "Jl. Pahlawan No. 5, Semarang"},
		{"Fitri Handayani", 25, "Female", "+6281666666666", "fitri.handayani@example.com", "Jl. Imam Bonjol No. 6, Medan"},
		{"Gilang Ramadhan", 38, "Male", "+6281777777777", "gilang.ramadhan@example.com", "Jl. Teuku Umar No. 7, Makassar"},
		{"Hesti Purwanti", 42, "Female", "+6281888888888", "hesti.purwanti@example.com", "Jl. Hasanuddin No. 8, Palembang"},
		{"Indra Gunawan", 29, "Male", "+6281999999999", "indra.gunawan@example.com", "Jl. Pattimura No. 9, Denpasar"},
		{"Jelita Sari", 33, "Female", "+6282111111111", "jelita.sari@example.com", "Jl. Ahmad Yani No. 10, Pontianak"},
	}

	for _, patientData := range patientsToCreate {
		patient := &domain.PatientEntity{
			Name:      patientData.Name,
			Age:       patientData.Age,
			Gender:    patientData.Gender,
			Phone:     patientData.Phone,
			Email:     patientData.Email,
			Address:   patientData.Address,
			CreatedBy: creatorObjID,
		}

		id, err := s.patientService.Create(ctx, patient)
		if err != nil {
			log.Printf("failed to create patient %s: %v, skipping", patient.Name, err)
			continue
		}
		newPatient, err := s.patientService.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve created patient %s: %w", patient.Name, err)
		}
		createdPatients = append(createdPatients, newPatient)
	}
	return createdPatients, nil
}

func (s *Seed) seedAppointments(ctx context.Context, creatorID string, doctors []*domain.DoctorEntity, patients []*domain.PatientEntity) ([]*domain.AppointmentEntity, error) {
	creatorObjID, _ := primitive.ObjectIDFromHex(creatorID)
	var createdAppointments []*domain.AppointmentEntity
	appointmentTypes := []domain.AppointmentType{
		domain.AppointmentTypeCheckUp,
		domain.AppointmentTypeConsultation,
		domain.AppointmentTypeFollowUp,
	}
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 20; i++ {
		patient := patients[rand.Intn(len(patients))]
		doctor := doctors[rand.Intn(len(doctors))]

		appointment := &domain.CreateAppointmentRequest{
			PatientID: patient.ID.Hex(),
			DoctorID:  doctor.ID.Hex(),
			Type:      appointmentTypes[rand.Intn(len(appointmentTypes))],
			DateTime:  time.Now().AddDate(0, 0, rand.Intn(60)-30),
			Duration:  30,
			Location:  fmt.Sprintf("Ruang %d", 101+rand.Intn(10)),
			Notes:     "Pemeriksaan rutin.",
		}

		id, err := s.appointmentService.Create(ctx, appointment, creatorObjID)
		if err != nil {
			log.Printf("Could not create appointment for patient %s: %v. Skipping.", patient.Name, err)
			continue
		}

		log.Printf("Created appointment with ID: %s", id)

		newAppt, err := s.appointmentService.GetByID(ctx, id)
		if err != nil {
			log.Printf("Failed to retrieve created appointment: %v. Skipping.", err)
			continue
		}

		createdAppointments = append(createdAppointments, newAppt)
	}

	return createdAppointments, nil
}

func (s *Seed) seedMedicalRecords(ctx context.Context, creatorID string, doctors []*domain.DoctorEntity, appointments []*domain.AppointmentEntity) ([]*domain.MedicalRecordEntity, error) {
	creatorObjID, _ := primitive.ObjectIDFromHex(creatorID)
	var createdRecords []*domain.MedicalRecordEntity
	rand.Seed(time.Now().UnixNano())

	for _, appt := range appointments {
		if appt.Status != domain.AppointmentStatusCompleted {
			continue
		}

		record := &domain.MedicalRecordEntity{
			PatientID:   appt.PatientID,
			DoctorID:    appt.DoctorID,
			Date:        appt.DateTime,
			RecordType:  domain.RecordTypeCheckUp,
			Description: "Pasien datang dengan keluhan...",
			Diagnosis:   "Hipertensi",
			Treatment:   "Pemberian obat A",
			Notes:       "Kontrol kembali 1 minggu lagi.",
			CreatedBy:   creatorObjID,
		}

		id, err := s.medicalRecordService.Create(ctx, record)
		if err != nil {
			log.Printf("Could not create medical record for patient %s: %v. Skipping.", appt.PatientID.Hex(), err)
			continue
		}

		newRecord, err := s.medicalRecordService.GetByID(ctx, id)
		if err != nil {
			log.Printf("Failed to retrieve medical record: %v. Skipping.", err)
			continue
		}

		createdRecords = append(createdRecords, newRecord)
	}

	// Tambahkan 5 medical record manual
	if len(appointments) >= 10 {
		manualRecords := []domain.MedicalRecordEntity{
			{
				PatientID:   appointments[5].PatientID,
				DoctorID:    appointments[5].DoctorID,
				Date:        time.Now().Add(-96 * time.Hour),
				RecordType:  domain.RecordTypeCheckUp,
				Description: "Pasien mengeluh nyeri dada ringan",
				Diagnosis:   "Asam lambung",
				Treatment:   "Obat antasida dan diet rendah asam",
				Notes:       "Hindari makanan pedas dan asam",
				CreatedBy:   creatorObjID,
			},
			{
				PatientID:   appointments[6].PatientID,
				DoctorID:    appointments[6].DoctorID,
				Date:        time.Now().Add(-72 * time.Hour),
				RecordType:  domain.RecordTypeFollowUp,
				Description: "Pemeriksaan lanjutan pasca operasi ringan",
				Diagnosis:   "Pemulihan baik",
				Treatment:   "Saran latihan ringan dan kontrol ulang",
				Notes:       "Kembali untuk cek luka 5 hari lagi",
				CreatedBy:   creatorObjID,
			},
			{
				PatientID:   appointments[7].PatientID,
				DoctorID:    appointments[7].DoctorID,
				Date:        time.Now().Add(-48 * time.Hour),
				RecordType:  domain.RecordTypeCheckUp,
				Description: "Pemeriksaan luka ringan karena jatuh",
				Diagnosis:   "Luka lecet",
				Treatment:   "Pembersihan luka dan salep antibiotik",
				Notes:       "Ganti perban tiap hari",
				CreatedBy:   creatorObjID,
			},
			{
				PatientID:   appointments[8].PatientID,
				DoctorID:    appointments[8].DoctorID,
				Date:        time.Now().Add(-24 * time.Hour),
				RecordType:  domain.RecordTypeCheckUp,
				Description: "Pasien mengalami batuk 3 hari",
				Diagnosis:   "Infeksi Saluran Pernapasan Atas",
				Treatment:   "Obat batuk dan istirahat",
				Notes:       "Kontrol ulang jika tidak membaik dalam 3 hari",
				CreatedBy:   creatorObjID,
			},
			{
				PatientID:   appointments[9].PatientID,
				DoctorID:    appointments[9].DoctorID,
				Date:        time.Now(),
				RecordType:  domain.RecordTypeFollowUp,
				Description: "Follow-up terapi hipertensi",
				Diagnosis:   "Hipertensi terkendali",
				Treatment:   "Lanjutkan obat dan jaga pola makan",
				Notes:       "Cek tekanan darah mingguan",
				CreatedBy:   creatorObjID,
			},
		}

		for _, record := range manualRecords {
			id, err := s.medicalRecordService.Create(ctx, &record)
			if err != nil {
				log.Printf("Manual record failed: %v", err)
				continue
			}
			newRecord, err := s.medicalRecordService.GetByID(ctx, id)
			if err != nil {
				log.Printf("Failed to retrieve manual record: %v", err)
				continue
			}
			createdRecords = append(createdRecords, newRecord)
		}
	}

	return createdRecords, nil
}
