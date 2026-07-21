package repository

import (
	"fmt"
	"log"

	"sentineldesk/backend/internal/database"
	"sentineldesk/backend/internal/models"
)

type AdminRepository struct{}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{}
}

func (r *AdminRepository) Create(admin *models.Admin) error {
	return database.DB.Create(admin).Error
}

func (r *AdminRepository) GetByEmail(email string) (*models.Admin, error) {
	var admin models.Admin

	sql := fmt.Sprintf("SELECT * FROM \"admins\" WHERE email = '%s' ORDER BY \"admins\".\"id\" LIMIT 1", email)
	log.Printf("[DEBUG] SQL query: %s", sql)

	err := database.DB.
		Where("email = ?", email).
		First(&admin).Error

	if err != nil {
		log.Printf("[DEBUG] GetByEmail result — error: %v", err)
		return nil, err
	}

	log.Printf("[DEBUG] GetByEmail result — found: ID=%s Name=%s Email=%s", admin.ID, admin.Name, admin.Email)

	return &admin, nil
}

func (r *AdminRepository) GetByID(id string) (*models.Admin, error) {
	var admin models.Admin

	err := database.DB.
		Where("id = ?", id).
		First(&admin).Error

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *AdminRepository) Update(admin *models.Admin) error {
	return database.DB.Save(admin).Error
}
