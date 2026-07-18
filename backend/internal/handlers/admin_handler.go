package handlers

import "github.com/gofiber/fiber/v2"

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) Profile(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"message": "Profile fetched successfully",
		"user": fiber.Map{
			"id":    c.Locals("userID"),
			"name":  c.Locals("name"),
			"email": c.Locals("email"),
			"role":  c.Locals("role"),
		},
	})
}
