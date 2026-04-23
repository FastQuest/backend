package jwt_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
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

func TestParseAccessTokenRejectsTokenWithoutExp(t *testing.T) {
	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)

	tokenStr := signCustomToken(t, privateKeyPEM, map[string]any{
		"sub":  uint(1),
		"role": "Aluno",
		"iat":  time.Now().Unix(),
	})

	if _, err := jwtsec.ParseAccessToken(tokenStr, publicKeyPEM); err == nil {
		t.Fatal("expected parsing token without exp to fail")
	}
}

func TestParseAccessTokenRejectsTokenMissingIdentityClaim(t *testing.T) {
	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)

	tokenStr := signCustomToken(t, privateKeyPEM, map[string]any{
		"role": "Aluno",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})

	if _, err := jwtsec.ParseAccessToken(tokenStr, publicKeyPEM); err == nil {
		t.Fatal("expected parsing token without sub to fail")
	}
}

func TestParseAccessTokenRejectsTokenMissingOrEmptyRole(t *testing.T) {
	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)
	now := time.Now()

	testCases := []struct {
		name    string
		payload map[string]any
	}{
		{
			name: "missing role",
			payload: map[string]any{
				"sub": uint(1),
				"exp": now.Add(time.Hour).Unix(),
				"iat": now.Unix(),
			},
		},
		{
			name: "empty role",
			payload: map[string]any{
				"sub":  uint(1),
				"role": "",
				"exp":  now.Add(time.Hour).Unix(),
				"iat":  now.Unix(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenStr := signCustomToken(t, privateKeyPEM, tc.payload)

			if _, err := jwtsec.ParseAccessToken(tokenStr, publicKeyPEM); err == nil {
				t.Fatal("expected parsing token with invalid role to fail")
			}
		})
	}
}

func TestParseAccessTokenRejectsExpiredToken(t *testing.T) {
	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)

	tokenStr, err := jwtsec.SignAccessToken(jwtsec.AuthClaims{UserID: 1, Role: "Aluno"}, -time.Minute, privateKeyPEM)
	if err != nil {
		t.Fatalf("expected no error signing token, got %v", err)
	}

	if _, err := jwtsec.ParseAccessToken(tokenStr, publicKeyPEM); err == nil {
		t.Fatal("expected expired token parsing to fail")
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

func signCustomToken(t *testing.T, privateKeyPEM string, payload map[string]any) string {
	t.Helper()

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		t.Fatal("failed to decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse private key: %v", err)
	}

	headerJSON, err := json.Marshal(map[string]string{"alg": "RS256", "typ": "JWT"})
	if err != nil {
		t.Fatalf("failed to marshal header: %v", err)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	header := base64.RawURLEncoding.EncodeToString(headerJSON)
	body := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := header + "." + body

	hash := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature)
}
