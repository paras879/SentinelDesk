package utils

import (
	"os"
	"testing"

	"github.com/google/uuid"

	"sentineldesk/backend/internal/models"
)

const testJWTSecret = "test-secret-key-for-unit-tests-min32chars"

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", testJWTSecret)
	os.Exit(m.Run())
}

func newTestAdmin() *models.Admin {
	return &models.Admin{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Name:      "Test Admin",
		Email:     "test@example.com",
		Password:  "hashed",
		Role:      "admin",
	}
}

func TestGenerateJWT_NonEmpty(t *testing.T) {
	admin := newTestAdmin()
	token, err := GenerateJWT(admin)
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateJWT returned empty token")
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	admin := newTestAdmin()
	token, err := GenerateJWT(admin)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT returned error for valid token: %v", err)
	}

	if email, ok := claims["email"].(string); !ok || email != admin.Email {
		t.Errorf("expected email %q, got %q", admin.Email, email)
	}
	if role, ok := claims["role"].(string); !ok || role != admin.Role {
		t.Errorf("expected role %q, got %q", admin.Role, role)
	}
	if id, ok := claims["id"].(string); !ok || id != admin.ID.String() {
		t.Errorf("expected id %q, got %q", admin.ID.String(), id)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := ValidateJWT("this.is.not.a.valid.jwt.token")
	if err == nil {
		t.Fatal("expected error for tampered/invalid token, got nil")
	}
}

func TestValidateJWT_EmptyString(t *testing.T) {
	_, err := ValidateJWT("")
	if err == nil {
		t.Fatal("expected error for empty token string, got nil")
	}
}
