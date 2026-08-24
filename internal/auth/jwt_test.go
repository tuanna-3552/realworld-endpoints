package auth

import (
	"testing"
)

func TestGenerateToken_Success(t *testing.T) {
	token, err := GenerateToken(1, "test@example.com", "secret")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Errorf("expected non-empty token")
	}
}

func TestParseToken_Valid(t *testing.T) {
	secret := "mysecret"
	userID := uint(1)
	email := "test@example.com"
	token, err := GenerateToken(userID, email, secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, err := ParseToken("garbage.token.string", "secret")
	if err == nil {
		t.Errorf("expected error for invalid token, got none")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, _ := GenerateToken(1, "test@example.com", "secret1")
	_, err := ParseToken(token, "secret2")
	if err == nil {
		t.Errorf("expected error when parsing with wrong secret")
	}
}

func TestHashPassword_Success(t *testing.T) {
	password := "my-secure-password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == password {
		t.Errorf("expected hash to be different from password")
	}
	if hash == "" {
		t.Errorf("expected non-empty hash")
	}
}

func TestCheckPasswordHash_Valid(t *testing.T) {
	password := "my-secure-password"
	hash, _ := HashPassword(password)
	if !CheckPasswordHash(password, hash) {
		t.Errorf("expected CheckPasswordHash to return true for correct password")
	}
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	password := "my-secure-password"
	hash, _ := HashPassword(password)
	if CheckPasswordHash("wrong-password", hash) {
		t.Errorf("expected CheckPasswordHash to return false for incorrect password")
	}
}
