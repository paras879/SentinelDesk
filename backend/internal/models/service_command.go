package models

import (
	"time"

	"github.com/google/uuid"
)

type ServiceCommand struct {
	BaseModel

	DeviceID    uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ServiceName string    `gorm:"size:255;not null" json:"service_name"`
	Action      string    `gorm:"size:50;not null" json:"action"`

	Status       string     `gorm:"size:50;default:pending" json:"status"`
	Result       string     `gorm:"size:50" json:"result"`
	ErrorMessage string     `gorm:"size:1000" json:"error_message"`
	ExecutedAt   *time.Time `json:"executed_at"`
}
