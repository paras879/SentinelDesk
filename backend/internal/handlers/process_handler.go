package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/services"
)

type ProcessHandler struct {
	service *services.ProcessService
}

func NewProcessHandler() *ProcessHandler {
	return &ProcessHandler{
		service: services.NewProcessService(),
	}
}

type ProcessItem struct {
	PID             int32   `json:"pid"`
	Name            string  `json:"name"`
	ExecutablePath  string  `json:"executable_path"`
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryBytes     uint64  `json:"memory_bytes"`
	MemoryPercent   float32 `json:"memory_percent"`
	ThreadCount     int32   `json:"thread_count"`
	HandleCount     int32   `json:"handle_count"`
	StartTime       string  `json:"start_time"`
	Username        string  `json:"username"`
	Status          string  `json:"status"`
}

type SaveProcessRequest struct {
	DeviceID  string        `json:"device_id"`
	Processes []ProcessItem `json:"processes"`
}

func (h *ProcessHandler) Save(c *fiber.Ctx) error {

	var req SaveProcessRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	deviceID, err := validateDeviceExists(c, req.DeviceID)
	if err != nil {
		return err
	}

	processes := make([]models.Process, len(req.Processes))

	for i, item := range req.Processes {

		startTime := time.Time{}
		if item.StartTime != "" {
			parsed, err := time.Parse(time.RFC3339, item.StartTime)
			if err == nil {
				startTime = parsed
			}
		}

		processes[i] = models.Process{
			DeviceID:       deviceID,
			PID:            item.PID,
			Name:           item.Name,
			ExecutablePath: item.ExecutablePath,
			CPUUsage:       item.CPUUsage,
			MemoryBytes:    item.MemoryBytes,
			MemoryPercent:  item.MemoryPercent,
			ThreadCount:    item.ThreadCount,
			HandleCount:    item.HandleCount,
			StartTime:      startTime,
			Username:       item.Username,
			Status:         "Running",
		}
	}

	if err := h.service.Replace(deviceID, processes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save process list",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Process list saved successfully",
		"count":    len(processes),
	})
}

func (h *ProcessHandler) GetByDeviceID(c *fiber.Ctx) error {

	deviceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid device UUID",
		})
	}

	processes, err := h.service.GetByDeviceID(deviceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Process list not found",
		})
	}

	return c.JSON(fiber.Map{
		"count":     len(processes),
		"processes": processes,
	})
}
