package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/services"
)

type DeviceSystemInfoHandler struct {
	service *services.DeviceSystemInfoService
}

func NewDeviceSystemInfoHandler() *DeviceSystemInfoHandler {
	return &DeviceSystemInfoHandler{
		service: services.NewDeviceSystemInfoService(),
	}
}

type SaveSystemInfoRequest struct {
	DeviceID string `json:"device_id"`

	CPUName  string  `json:"cpu_name"`
	CPUUsage float64 `json:"cpu_usage"`

	TotalRAM uint64 `json:"total_ram"`
	UsedRAM  uint64 `json:"used_ram"`
	FreeRAM  uint64 `json:"free_ram"`

	TotalDisk uint64 `json:"total_disk"`
	UsedDisk  uint64 `json:"used_disk"`
	FreeDisk  uint64 `json:"free_disk"`

	OSVersion string `json:"os_version"`

	LocalIP  string `json:"local_ip"`
	PublicIP string `json:"public_ip"`

	MACAddress string `json:"mac_address"`

	Uptime int64 `json:"uptime"`

	LastBoot time.Time `json:"last_boot"`

	AgentVersion string `json:"agent_version"`
}

// POST /api/v1/devices/system-info
func (h *DeviceSystemInfoHandler) Save(c *fiber.Ctx) error {

	var req SaveSystemInfoRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	deviceID, err := validateDeviceExists(c, req.DeviceID)
	if err != nil {
		return err
	}

	info := &models.DeviceSystemInfo{
		DeviceID: deviceID,

		CPUName:  req.CPUName,
		CPUUsage: req.CPUUsage,

		TotalRAM: req.TotalRAM,
		UsedRAM:  req.UsedRAM,
		FreeRAM:  req.FreeRAM,

		TotalDisk: req.TotalDisk,
		UsedDisk:  req.UsedDisk,
		FreeDisk:  req.FreeDisk,

		OSVersion: req.OSVersion,

		LocalIP:  req.LocalIP,
		PublicIP: req.PublicIP,

		MACAddress: req.MACAddress,

		Uptime: req.Uptime,

		LastBoot: req.LastBoot,

		AgentVersion: req.AgentVersion,
	}

	if err := h.service.Save(info); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save system information",
		})
	}

	return c.JSON(fiber.Map{
		"message": "System information saved successfully",
	})
}

// GET /api/v1/devices/:id/system-info
func (h *DeviceSystemInfoHandler) GetByDeviceID(c *fiber.Ctx) error {

	deviceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid device UUID",
		})
	}

	info, err := h.service.GetByDeviceID(deviceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "System information not found",
		})
	}

	return c.JSON(info)
}
