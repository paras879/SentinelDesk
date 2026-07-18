package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"sentineldesk/backend/internal/services"
)

type ServiceCommandHandler struct {
	cmdService *services.ServiceCommandService
}

func NewServiceCommandHandler() *ServiceCommandHandler {
	return &ServiceCommandHandler{
		cmdService: services.NewServiceCommandService(),
	}
}

type ReportCommandResultRequest struct {
	CommandID    string `json:"command_id"`
	Result       string `json:"result"`
	ErrorMessage string `json:"error_message"`
}

func (h *ServiceCommandHandler) GetPendingCommands(c *fiber.Ctx) error {

	deviceIDParam := c.Query("device_id")
	if deviceIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "device_id query parameter is required",
		})
	}

	deviceID, err := uuid.Parse(deviceIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Device ID format",
		})
	}

	commands, err := h.cmdService.GetPendingByDeviceID(deviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pending commands",
		})
	}

	return c.JSON(fiber.Map{
		"count":    len(commands),
		"commands": commands,
	})
}

func (h *ServiceCommandHandler) ReportResult(c *fiber.Ctx) error {

	var req ReportCommandResultRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	cmdID, err := uuid.Parse(req.CommandID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid command ID",
		})
	}

	if req.Result == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Result is required",
		})
	}

	if err := h.cmdService.MarkCompleted(cmdID, req.Result, req.ErrorMessage); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update command status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Command result reported successfully",
	})
}
