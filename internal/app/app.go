package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/ekastn/hms-api/internal/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type App struct {
	f   *fiber.App
	cfg config
	db  *mongo.Database
}

type config struct {
	addr     string
	mongoCfg mongoDbCfg
}

type mongoDbCfg struct {
	addr string
	db   string
}

func (a *App) Run() {
	a.loadConfig()
	a.connectDb()
	a.mount()
	a.setupRoutes()

	go func() {
		if err := a.f.Listen(a.cfg.addr); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Println("shutting down server...")
	if err := a.f.Shutdown(); err != nil {
		log.Fatal("error while shutting down server: ", err)
	}

	log.Println("server shut down")
}

func (a *App) mount() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174"),
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE",
	}))

	a.f = app
}

func (a *App) connectDb() {
	log.Println("connecting to database", a.cfg.mongoCfg.addr)
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
	log.Println("database ping")

	a.db = client.Database(a.cfg.mongoCfg.db)
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
	}

	a.cfg = cfg
}
