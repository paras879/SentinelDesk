package routes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"sentineldesk/backend/internal/liveview"
	"sentineldesk/backend/internal/middleware"
	"sentineldesk/backend/internal/utils"
)

func SetupLiveViewRoutes(app *fiber.App) {

	app.Get("/ws/live/stream", middleware.AgentProtected(), func(c *fiber.Ctx) error {
		deviceID := c.Query("device_id")
		if deviceID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "device_id required"})
		}
		c.Locals("deviceID", deviceID)
		return c.Next()
	}, websocket.New(liveview.HandleAgentStream))

	app.Get("/ws/live/view/:deviceID", func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing token"})
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		if claims["role"] != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}

		deviceID := c.Params("deviceID")
		_, err = uuid.Parse(deviceID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid device id"})
		}

		c.Locals("deviceID", deviceID)
		return c.Next()
	}, websocket.New(liveview.HandleAdminView))
}
