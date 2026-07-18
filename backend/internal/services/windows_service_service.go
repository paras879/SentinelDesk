package services

import (
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/repository"
)

type WindowsServiceService struct {
	repo *repository.WindowsServiceRepository
}

func NewWindowsServiceService() *WindowsServiceService {
	return &WindowsServiceService{
		repo: repository.NewWindowsServiceRepository(),
	}
}

func (s *WindowsServiceService) Replace(deviceID uuid.UUID, services []models.WindowsService) error {
	return s.repo.Replace(deviceID, services)
}

func (s *WindowsServiceService) GetByDeviceID(deviceID uuid.UUID) ([]models.WindowsService, error) {
	return s.repo.GetByDeviceID(deviceID)
}
