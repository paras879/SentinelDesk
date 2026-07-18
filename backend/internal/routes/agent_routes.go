package routes

import (
	"github.com/gofiber/fiber/v2"

	"sentineldesk/backend/internal/handlers"
	"sentineldesk/backend/internal/middleware"
)

func SetupAgentRoutes(app *fiber.App) {

	handler := handlers.NewServiceCommandHandler()

	api := app.Group("/api/v1/agent")

	api.Use(middleware.AgentProtected())

	api.Get("/commands", handler.GetPendingCommands)
	api.Post("/commands/result", handler.ReportResult)
}
