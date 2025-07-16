package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ekastn/hms-api/internal/env"
	"github.com/ekastn/hms-api/internal/seed"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// NOTE: This seeder is for development purposes only.
	// The initial admin user is now created on application startup.

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

	seeder := seed.NewSeeder(db)
	fmt.Println("Starting database seeding...")

	if err := seeder.SeedAll(context.Background()); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	fmt.Println("Database seeded successfully!")
}