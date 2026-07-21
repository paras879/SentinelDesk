package routes

import (
	"github.com/gofiber/fiber/v2"

	"sentineldesk/backend/internal/handlers"
	"sentineldesk/backend/internal/middleware"
)

func SetupDeviceRoutes(app *fiber.App) {

	device := handlers.NewDeviceHandler()

	api := app.Group("/api/v1/devices")

	// ==========================
	// Agent Routes
	// ==========================
	api.Post(
		"/register",
		middleware.AgentProtected(),
		device.Register,
	)

	api.Post(
		"/heartbeat",
		middleware.AgentProtected(),
		device.Heartbeat,
	)

	// ==========================
	// Admin Routes
	// ==========================
	api.Get("/", middleware.JWTProtected(), device.GetAll)
	api.Get("/network/:networkID", middleware.JWTProtected(), device.GetByNetwork)
	api.Get("/:id", middleware.JWTProtected(), device.GetByID)
	api.Delete("/:id", middleware.JWTProtected(), device.Delete)
}
