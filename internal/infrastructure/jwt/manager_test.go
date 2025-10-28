package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateToken(t *testing.T) {
	manager := NewManager("secret", time.Hour, "issuer")

	token, err := manager.GenerateToken("user-id")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	subject, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}
	if subject != "user-id" {
		t.Fatalf("expected subject user-id got %s", subject)
	}
}

func TestValidateTokenErrors(t *testing.T) {
	manager := NewManager("secret", time.Second, "issuer")

	// malformed token
	if _, err := manager.ValidateToken("bad.token"); err == nil {
		t.Fatalf("expected error for malformed token")
	}

	// expired token
	short := NewManager("secret", -time.Second, "issuer")
	token, err := short.GenerateToken("user")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if _, err := manager.ValidateToken(token); err == nil {
		t.Fatalf("expected expiration error")
	}
}

func TestValidateTokenNotYetValid(t *testing.T) {
	manager := NewManager("secret", time.Hour, "issuer")
	claims := jwt.RegisteredClaims{
		Subject:   "user",
		Issuer:    "issuer",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	if _, err := manager.ValidateToken(str); err == nil {
		t.Fatalf("expected not yet valid error")
	}
}
