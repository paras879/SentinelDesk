package repository

import (
	"github.com/google/uuid"

	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

type WindowsServiceRepository struct{}

func NewWindowsServiceRepository() *WindowsServiceRepository {
	return &WindowsServiceRepository{}
}

func (r *WindowsServiceRepository) Replace(deviceID uuid.UUID, services []models.WindowsService) error {

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("device_id = ?", deviceID).
		Delete(&models.WindowsService{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range services {
		services[i].DeviceID = deviceID
	}

	if len(services) > 0 {
		if err := tx.Create(&services).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *WindowsServiceRepository) GetByDeviceID(deviceID uuid.UUID) ([]models.WindowsService, error) {

	var services []models.WindowsService

	err := database.DB.
		Where("device_id = ?", deviceID).
		Order("display_name ASC").
		Find(&services).Error

	if err != nil {
		return nil, err
	}

	return services, nil
}
