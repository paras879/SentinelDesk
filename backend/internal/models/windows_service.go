package models

import (
	"github.com/google/uuid"
)

type WindowsService struct {
	BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_device_service" json:"device_id"`
	ServiceName string `gorm:"size:255;not null;uniqueIndex:idx_device_service" json:"service_name"`

	DisplayName     string `gorm:"size:255" json:"display_name"`
	Status          string `gorm:"size:50" json:"status"`
	StartType       string `gorm:"size:50" json:"start_type"`
	ExecutablePath  string `gorm:"size:500" json:"executable_path"`
	PID             int32  `json:"pid"`
	ServiceAccount  string `gorm:"size:255" json:"service_account"`
	Description     string `gorm:"size:1000" json:"description"`
	CanStop         bool   `json:"can_stop"`
	CanPause        bool   `json:"can_pause"`
	AcceptShutdown  bool   `json:"accept_shutdown"`

	Device Device `gorm:"foreignKey:DeviceID;references:ID"`
}
