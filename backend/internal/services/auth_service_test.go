package services

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcrypt_CorrectPassword(t *testing.T) {
	password := "SuperSecretPassword123!"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword failed: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		t.Errorf("CompareHashAndPassword failed for correct password: %v", err)
	}
}

func TestBcrypt_WrongPassword(t *testing.T) {
	password := "CorrectPassword"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword failed: %v", err)
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte("WrongPassword"))
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}

func TestBcrypt_EmptyPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("somepassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword failed: %v", err)
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(""))
	if err == nil {
		t.Error("expected error for empty password comparison, got nil")
	}
}

func TestBcrypt_HashIsNotPlaintext(t *testing.T) {
	password := "plaintext"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword failed: %v", err)
	}
	if string(hash) == password {
		t.Error("hash must not equal the original plaintext password")
	}
}
