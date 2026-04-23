package password_test

import (
	"testing"

	"flashquest/pkg/security/password"
)

func TestHashAndComparePassword(t *testing.T) {
	hash, err := password.Hash("Secret@123")
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	if hash == "Secret@123" {
		t.Fatal("expected hash to differ from raw password")
	}

	if err := password.Compare(hash, "Secret@123"); err != nil {
		t.Fatalf("expected password comparison success, got %v", err)
	}

	if err := password.Compare(hash, "wrong"); err == nil {
		t.Fatal("expected password comparison to fail for wrong password")
	}
}
