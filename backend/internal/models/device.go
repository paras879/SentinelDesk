package models

import "time"

type Device struct {
	BaseModel

	DeviceName string `gorm:"size:150;not null"`
	Hostname   string `gorm:"size:150"`
	Username   string `gorm:"size:100"`
	OS         string `gorm:"size:100"`
	OSVersion  string `gorm:"size:100"`
	IPAddress  string `gorm:"size:50"`
	MACAddress string `gorm:"size:20"`

	ConnectedSubnet string `gorm:"size:50"`
	NetworkAdapters string `gorm:"type:jsonb"`
	DefaultGateway  string `gorm:"size:45"`
	NetworkGroupID  string `gorm:"type:varchar(64);index"`
	LocationType    string `gorm:"size:50;default:'Unassigned'"`

	Status string `gorm:"default:offline"`

	LastSeen *time.Time
}
