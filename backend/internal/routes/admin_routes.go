package routes

import (
	"github.com/gofiber/fiber/v2"

	"sentineldesk/backend/internal/handlers"
	"sentineldesk/backend/internal/middleware"
)

func SetupAdminRoutes(app *fiber.App) {

	admin := handlers.NewAdminHandler()

	api := app.Group("/api/v1/admin")

	// Protected Routes
	api.Use(middleware.JWTProtected())

	api.Get("/profile", admin.Profile)
}
