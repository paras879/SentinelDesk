package scheduler

import (
	"log"
	"time"

	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

func StartDeviceStatusScheduler() {

	log.Println("Device Status Scheduler Started")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {

		var devices []models.Device

		err := database.DB.Find(&devices).Error
		if err != nil {
			log.Println("Scheduler Error:", err)
			continue
		}

		for _, device := range devices {

			if device.LastSeen == nil {
				continue
			}

			if time.Since(*device.LastSeen) > 3*time.Minute {

				err := database.DB.Model(&device).
					Update("status", "offline").Error

				if err != nil {
					log.Println("Failed to update device status:", err)
				}
			}
		}
	}
}
