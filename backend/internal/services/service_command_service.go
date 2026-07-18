package services

import (
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/repository"
)

type ServiceCommandService struct {
	repo *repository.ServiceCommandRepository
}

func NewServiceCommandService() *ServiceCommandService {
	return &ServiceCommandService{
		repo: repository.NewServiceCommandRepository(),
	}
}

func (s *ServiceCommandService) Create(deviceID uuid.UUID, serviceName, action string) (*models.ServiceCommand, error) {

	cmd := &models.ServiceCommand{
		DeviceID:    deviceID,
		ServiceName: serviceName,
		Action:      action,
		Status:      "pending",
	}

	return cmd, s.repo.Create(cmd)
}

func (s *ServiceCommandService) GetPendingByDeviceID(deviceID uuid.UUID) ([]models.ServiceCommand, error) {
	return s.repo.GetPendingByDeviceID(deviceID)
}

func (s *ServiceCommandService) MarkCompleted(cmdID uuid.UUID, result, errorMessage string) error {
	return s.repo.MarkCompleted(cmdID, result, errorMessage)
}
