package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/ekastn/hms-api/internal/env"
)

type App struct {
	f *fiber.App
	cfg config
}

type config struct {
	addr string
	mongoCfg mongoDbCfg
}

type mongoDbCfg struct {
	uri string
	db  string
}

func (a *App) Run() {
	a.loadConfig()
	a.mount()
	a.SetupRoutes()

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

func (a *App) loadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config{
		addr: env.GetString("APP_ADDR", ":3000"),
		mongoCfg: mongoDbCfg{
			uri: env.GetString("MONGO_URI", "mongodb://localhost:27017"),
			db:  env.GetString("MONGO_DB", "hms"),
		},
	}

	a.cfg = cfg
}
