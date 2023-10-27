package security

import (
	"golang.org/x/crypto/bcrypt"
)

// The function HashPassword takes a raw password as input and returns its hashed version using bcrypt
// algorithm.
func HashPassword(password_raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password_raw), 7)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// The function `ComparePassword` compares a raw password with a hashed password and returns true if
// they match, and false otherwise.
func ComparePassword(password_raw string, password_hashed string) bool {
	err1 := bcrypt.CompareHashAndPassword([]byte(password_hashed), []byte(password_raw))
	if err1 != nil {
		return false
	}
	return true
}
