package routes

import (
	"sentineldesk/backend/internal/handlers"
	"sentineldesk/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupDashboardRoutes(app *fiber.App) {

	handler := handlers.NewDashboardHandler()

	api := app.Group("/api/v1/dashboard")

	api.Use(middleware.JWTProtected())

	api.Get("/summary", handler.Summary)
	api.Get("/my-ip", handler.MyIP)
}
