package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "super-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned an error: %v", err)
	}

	gotUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() returned an error: %v", err)
	}

	if gotUserID != userID {
		t.Fatalf("ValidateJWT() returned user ID %v, want %v", gotUserID, userID)
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned an error: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("ValidateJWT() expected an error with the wrong secret")
	}
}

func TestValidateJWTExpired(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "secret", -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned an error: %v", err)
	}

	_, err = ValidateJWT(token, "secret")
	if err == nil {
		t.Fatal("ValidateJWT() expected an error for an expired token")
	}
}
