package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"sentineldesk/backend/internal/services"
	"sentineldesk/backend/internal/utils"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service: services.NewAuthService(),
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {

	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Admin registered successfully",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("[DEBUG] Login body parse error: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("[DEBUG] Login request — Email: %q | Password: %q", req.Email, req.Password)

	if req.Email == "" || req.Password == "" {
		log.Printf("[DEBUG] Login validation failed — email empty: %v | password empty: %v", req.Email == "", req.Password == "")
		return c.Status(400).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	admin, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		log.Printf("[DEBUG] Login service returned error for %q — reason: %v", req.Email, err)
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := utils.GenerateJWT(admin)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user": fiber.Map{
			"id":    admin.ID,
			"name":  admin.Name,
			"email": admin.Email,
			"role":  admin.Role,
		},
	})
}
