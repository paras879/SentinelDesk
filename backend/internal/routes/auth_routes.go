package routes

import (
	"github.com/gofiber/fiber/v2"

	"sentineldesk/backend/internal/handlers"
)

func SetupAuthRoutes(app *fiber.App) {

	auth := handlers.NewAuthHandler()

	api := app.Group("/api/v1/auth")

	api.Post("/register", auth.Register)
	api.Post("/login", auth.Login) // 👈 Ye line add karo
}
