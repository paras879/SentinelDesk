package handlers

import "github.com/gofiber/fiber/v2"

func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
		"server": "running",
		"uptime": "active",
		"http":   200,
	})
}
