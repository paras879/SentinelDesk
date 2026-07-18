package repository

import (
	"github.com/google/uuid"

	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

type DeviceSystemInfoRepository struct{}

func NewDeviceSystemInfoRepository() *DeviceSystemInfoRepository {
	return &DeviceSystemInfoRepository{}
}

// Create or Update System Information
func (r *DeviceSystemInfoRepository) Upsert(info *models.DeviceSystemInfo) error {

	var existing models.DeviceSystemInfo

	err := database.DB.
		Where("device_id = ?", info.DeviceID).
		First(&existing).Error

	if err == nil {

		info.ID = existing.ID

		return database.DB.
			Save(info).Error
	}

	return database.DB.Create(info).Error
}

// Get System Information By Device ID
func (r *DeviceSystemInfoRepository) GetByDeviceID(deviceID uuid.UUID) (*models.DeviceSystemInfo, error) {

	var info models.DeviceSystemInfo

	err := database.DB.
		Preload("Device").
		Where("device_id = ?", deviceID).
		First(&info).Error

	if err != nil {
		return nil, err
	}

	return &info, nil
}
