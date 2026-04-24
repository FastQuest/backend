package token

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
)

const refreshTokenBytes = 32

func Generate() (string, error) {
	buf := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func Hash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func GenerateAndHash() (string, string, error) {
	raw, err := Generate()
	if err != nil {
		return "", "", err
	}

	return raw, Hash(raw), nil
}

func VerifyHash(raw, hashed string) bool {
	rawHash := Hash(raw)
	if len(rawHash) != len(hashed) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(rawHash), []byte(hashed)) == 1
}
