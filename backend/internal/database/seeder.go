package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"sentineldesk/backend/internal/models"
)

func SeedAdmin() {
	var count int64
	DB.Model(&models.Admin{}).Count(&count)

	if count > 0 {
		log.Println("✅ Admin already exists, skipping seed")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Failed to hash admin password:", err)
	}

	admin := &models.Admin{
		Name:     "Admin",
		Email:    "admin@sentineldesk.com",
		Password: string(hash),
		Role:     "admin",
	}

	if err := DB.Create(admin).Error; err != nil {
		log.Fatal("❌ Failed to seed admin:", err)
	}

	log.Println("✅ Default admin created (admin@sentineldesk.com)")
}
