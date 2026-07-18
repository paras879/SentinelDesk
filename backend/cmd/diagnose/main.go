package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"sentineldesk/backend/internal/config"
	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

func main() {
	cfg := config.LoadConfig()

	if err := database.Connect(cfg); err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// 1. Total admins
	var total int64
	database.DB.Model(&models.Admin{}).Count(&total)
	fmt.Printf("Total admins: %d\n\n", total)

	// 2. List all admins
	var admins []models.Admin
	database.DB.Find(&admins)
	for _, a := range admins {
		fmt.Printf("  ID:       %s\n", a.ID)
		fmt.Printf("  Name:     %s\n", a.Name)
		fmt.Printf("  Email:    %s\n", a.Email)
		fmt.Printf("  Role:     %s\n", a.Role)
		fmt.Printf("  Hash len: %d\n", len(a.Password))
		fmt.Printf("  Hash pre: %s\n", a.Password[:min(len(a.Password), 30)])
		fmt.Println()
	}

	// 3. Find specific admin
	var admin models.Admin
	err := database.DB.Where("email = ?", "admin@sentineldesk.com").First(&admin).Error
	if err != nil {
		fmt.Println("admin@sentineldesk.com: NOT FOUND")
		fmt.Println("Root cause: admin with that email does not exist in the database.")

		fmt.Println("\nExisting admin emails:")
		for _, a := range admins {
			fmt.Printf("  - %s\n", a.Email)
		}
		return
	}

	fmt.Println("admin@sentineldesk.com: FOUND")
	fmt.Printf("  Name:      %s\n", admin.Name)
	fmt.Printf("  Hash:      %s (%d bytes)\n", admin.Password[:min(len(admin.Password), 20)], len(admin.Password))

	// 4. Test bcrypt with "Admin@123"
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte("Admin@123"))
	if err != nil {
		fmt.Printf("bcrypt.CompareHashAndPassword FAILED: %v\n", err)
		fmt.Println("Root cause: stored password hash does not match 'Admin@123'.")

		fmt.Println("\nResetting password to Admin@123 ...")
		hash, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}
		result := database.DB.Model(&admin).Update("password", string(hash))
		if result.Error != nil {
			log.Fatal("Failed to update password:", result.Error)
		}
		if result.RowsAffected != 1 {
			log.Fatal("Expected 1 row affected, got", result.RowsAffected)
		}
		fmt.Println("Password reset successfully!")
	} else {
		fmt.Println("bcrypt.CompareHashAndPassword: PASSED")
		fmt.Println("Password hash matches 'Admin@123'. Login should work.")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
