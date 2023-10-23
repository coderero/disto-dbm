package security

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password_raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password_raw), 8)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(password_raw string, password_hashed string) bool {
	err1 := bcrypt.CompareHashAndPassword([]byte(password_hashed), []byte(password_raw))
	if err1 != nil {
		return false
	}
	return true
}
