package models

import (
	"time"

	"github.com/google/uuid"
)

type DeviceSystemInfo struct {
	BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	CPUName  string  `gorm:"size:255" json:"cpu_name"`
	CPUUsage float64 `json:"cpu_usage"`

	TotalRAM uint64 `json:"total_ram"`
	UsedRAM  uint64 `json:"used_ram"`
	FreeRAM  uint64 `json:"free_ram"`

	TotalDisk uint64 `json:"total_disk"`
	UsedDisk  uint64 `json:"used_disk"`
	FreeDisk  uint64 `json:"free_disk"`

	OSVersion string `gorm:"size:255" json:"os_version"`

	LocalIP  string `gorm:"size:100" json:"local_ip"`
	PublicIP string `gorm:"size:100" json:"public_ip"`

	MACAddress string `gorm:"size:100" json:"mac_address"`

	Uptime int64 `json:"uptime"`

	LastBoot time.Time `json:"last_boot"`

	AgentVersion string `gorm:"size:50" json:"agent_version"`

	Device Device `gorm:"foreignKey:DeviceID;references:ID"`
}
