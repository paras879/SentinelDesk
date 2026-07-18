package routes

import (
	"sentineldesk/backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// ==========================
	// Public Routes
	// ==========================
	app.Get("/", handlers.Home)
	app.Get("/health", handlers.Health)
	app.Get("/version", handlers.Version)

	// ==========================
	// Authentication
	// ==========================
	SetupAuthRoutes(app)

	// ==========================
	// Admin
	// ==========================
	SetupAdminRoutes(app)

	// ==========================
	// Device System Information
	// (registered before Device routes so the static path /system-info
	// takes priority over the parameterized route /:id)
	// ==========================
	SetupDeviceSystemInfoRoutes(app)

	// ==========================
	// Software Inventory
	// ==========================
	SetupSoftwareInventoryRoutes(app)

	// ==========================
	// Running Processes
	// ==========================
	SetupProcessRoutes(app)

	// ==========================
	// Windows Services
	// ==========================
	SetupWindowsServiceRoutes(app)

	// ==========================
	// Agent Command Polling
	// ==========================
	SetupAgentRoutes(app)

	// ==========================
	// Device
	// ==========================
	SetupDeviceRoutes(app)

	// ==========================
	// Dashboard
	// ==========================
	SetupDashboardRoutes(app)

	// ==========================
	// Live View
	// ==========================
	SetupLiveViewRoutes(app)
}
