package services

import (
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/repository"
)

type ProcessService struct {
	repo *repository.ProcessRepository
}

func NewProcessService() *ProcessService {
	return &ProcessService{
		repo: repository.NewProcessRepository(),
	}
}

func (s *ProcessService) Replace(deviceID uuid.UUID, processes []models.Process) error {
	return s.repo.Replace(deviceID, processes)
}

func (s *ProcessService) GetByDeviceID(deviceID uuid.UUID) ([]models.Process, error) {
	return s.repo.GetByDeviceID(deviceID)
}
