package handlers

import (
	"github.com/gofiber/fiber/v2"
	"sentineldesk/backend/internal/services"
)

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

func (h *AdminHandler) UpdateProfile(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID := c.Locals("userID").(string)

	authService := services.NewAuthService()
	if err := authService.UpdateProfile(userID, req.Name, req.Email, req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}
