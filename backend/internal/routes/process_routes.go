package routes

import (
	"github.com/gofiber/fiber/v2"

	"sentineldesk/backend/internal/handlers"
	"sentineldesk/backend/internal/middleware"
)

func SetupProcessRoutes(app *fiber.App) {

	handler := handlers.NewProcessHandler()

	api := app.Group("/api/v1/devices")

	// ==========================
	// Agent APIs
	// ==========================
	api.Post(
		"/processes",
		middleware.AgentProtected(),
		handler.Save,
	)

	// ==========================
	// Admin APIs
	// ==========================
	api.Get(
		"/:id/processes",
		middleware.JWTProtected(),
		handler.GetByDeviceID,
	)
}
