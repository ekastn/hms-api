
package app

import (
	"context"
	"log"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateInitialAdmin(db *mongo.Database) {
	initialAdminEmail := env.GetString("INITIAL_ADMIN_EMAIL", "")
	initialAdminPassword := env.GetString("INITIAL_ADMIN_PASSWORD", "")

	if initialAdminEmail == "" || initialAdminPassword == "" {
		log.Println("Initial admin setup skipped: environment variables not set")
		return
	}

	users := db.Collection("users")
	filter := bson.M{"email": initialAdminEmail}
	count, err := users.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Printf("Error checking for initial admin: %v", err)
		return
	}

	if count > 0 {
		log.Println("Initial admin user already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(initialAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing initial admin password: %v", err)
		return
	}

	adminUser := &domain.UserEntity{
		Name:      "Super Admin",
		Email:     initialAdminEmail,
		Password:  string(hashedPassword),
		Role:      domain.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = users.InsertOne(context.Background(), adminUser)
	if err != nil {
		log.Printf("Error creating initial admin user: %v", err)
		return
	}

	log.Println("Initial admin user created successfully")
}
