package handlers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/services"
)

// validateDeviceExists checks that a device with the given device_id
// exists in the devices table before processing child records.
// Returns the parsed UUID or a 404 response.
func validateDeviceExists(c *fiber.Ctx, deviceIDStr string) (uuid.UUID, error) {
	if deviceIDStr == "" {
		return uuid.Nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Device ID is required",
		})
	}

	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return uuid.Nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Device ID format",
		})
	}

	deviceService := services.NewDeviceService()
	device, err := deviceService.GetByID(deviceID)
	if err != nil || device == nil {
		return uuid.Nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Device not registered. Please register the device first.",
		})
	}

	return deviceID, nil
}

type DeviceHandler struct {
	service *services.DeviceService
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{
		service: services.NewDeviceService(),
	}
}

type RegisterDeviceRequest struct {
	DeviceID        string          `json:"device_id"`
	DeviceName      string          `json:"device_name"`
	Hostname        string          `json:"hostname"`
	Username        string          `json:"username"`
	OS              string          `json:"os"`
	OSVersion       string          `json:"os_version"`
	IPAddress       string          `json:"ip_address"`
	MACAddress      string          `json:"mac_address"`
	ConnectedSubnet string          `json:"connected_subnet"`
	NetworkAdapters json.RawMessage `json:"network_adapters"`
	DefaultGateway  string          `json:"default_gateway"`
	NetworkGroupID  string          `json:"network_group_id"`
}

type HeartbeatRequest struct {
	DeviceID        string `json:"device_id"`
	Hostname        string `json:"hostname"`
	IPAddress       string `json:"ip_address"`
	MACAddress      string `json:"mac_address"`
	ConnectedSubnet string `json:"connected_subnet"`
	DefaultGateway  string `json:"default_gateway"`
	NetworkGroupID  string `json:"network_group_id"`
}

// ==========================
// Register Device
// ==========================
func (h *DeviceHandler) Register(c *fiber.Ctx) error {

	var req RegisterDeviceRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Register: BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.DeviceID == "" || req.DeviceName == "" {
		log.Printf("Register: Missing required fields: device_id=%q device_name=%q",
			req.DeviceID, req.DeviceName)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Device ID and Device Name are required",
		})
	}

	deviceID, err := uuid.Parse(req.DeviceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Device ID format",
		})
	}

	isUpdate, err := h.service.Register(
		deviceID,
		req.DeviceName,
		req.Hostname,
		req.Username,
		req.OS,
		req.OSVersion,
		req.IPAddress,
		req.MACAddress,
		req.ConnectedSubnet,
		string(req.NetworkAdapters),
		req.DefaultGateway,
		req.NetworkGroupID,
	)

	if err != nil {
		log.Printf("Register: Service error for device_id=%s: %v", req.DeviceID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if isUpdate {
		log.Printf("Register: Device already registered, record updated: device_id=%s", req.DeviceID)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Device already registered, record updated",
		})
	}

	log.Printf("Register: New device registered: device_id=%s", req.DeviceID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Device registered successfully",
	})
}

// ==========================
// Get All Devices
// Supports ?network_group_id= for site filtering
// Supports ?network_id= for legacy subnet-based filtering
// ==========================
func (h *DeviceHandler) GetAll(c *fiber.Ctx) error {

	networkGroupID := c.Query("network_group_id")
	networkID := c.Query("network_id")

	var devices []models.Device
	var err error

	switch {
	case networkGroupID != "":
		devices, err = h.service.GetByNetworkGroupID(networkGroupID)
	case networkID != "":
		devices, err = h.service.GetByNetworkID(networkID)
	default:
		devices, err = h.service.GetAll()
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch devices",
		})
	}

	return c.JSON(fiber.Map{
		"count":   len(devices),
		"devices": devices,
	})
}

// ==========================
// Get Devices By Network ID
// ==========================
func (h *DeviceHandler) GetByNetwork(c *fiber.Ctx) error {
	networkID := c.Params("networkID")
	if networkID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Network ID is required",
		})
	}

	devices, err := h.service.GetByNetworkID(networkID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch devices",
		})
	}

	return c.JSON(fiber.Map{
		"count":   len(devices),
		"devices": devices,
	})
}

// ==========================
// Heartbeat
// ==========================
func (h *DeviceHandler) Heartbeat(c *fiber.Ctx) error {

	var req HeartbeatRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.DeviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Device ID is required",
		})
	}

	deviceID, err := uuid.Parse(req.DeviceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Device ID format",
		})
	}

	err = h.service.Heartbeat(
		deviceID,
		req.IPAddress,
		req.MACAddress,
		req.ConnectedSubnet,
		req.DefaultGateway,
		req.NetworkGroupID,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Heartbeat received successfully",
	})
}

// ==========================
// Get Device By UUID
// ==========================
func (h *DeviceHandler) GetByID(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid device UUID",
		})
	}

	device, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Device not found",
		})
	}

	return c.JSON(device)
}
