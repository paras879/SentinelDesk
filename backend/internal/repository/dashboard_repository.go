package repository

import (
	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

type DashboardRepository struct{}

func NewDashboardRepository() *DashboardRepository {
	return &DashboardRepository{}
}

type DashboardSummary struct {
	TotalDevices   int64 `json:"totalDevices"`
	OnlineDevices  int64 `json:"onlineDevices"`
	OfflineDevices int64 `json:"offlineDevices"`
}

// FindDeviceByIP returns the most recently seen device with the given IP address.
func (r *DashboardRepository) FindDeviceByIP(ip string) (*models.Device, error) {
	var device models.Device
	err := database.DB.
		Where("ip_address = ?", ip).
		Order("last_seen DESC NULLS LAST").
		First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DashboardRepository) GetSummary() (*DashboardSummary, error) {

	var total int64
	var online int64
	var offline int64

	if err := database.DB.Model(&models.Device{}).
		Count(&total).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Model(&models.Device{}).
		Where("status = ?", "online").
		Count(&online).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Model(&models.Device{}).
		Where("status = ?", "offline").
		Count(&offline).Error; err != nil {
		return nil, err
	}

	return &DashboardSummary{
		TotalDevices:   total,
		OnlineDevices:  online,
		OfflineDevices: offline,
	}, nil
}
