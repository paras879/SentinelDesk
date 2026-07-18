package services

import (
	"sentineldesk/backend/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService() *DashboardService {
	return &DashboardService{
		repo: repository.NewDashboardRepository(),
	}
}

func (s *DashboardService) GetSummary() (*repository.DashboardSummary, error) {
	return s.repo.GetSummary()
}

// GetAdminNetworkGroupID returns the NetworkGroupID of the device that has the given IP.
// If no device has that IP, returns empty string → frontend falls back to All Devices.
func (s *DashboardService) GetAdminNetworkGroupID(adminIP string) string {
	if adminIP == "" || adminIP == "127.0.0.1" || adminIP == "::1" {
		return ""
	}

	device, err := s.repo.FindDeviceByIP(adminIP)
	if err != nil || device == nil {
		return ""
	}

	return device.NetworkGroupID
}
