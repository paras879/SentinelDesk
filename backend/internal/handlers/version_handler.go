package handlers

import "github.com/gofiber/fiber/v2"

func Version(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"application": "SentinelDesk",
		"version":     "1.0.0",
		"api":         "v1",
	})
}
