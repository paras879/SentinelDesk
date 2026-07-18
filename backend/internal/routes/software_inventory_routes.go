package routes

import (
	"github.com/gofiber/fiber/v2"

	"sentineldesk/backend/internal/handlers"
	"sentineldesk/backend/internal/middleware"
)

func SetupSoftwareInventoryRoutes(app *fiber.App) {

	handler := handlers.NewSoftwareInventoryHandler()

	api := app.Group("/api/v1/devices")

	// ==========================
	// Agent APIs
	// ==========================
	api.Post(
		"/software",
		middleware.AgentProtected(),
		handler.Save,
	)

	// ==========================
	// Admin APIs
	// ==========================
	api.Get(
		"/:id/software",
		middleware.JWTProtected(),
		handler.GetByDeviceID,
	)
}
