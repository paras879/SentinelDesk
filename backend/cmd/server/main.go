package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"sentineldesk/backend/internal/config"
	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/routes"
	"sentineldesk/backend/internal/scheduler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	cfg := config.LoadConfig()

	// Database Connection
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Database Connection Failed:", err)
	}

	// Database Migration
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Migration Failed:", err)
	}

	// Seed default admin
	database.SeedAdmin()

	// Start Background Scheduler
	go scheduler.StartDeviceStatusScheduler()

	app := fiber.New()

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())

	// CORS - Allow frontend from configured origins
	allowOrigins := cfg.AllowedOrigins
	allowCredentials := allowOrigins != "*" && !strings.Contains(allowOrigins, "*")

	if !allowCredentials {
		log.Println("CORS: AllowOrigins is wildcard, AllowCredentials set to false")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Agent-Key",
		AllowCredentials: allowCredentials,
	}))

	// Routes
	routes.Setup(app)

	// Graceful Shutdown
	go func() {
		log.Printf("%s v%s is running on http://localhost:%s",
			cfg.AppName,
			cfg.AppVersion,
			cfg.AppPort,
		)

		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	log.Println("Server exited cleanly")
}
