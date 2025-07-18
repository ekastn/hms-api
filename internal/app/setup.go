package app

import (
	"context"
	"log"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/env"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (a *App) connectDb() {
	log.Println("connecting to", a.cfg.mongoCfg.addr)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(a.cfg.mongoCfg.addr))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	a.db = client.Database(a.cfg.mongoCfg.db)

	a.createInitialAdmin()
}

func (a *App) loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := config{
		addr: env.GetString("APP_ADDR", ":3000"),
		mongoCfg: mongoDbCfg{
			addr: env.GetString("MONGO_ADDR", "mongodb://localhost:27017"),
			db:   env.GetString("MONGO_DB", "hms"),
		},
		jwtSecret: env.GetString("JWT_SECRET", "your-secret-key"),
	}

	a.cfg = cfg
}

func (a *App) createInitialAdmin() {
	initialAdminEmail := env.GetString("INITIAL_ADMIN_EMAIL", "")
	initialAdminPassword := env.GetString("INITIAL_ADMIN_PASSWORD", "")

	if initialAdminEmail == "" || initialAdminPassword == "" {
		log.Println("Initial admin setup skipped: environment variables not set")
		return
	}

	log.Printf("Initial admin user: email=%s, password=%s\n", initialAdminEmail, initialAdminPassword)

	users := a.db.Collection("users")
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

	hashedPassword, err := utils.HashPassword(initialAdminPassword)
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
