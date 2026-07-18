package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/services"
)

type WindowsServiceHandler struct {
	service    *services.WindowsServiceService
	cmdService *services.ServiceCommandService
}

func NewWindowsServiceHandler() *WindowsServiceHandler {
	return &WindowsServiceHandler{
		service:    services.NewWindowsServiceService(),
		cmdService: services.NewServiceCommandService(),
	}
}

type WindowsServiceItem struct {
	ServiceName    string `json:"service_name"`
	DisplayName    string `json:"display_name"`
	Status         string `json:"status"`
	StartType      string `json:"start_type"`
	ExecutablePath string `json:"executable_path"`
	PID            int32  `json:"pid"`
	ServiceAccount string `json:"service_account"`
	Description    string `json:"description"`
	CanStop        bool   `json:"can_stop"`
	CanPause       bool   `json:"can_pause"`
	AcceptShutdown bool   `json:"accept_shutdown"`
}

type SaveServicesRequest struct {
	DeviceID string              `json:"device_id"`
	Services []WindowsServiceItem `json:"services"`
}

type ServiceActionRequest struct {
	ServiceName string `json:"service_name"`
}

func (h *WindowsServiceHandler) Save(c *fiber.Ctx) error {

	var req SaveServicesRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	deviceID, err := validateDeviceExists(c, req.DeviceID)
	if err != nil {
		return err
	}

	services := make([]models.WindowsService, len(req.Services))

	for i, item := range req.Services {
		services[i] = models.WindowsService{
			DeviceID:       deviceID,
			ServiceName:    item.ServiceName,
			DisplayName:    item.DisplayName,
			Status:         item.Status,
			StartType:      item.StartType,
			ExecutablePath: item.ExecutablePath,
			PID:            item.PID,
			ServiceAccount: item.ServiceAccount,
			Description:    item.Description,
			CanStop:        item.CanStop,
			CanPause:       item.CanPause,
			AcceptShutdown: item.AcceptShutdown,
		}
	}

	if err := h.service.Replace(deviceID, services); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save Windows services",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Windows services saved successfully",
		"count":   len(services),
	})
}

func (h *WindowsServiceHandler) GetByDeviceID(c *fiber.Ctx) error {

	deviceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid device UUID",
		})
	}

	services, err := h.service.GetByDeviceID(deviceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Services not found",
		})
	}

	return c.JSON(fiber.Map{
		"count":    len(services),
		"services": services,
	})
}

func (h *WindowsServiceHandler) StartService(c *fiber.Ctx) error {
	return h.createServiceCommand(c, "start")
}

func (h *WindowsServiceHandler) StopService(c *fiber.Ctx) error {
	return h.createServiceCommand(c, "stop")
}

func (h *WindowsServiceHandler) RestartService(c *fiber.Ctx) error {
	return h.createServiceCommand(c, "restart")
}

func (h *WindowsServiceHandler) createServiceCommand(c *fiber.Ctx, action string) error {

	deviceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid device UUID",
		})
	}

	var req ServiceActionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ServiceName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Service name is required",
		})
	}

	cmd, err := h.cmdService.Create(deviceID, req.ServiceName, action)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create service command",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Service command queued",
		"command": fiber.Map{
			"id":           cmd.ID,
			"service_name": cmd.ServiceName,
			"action":       cmd.Action,
			"status":       cmd.Status,
		},
	})
}
