package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT_Success(t *testing.T) {
	token, err := GenerateJWT("user-1", "resident", "test@example.com", "A-101")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token string")
	}
}

func TestVerifyJwt_ValidToken(t *testing.T) {
	tokenStr, err := GenerateJWT("user-2", "officer", "officer@example.com", "B-202")
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}

	token, err := VerifyJwt(tokenStr)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if !token.Valid {
		t.Fatal("expected token to be valid")
	}
}

func TestVerifyJwt_InvalidSignature(t *testing.T) {
	tokenStr, _ := GenerateJWT("user-3", "resident", "user3@example.com", "C-303")

	badToken := tokenStr[:len(tokenStr)-1] + "x"

	_, err := VerifyJwt(badToken)
	if err == nil {
		t.Fatal("expected error for invalid signature, got nil")
	}
}

func TestVerifyJwt_InvalidSigningMethod(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": "user-4",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	// Use none signing method
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, _ := token.SignedString(jwtSecret)

	_, err := VerifyJwt(tokenStr)
	if err == nil {
		t.Fatal("expected error for invalid signing method, got nil")
	}
}

func TestVerifyJwt_ExpiredToken(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": "user-5",
		"exp":     time.Now().Add(-time.Hour).Unix(), // expired
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(jwtSecret)

	_, err := VerifyJwt(tokenStr)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestVerifyJwt_MalformedToken(t *testing.T) {
	_, err := VerifyJwt("not-a-token")
	if err == nil {
		t.Fatal("expected error for malformed token, got nil")
	}

	_, err = VerifyJwt("")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
