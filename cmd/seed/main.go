package main

import (
	"context"
	"log"
	"time"

	"github.com/ekastn/hms-api/internal/env"
	"github.com/ekastn/hms-api/internal/repository"
	"github.com/ekastn/hms-api/internal/seed"
	"github.com/ekastn/hms-api/internal/service"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := env.GetString("MONGO_ADDR", "mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Connecting to MongoDB", mongoURI)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	dbName := env.GetString("MONGO_DB", "test_db")
	db := client.Database(dbName)

	userRepo := repository.NewUserRepository(db.Collection("users"))
	doctorRepo := repository.NewDoctorRepository(db.Collection("doctors"))
	patientRepo := repository.NewPatientRepository(db.Collection("patients"))
	appointmentRepo := repository.NewAppointmentRepository(db.Collection("appointments"))
	medicalRecordRepo := repository.NewMedicalRecordRepository(db.Collection("medical_records"))
	activityRepo := repository.NewActivityRepository(db.Collection("activities"))

	activityService := service.NewActivityService(activityRepo)
	userService := service.NewUserService(userRepo)
	doctorService := service.NewDoctorService(doctorRepo, appointmentRepo, patientRepo, activityService)
	patientService := service.NewPatientService(patientRepo, appointmentRepo, medicalRecordRepo, activityService)
	appointmentService := service.NewAppointmentService(appointmentRepo, patientRepo, medicalRecordRepo, activityService, client)
	medicalRecordService := service.NewMedicalRecordService(medicalRecordRepo, activityService)

	seeder := seed.NewSeeder(db, userService, doctorService, patientService, appointmentService, medicalRecordService)

	seeder.Seed(context.Background())
}
