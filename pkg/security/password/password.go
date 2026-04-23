package password

import "golang.org/x/crypto/bcrypt"

const defaultCost = bcrypt.DefaultCost

func Hash(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), defaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func Compare(hash, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
}
