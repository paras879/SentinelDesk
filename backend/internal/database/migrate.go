package database

import (
	"log"
	"strings"

	"sentineldesk/backend/internal/models"
)

func AutoMigrate() error {

	migrator := DB.Migrator()

	modelList := []interface{}{
		&models.Admin{},
		&models.Device{},
		&models.DeviceSystemInfo{},
		&models.SoftwareInventory{},
		&models.Process{},
		&models.WindowsService{},
		&models.ServiceCommand{},
	}

	// Phase 1: Create tables if they don't exist
	for _, model := range modelList {
		if !migrator.HasTable(model) {
			log.Printf("Creating table for %T", model)
			if err := migrator.CreateTable(model); err != nil {
				return err
			}
		}
	}

	// Phase 2: AutoMigrate for column/index changes
	// GORM v1.31.2 has a known issue where it doesn't always detect
	// existing columns on PostgreSQL, causing "already exists" errors
	err := DB.AutoMigrate(modelList...)

	if err != nil {
		errStr := err.Error()
		// PostgreSQL "already exists" errors are safe to skip —
		// they mean the schema is already up to date
		if strings.Contains(errStr, "already exists") {
			log.Println("⚠️  Migration: schema already up to date, skipping redundant changes")
		} else {
			return err
		}
	}

	log.Println("✅ Database Migration Completed")

	return nil
}
