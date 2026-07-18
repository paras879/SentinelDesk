package models

import (
	"time"

	"github.com/google/uuid"
)

type Process struct {
	BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_device_pid" json:"device_id"`
	PID      int32     `gorm:"not null;uniqueIndex:idx_device_pid" json:"pid"`

	Name            string  `gorm:"size:255" json:"name"`
	ExecutablePath  string  `gorm:"size:500" json:"executable_path"`
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryBytes     uint64  `json:"memory_bytes"`
	MemoryPercent   float32 `json:"memory_percent"`
	ThreadCount     int32   `json:"thread_count"`
	HandleCount     int32   `json:"handle_count"`
	StartTime       time.Time `json:"start_time"`
	Username        string  `gorm:"size:150" json:"username"`
	Status          string  `gorm:"size:50;default:Running" json:"status"`

	Device Device `gorm:"foreignKey:DeviceID;references:ID"`
}
