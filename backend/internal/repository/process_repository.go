package repository

import (
	"github.com/google/uuid"

	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

type ProcessRepository struct{}

func NewProcessRepository() *ProcessRepository {
	return &ProcessRepository{}
}

func (r *ProcessRepository) Replace(deviceID uuid.UUID, processes []models.Process) error {

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("device_id = ?", deviceID).
		Delete(&models.Process{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range processes {
		processes[i].DeviceID = deviceID
	}

	if len(processes) > 0 {
		if err := tx.Create(&processes).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *ProcessRepository) GetByDeviceID(deviceID uuid.UUID) ([]models.Process, error) {

	var processes []models.Process

	err := database.DB.
		Where("device_id = ?", deviceID).
		Order("name ASC").
		Find(&processes).Error

	if err != nil {
		return nil, err
	}

	return processes, nil
}
