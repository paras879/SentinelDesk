package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"sentineldesk/backend/internal/models"
)

func SeedAdmin() {
	hash, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Failed to hash admin password:", err)
	}

	var admin models.Admin
	result := DB.Where("email = ?", "admin@sentineldesk.com").First(&admin)

	if result.Error == nil {
		admin.Password = string(hash)
		if err := DB.Save(&admin).Error; err != nil {
			log.Fatal("❌ Failed to update admin password:", err)
		}
		log.Println("✅ Admin password updated (admin@sentineldesk.com)")
		return
	}

	admin = models.Admin{
		Name:     "Admin",
		Email:    "admin@sentineldesk.com",
		Password: string(hash),
		Role:     "admin",
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Fatal("❌ Failed to seed admin:", err)
	}

	log.Println("✅ Default admin created (admin@sentineldesk.com)")
}
