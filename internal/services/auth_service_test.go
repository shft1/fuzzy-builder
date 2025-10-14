package services

import (
	"testing"
	"time"
)

func TestPasswordHasher(t *testing.T) {
	h := NewPasswordHasher()
	hash, err := h.HashPassword("secret")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if err := h.VerifyPassword(hash, "secret"); err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if err := h.VerifyPassword(hash, "bad"); err == nil {
		t.Fatalf("expected verify to fail for wrong password")
	}
}

func TestJWTIssuer(t *testing.T) {
	j := NewJWTIssuer("test-secret", "issuer", time.Minute)
	token, err := j.GenerateToken(42, "manager")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	uid, role, err := j.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if uid != 42 || role != "manager" {
		t.Fatalf("unexpected claims: uid=%d role=%s", uid, role)
	}
}
