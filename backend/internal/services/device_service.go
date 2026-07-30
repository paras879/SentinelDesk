package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/repository"
)

type DeviceService struct {
	repo *repository.DeviceRepository
}

func NewDeviceService() *DeviceService {
	return &DeviceService{
		repo: repository.NewDeviceRepository(),
	}
}

// Register Device using Device ID (UUID)
func (s *DeviceService) Register(
	deviceID uuid.UUID,
	deviceName string,
	hostname string,
	username string,
	os string,
	osVersion string,
	ipAddress string,
	macAddress string,
	connectedSubnet string,
	networkAdaptersJSON string,
	defaultGateway string,
	networkGroupID string,
) (bool, error) {

	existing, err := s.repo.GetByDeviceID(deviceID)
	if err == nil && existing != nil {
		existing.DeviceName = deviceName
		existing.Hostname = hostname
		existing.Username = username
		existing.OS = os
		existing.OSVersion = osVersion
		existing.IPAddress = ipAddress
		existing.MACAddress = macAddress
		existing.ConnectedSubnet = connectedSubnet
		existing.NetworkAdapters = networkAdaptersJSON
		existing.DefaultGateway = defaultGateway
		existing.NetworkGroupID = networkGroupID
		existing.Status = "online"
		now := time.Now()
		existing.LastSeen = &now
		return true, s.repo.Update(existing)
	}

	device := &models.Device{
		BaseModel:       models.BaseModel{ID: deviceID},
		DeviceName:      deviceName,
		Hostname:        hostname,
		Username:        username,
		OS:              os,
		OSVersion:       osVersion,
		IPAddress:       ipAddress,
		MACAddress:      macAddress,
		ConnectedSubnet: connectedSubnet,
		NetworkAdapters: networkAdaptersJSON,
		DefaultGateway:  defaultGateway,
		NetworkGroupID:  networkGroupID,
		Status:          "online",
	}

	return false, s.repo.Create(device)
}

// Get All Devices
func (s *DeviceService) GetAll() ([]models.Device, error) {
	return s.repo.GetAll()
}

// Get Devices By Network ID
func (s *DeviceService) GetByNetworkID(networkID string) ([]models.Device, error) {
	if networkID == "" {
		return s.repo.GetAll()
	}
	return s.repo.GetByNetworkID(networkID)
}

// Heartbeat by Device ID
func (s *DeviceService) Heartbeat(
	deviceID uuid.UUID,
	ipAddress string,
	macAddress string,
	connectedSubnet string,
	defaultGateway string,
	networkGroupID string,
) error {

	device, err := s.repo.GetByDeviceID(deviceID)
	if err != nil || device == nil {
		return errors.New("device not found")
	}

	return s.repo.UpdateHeartbeat(deviceID, ipAddress, macAddress, connectedSubnet, defaultGateway, networkGroupID)
}

// Get Devices By Network Group ID
func (s *DeviceService) GetByNetworkGroupID(groupID string) ([]models.Device, error) {
	if groupID == "" {
		return s.repo.GetAll()
	}
	return s.repo.GetByNetworkGroupID(groupID)
}

// Get Devices By Location Type
func (s *DeviceService) GetByLocationType(locationType string) ([]models.Device, error) {
	if locationType == "" {
		return s.repo.GetAll()
	}
	return s.repo.GetByLocationType(locationType)
}

// Update Device Location Type
func (s *DeviceService) UpdateLocationType(id uuid.UUID, locationType string) error {
	device, err := s.repo.GetByDeviceID(id)
	if err != nil || device == nil {
		return errors.New("device not found")
	}
	device.LocationType = locationType
	return s.repo.Update(device)
}

// Find Device By IP
func (s *DeviceService) FindDeviceByIP(ip string) (*models.Device, error) {
	return s.repo.FindDeviceByIP(ip)
}

// Get Device By UUID
func (s *DeviceService) GetByID(id uuid.UUID) (*models.Device, error) {
	return s.repo.GetByID(id)
}

// Delete Device By UUID
func (s *DeviceService) DeleteDevice(id uuid.UUID) error {
	return s.repo.Delete(id)
}
