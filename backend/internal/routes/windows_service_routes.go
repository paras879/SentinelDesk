package routes

import (
	"github.com/gofiber/fiber/v2"

	"sentineldesk/backend/internal/handlers"
	"sentineldesk/backend/internal/middleware"
)

func SetupWindowsServiceRoutes(app *fiber.App) {

	handler := handlers.NewWindowsServiceHandler()

	api := app.Group("/api/v1/devices")

	// ==========================
	// Agent APIs
	// ==========================
	api.Post(
		"/services",
		middleware.AgentProtected(),
		handler.Save,
	)

	// ==========================
	// Admin APIs
	// ==========================
	api.Get(
		"/:id/services",
		middleware.JWTProtected(),
		handler.GetByDeviceID,
	)

	api.Post(
		"/:id/services/start",
		middleware.JWTProtected(),
		handler.StartService,
	)

	api.Post(
		"/:id/services/stop",
		middleware.JWTProtected(),
		handler.StopService,
	)

	api.Post(
		"/:id/services/restart",
		middleware.JWTProtected(),
		handler.RestartService,
	)
}
