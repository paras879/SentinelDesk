package database

import (
	"log"

	"sentineldesk/backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var dsn string

	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
		log.Println("ℹ️  Connecting using DATABASE_URL")
	} else {
		dsn = "host=" + cfg.DBHost +
			" user=" + cfg.DBUser +
			" password=" + cfg.DBPassword +
			" dbname=" + cfg.DBName +
			" port=" + cfg.DBPort +
			" sslmode=" + cfg.DBSSLMode
		log.Println("ℹ️  Connecting using individual DB_* variables (local development)")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	log.Println("✅ PostgreSQL Connected Successfully")

	return nil
}
