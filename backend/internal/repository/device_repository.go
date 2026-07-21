package repository

import (
	"time"

	"github.com/google/uuid"

	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

type DeviceRepository struct{}

func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{}
}

// Create Device
func (r *DeviceRepository) Create(device *models.Device) error {
	return database.DB.Create(device).Error
}

// Get Device By Device ID (UUID primary key)
func (r *DeviceRepository) GetByDeviceID(id uuid.UUID) (*models.Device, error) {

	var device models.Device

	err := database.DB.
		Where("id = ?", id).
		First(&device).Error

	if err != nil {
		return nil, err
	}

	return &device, nil
}

// Get Devices By Network ID (legacy subnet-based filtering)
func (r *DeviceRepository) GetByNetworkID(networkID string) ([]models.Device, error) {
	var devices []models.Device
	err := database.DB.
		Where("connected_subnet = ?", networkID).
		Order("created_at DESC").
		Find(&devices).Error
	return devices, err
}

// Get Devices By Network Group ID
func (r *DeviceRepository) GetByNetworkGroupID(groupID string) ([]models.Device, error) {
	var devices []models.Device
	err := database.DB.
		Where("network_group_id = ?", groupID).
		Order("created_at DESC").
		Find(&devices).Error
	return devices, err
}

// Find Device By IP Address (returns most recently seen device with this IP)
func (r *DeviceRepository) FindDeviceByIP(ip string) (*models.Device, error) {
	var device models.Device
	err := database.DB.
		Where("ip_address = ?", ip).
		Order("last_seen DESC NULLS LAST").
		First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// Get All Devices
func (r *DeviceRepository) GetAll() ([]models.Device, error) {

	var devices []models.Device

	err := database.DB.
		Order("created_at DESC").
		Find(&devices).Error

	return devices, err
}

// Update Heartbeat by Device ID
func (r *DeviceRepository) UpdateHeartbeat(
	deviceID uuid.UUID,
	ipAddress string,
	macAddress string,
	connectedSubnet string,
	defaultGateway string,
	networkGroupID string,
) error {

	return database.DB.
		Model(&models.Device{}).
		Where("id = ?", deviceID).
		Updates(map[string]interface{}{
			"status":            "online",
			"ip_address":        ipAddress,
			"mac_address":       macAddress,
			"connected_subnet":  connectedSubnet,
			"default_gateway":   defaultGateway,
			"network_group_id":  networkGroupID,
			"last_seen":         time.Now(),
		}).Error
}

// Update Device
func (r *DeviceRepository) Update(device *models.Device) error {
	return database.DB.Save(device).Error
}

// Get Device By UUID (alias for GetByDeviceID)
func (r *DeviceRepository) GetByID(id uuid.UUID) (*models.Device, error) {
	return r.GetByDeviceID(id)
}

// Delete Device
func (r *DeviceRepository) Delete(id uuid.UUID) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Manually cascade delete child records to prevent foreign key constraint errors
	tx.Where("device_id = ?", id).Delete(&models.DeviceSystemInfo{})
	tx.Where("device_id = ?", id).Delete(&models.Process{})
	tx.Where("device_id = ?", id).Delete(&models.SoftwareInventory{})
	tx.Where("device_id = ?", id).Delete(&models.WindowsService{})
	tx.Where("device_id = ?", id).Delete(&models.ServiceCommand{})

	if err := tx.Where("id = ?", id).Delete(&models.Device{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
