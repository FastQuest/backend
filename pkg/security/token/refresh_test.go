package token_test

import (
	"testing"

	"flashquest/pkg/security/token"
)

func TestGenerateAndHashRefreshToken(t *testing.T) {
	raw, hashed, err := token.GenerateAndHash()
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	if raw == "" {
		t.Fatal("expected raw token not empty")
	}

	if hashed == "" {
		t.Fatal("expected hashed token not empty")
	}

	if raw == hashed {
		t.Fatal("expected raw token and hashed token to differ")
	}

	if !token.VerifyHash(raw, hashed) {
		t.Fatal("expected hash verification to succeed")
	}
}

func TestVerifyHashRejectsWrongToken(t *testing.T) {
	raw, hashed, err := token.GenerateAndHash()
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	if token.VerifyHash(raw+"x", hashed) {
		t.Fatal("expected hash verification to fail for wrong token")
	}
}
