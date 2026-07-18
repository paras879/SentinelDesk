package routes

import (
	"github.com/gofiber/fiber/v2"

	"sentineldesk/backend/internal/handlers"
	"sentineldesk/backend/internal/middleware"
)

func SetupDeviceSystemInfoRoutes(app *fiber.App) {

	handler := handlers.NewDeviceSystemInfoHandler()

	api := app.Group("/api/v1/devices")

	// ==========================
	// Agent APIs
	// ==========================
	api.Post(
		"/system-info",
		middleware.AgentProtected(),
		handler.Save,
	)

	// ==========================
	// Admin APIs
	// ==========================
	api.Get(
		"/:id/system-info",
		middleware.JWTProtected(),
		handler.GetByDeviceID,
	)
}
