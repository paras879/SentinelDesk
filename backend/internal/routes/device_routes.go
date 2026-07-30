package routes

import (
	"github.com/gofiber/contrib/websocket"
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

	app.Get(
		"/ws/heartbeat",
		middleware.AgentProtected(),
		websocket.New(device.HeartbeatWS),
	)

	// ==========================
	// Admin Routes
	// ==========================
	api.Get("/", middleware.JWTProtected(), device.GetAll)
	api.Get("/network/:networkID", middleware.JWTProtected(), device.GetByNetwork)
	api.Get("/:id", middleware.JWTProtected(), device.GetByID)
	api.Put("/:id/location", middleware.JWTProtected(), device.UpdateLocation)
	api.Delete("/:id", middleware.JWTProtected(), device.Delete)
}
