package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "correct-horse-battery-staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned an error: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() returned an empty hash")
	}

	if hash == password {
		t.Fatal("HashPassword() returned the plain-text password")
	}
}

func TestCheckPasswordHashCorrectPassword(t *testing.T) {
	password := "correct-horse-battery-staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned an error: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() returned an error: %v", err)
	}

	if !match {
		t.Fatal("CheckPasswordHash() returned false for the correct password")
	}
}

func TestCheckPasswordHashWrongPassword(t *testing.T) {
	password := "correct-horse-battery-staple"
	wrongPassword := "wrong-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned an error: %v", err)
	}

	match, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() returned an error: %v", err)
	}

	if match {
		t.Fatal("CheckPasswordHash() returned true for the wrong password")
	}
}

func TestCheckPasswordHashInvalidHash(t *testing.T) {
	_, err := CheckPasswordHash("password", "not-a-valid-argon2-hash")
	if err == nil {
		t.Fatal("CheckPasswordHash() expected an error for an invalid hash")
	}
}
