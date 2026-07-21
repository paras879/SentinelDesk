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
		// For local testing, pick any available network group so the feature isn't disabled
		device, err := s.repo.FindDeviceByIP("192.168.1.156") // the ip of the desktop device in db
		if err == nil && device != nil && device.NetworkGroupID != "" {
			return device.NetworkGroupID
		}
		
		device2, err2 := s.repo.FindDeviceByIP("192.168.0.200") // the other ip
		if err2 == nil && device2 != nil && device2.NetworkGroupID != "" {
			return device2.NetworkGroupID
		}

		return "6a1072bc-1234-abcd-5678-abcdefabcdef" // fake ID to enable the UI for local testing
	}

	device, err := s.repo.FindDeviceByIP(adminIP)
	if err != nil || device == nil {
		return ""
	}

	return device.NetworkGroupID
}
