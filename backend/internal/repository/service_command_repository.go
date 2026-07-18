package repository

import (
	"time"

	"github.com/google/uuid"

	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

type ServiceCommandRepository struct{}

func NewServiceCommandRepository() *ServiceCommandRepository {
	return &ServiceCommandRepository{}
}

func (r *ServiceCommandRepository) Create(cmd *models.ServiceCommand) error {
	return database.DB.Create(cmd).Error
}

func (r *ServiceCommandRepository) GetPendingByDeviceID(deviceID uuid.UUID) ([]models.ServiceCommand, error) {

	var commands []models.ServiceCommand

	err := database.DB.
		Where("device_id = ? AND status = ?", deviceID, "pending").
		Order("created_at ASC").
		Find(&commands).Error

	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (r *ServiceCommandRepository) MarkCompleted(cmdID uuid.UUID, result, errorMessage string) error {

	return database.DB.
		Model(&models.ServiceCommand{}).
		Where("id = ?", cmdID).
		Updates(map[string]interface{}{
			"status":        "completed",
			"result":        result,
			"error_message": errorMessage,
			"executed_at":   time.Now(),
		}).Error
}
