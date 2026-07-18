package middleware

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func AgentProtected() fiber.Handler {

	expectedKey := os.Getenv("AGENT_API_KEY")
	if expectedKey == "" {
		log.Fatal("AGENT_API_KEY environment variable is not set")
	}

	return func(c *fiber.Ctx) error {

		apiKey := c.Get("X-Agent-Key")

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Agent API Key",
			})
		}

		if apiKey != expectedKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Agent API Key",
			})
		}

		return c.Next()
	}
}
