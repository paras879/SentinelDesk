package services

import (
	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/repository"
)

type DeviceSystemInfoService struct {
	repo *repository.DeviceSystemInfoRepository
}

func NewDeviceSystemInfoService() *DeviceSystemInfoService {
	return &DeviceSystemInfoService{
		repo: repository.NewDeviceSystemInfoRepository(),
	}
}

// Save or Update System Information
func (s *DeviceSystemInfoService) Save(info *models.DeviceSystemInfo) error {
	return s.repo.Upsert(info)
}

// Get System Information By Device ID
func (s *DeviceSystemInfoService) GetByDeviceID(deviceID uuid.UUID) (*models.DeviceSystemInfo, error) {
	return s.repo.GetByDeviceID(deviceID)
}
