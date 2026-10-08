package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "super-secret"

	token, err := MakeJWT(userID, secret)
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

	token, err := MakeJWT(userID, "correct-secret")
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
	secret := "secret"

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
		ID:        uuid.NewString(),
		Audience:  []string{userID.String()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("SignedString() returned an error: %v", err)
	}

	_, err = ValidateJWT(signedToken, secret)
	if err == nil {
		t.Fatal("ValidateJWT() expected an error for an expired token")
	}
}
