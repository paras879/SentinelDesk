package services

import (
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/repository"
)

type SoftwareInventoryService struct {
	repo *repository.SoftwareInventoryRepository
}

func NewSoftwareInventoryService() *SoftwareInventoryService {
	return &SoftwareInventoryService{
		repo: repository.NewSoftwareInventoryRepository(),
	}
}

func (s *SoftwareInventoryService) Replace(deviceID uuid.UUID, software []models.SoftwareInventory) error {
	return s.repo.Replace(deviceID, software)
}

func (s *SoftwareInventoryService) GetByDeviceID(deviceID uuid.UUID) ([]models.SoftwareInventory, error) {
	return s.repo.GetByDeviceID(deviceID)
}
