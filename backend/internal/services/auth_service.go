package services

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"sentineldesk/backend/internal/models"
	"sentineldesk/backend/internal/repository"
)

type AuthService struct {
	repo *repository.AdminRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo: repository.NewAdminRepository(),
	}
}

func (s *AuthService) Register(name, email, password string) error {

	// Email already exists?
	existing, _ := s.repo.GetByEmail(email)
	if existing != nil {
		return errors.New("email already exists")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &models.Admin{
		Name:     name,
		Email:    email,
		Password: string(hash),
		Role:     "admin",
	}

	return s.repo.Create(admin)
}

// ----------------------
// Login
// ----------------------

func (s *AuthService) Login(email, password string) (*models.Admin, error) {

	log.Printf("[DEBUG] Login attempt for email: %s", email)

	admin, err := s.repo.GetByEmail(email)
	if err != nil {
		log.Printf("[DEBUG] Admin not found for email: %s — err: %v", email, err)
		return nil, errors.New("invalid email or password")
	}

	log.Printf("[DEBUG] Admin found: %s (%s)", admin.Name, admin.Email)

	err = bcrypt.CompareHashAndPassword(
		[]byte(admin.Password),
		[]byte(password),
	)

	if err != nil {
		log.Printf("[DEBUG] bcrypt comparison FAILED for %s — err: %v", email, err)
		return nil, errors.New("invalid email or password")
	}

	log.Printf("[DEBUG] bcrypt comparison PASSED for %s", email)

	return admin, nil
}
