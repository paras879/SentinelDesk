package repository

import (
	"github.com/google/uuid"

	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

type SoftwareInventoryRepository struct{}

func NewSoftwareInventoryRepository() *SoftwareInventoryRepository {
	return &SoftwareInventoryRepository{}
}

func (r *SoftwareInventoryRepository) Replace(deviceID uuid.UUID, software []models.SoftwareInventory) error {

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("device_id = ?", deviceID).
		Delete(&models.SoftwareInventory{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range software {
		software[i].DeviceID = deviceID
	}

	if len(software) > 0 {
		if err := tx.Create(&software).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *SoftwareInventoryRepository) GetByDeviceID(deviceID uuid.UUID) ([]models.SoftwareInventory, error) {

	var software []models.SoftwareInventory

	err := database.DB.
		Where("device_id = ?", deviceID).
		Order("name ASC").
		Find(&software).Error

	if err != nil {
		return nil, err
	}

	return software, nil
}
