package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	jwtsec "flashquest/pkg/security/jwt"
)

func TestSignAndParseRS256Token(t *testing.T) {
	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)

	tokenStr, err := jwtsec.SignAccessToken(jwtsec.AuthClaims{UserID: 1, Role: "Aluno"}, 72*time.Hour, privateKeyPEM)
	if err != nil {
		t.Fatalf("expected no error signing token, got %v", err)
	}

	claims, err := jwtsec.ParseAccessToken(tokenStr, publicKeyPEM)
	if err != nil {
		t.Fatalf("expected no error parsing token, got %v", err)
	}

	if claims.UserID != 1 {
		t.Fatalf("expected user id 1, got %d", claims.UserID)
	}

	if claims.Role != "Aluno" {
		t.Fatalf("expected role Aluno, got %q", claims.Role)
	}

	if claims.IssuedAt == 0 || claims.ExpiresAt == 0 {
		t.Fatal("expected iat and exp to be set")
	}

	ttl := time.Unix(claims.ExpiresAt, 0).Sub(time.Unix(claims.IssuedAt, 0))
	if ttl < 71*time.Hour || ttl > 73*time.Hour {
		t.Fatalf("expected ttl around 72h, got %s", ttl)
	}
}

func TestParseAccessTokenRejectsInvalidSignature(t *testing.T) {
	privateKeyPEM, _ := generateTestKeyPair(t)
	_, wrongPublicKeyPEM := generateTestKeyPair(t)

	tokenStr, err := jwtsec.SignAccessToken(jwtsec.AuthClaims{UserID: 1, Role: "Aluno"}, time.Hour, privateKeyPEM)
	if err != nil {
		t.Fatalf("expected no error signing token, got %v", err)
	}

	if _, err := jwtsec.ParseAccessToken(tokenStr, wrongPublicKeyPEM); err == nil {
		t.Fatal("expected signature validation to fail")
	}
}

func generateTestKeyPair(t *testing.T) (string, string) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	privateDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateDER})

	publicDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to encode RSA public key: %v", err)
	}

	publicPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})

	return string(privatePEM), string(publicPEM)
}
