package handlers

import "github.com/gofiber/fiber/v2"

func Home(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"project": "SentinelDesk",
		"status":  "Running",
		"version": "1.0.0",
	})
}
