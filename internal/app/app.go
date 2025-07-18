package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/ekastn/hms-api/internal/env"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	f   *fiber.App
	cfg config
	db  *mongo.Database
}

type config struct {
	addr      string
	mongoCfg  mongoDbCfg
	jwtSecret string
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
		AllowOrigins:     env.GetString("CORS_ALLOWED_ORIGINS", "http://localhost:5174"),
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE",
	}))

	a.f = app
}
