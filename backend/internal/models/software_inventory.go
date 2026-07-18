package models

import (
	"github.com/google/uuid"
)

type SoftwareInventory struct {
	BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_device_software" json:"device_id"`
	Name     string    `gorm:"size:255;not null;uniqueIndex:idx_device_software" json:"name"`
	Version  string    `gorm:"size:100" json:"version"`

	Publisher       string `gorm:"size:255" json:"publisher"`
	InstallDate     string `gorm:"size:50" json:"install_date"`
	InstallLocation string `gorm:"size:500" json:"install_location"`
	EstimatedSize   int64  `json:"estimated_size"`

	Device Device `gorm:"foreignKey:DeviceID;references:ID"`
}
