package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/services"
)

type SoftwareInventoryHandler struct {
	service *services.SoftwareInventoryService
}

func NewSoftwareInventoryHandler() *SoftwareInventoryHandler {
	return &SoftwareInventoryHandler{
		service: services.NewSoftwareInventoryService(),
	}
}

type SoftwareItem struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Publisher       string `json:"publisher"`
	InstallDate     string `json:"install_date"`
	InstallLocation string `json:"install_location"`
	EstimatedSize   int64  `json:"estimated_size"`
}

type SaveSoftwareRequest struct {
	DeviceID string        `json:"device_id"`
	Software []SoftwareItem `json:"software"`
}

func (h *SoftwareInventoryHandler) Save(c *fiber.Ctx) error {

	var req SaveSoftwareRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	deviceID, err := validateDeviceExists(c, req.DeviceID)
	if err != nil {
		return err
	}

	software := make([]models.SoftwareInventory, len(req.Software))

	for i, item := range req.Software {
		software[i] = models.SoftwareInventory{
			DeviceID:        deviceID,
			Name:            item.Name,
			Version:         item.Version,
			Publisher:       item.Publisher,
			InstallDate:     item.InstallDate,
			InstallLocation: item.InstallLocation,
			EstimatedSize:   item.EstimatedSize,
		}
	}

	if err := h.service.Replace(deviceID, software); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save software inventory",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Software inventory saved successfully",
		"count":    len(software),
	})
}

func (h *SoftwareInventoryHandler) GetByDeviceID(c *fiber.Ctx) error {

	deviceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid device UUID",
		})
	}

	software, err := h.service.GetByDeviceID(deviceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Software inventory not found",
		})
	}

	return c.JSON(fiber.Map{
		"count":    len(software),
		"software": software,
	})
}
