package security

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password_raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password_raw), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	hashOnly := strings.Replace(string(hash), "$2a$10$", "", 1)

	saltedHash := "$bitlone$" + salting(hashOnly)
	return saltedHash, nil
}

func ComparePassword(password_raw string, password_hashed string) bool {
	password_hashed = strings.Replace(password_hashed, "$bitlone$", "", 1)
	password_sallt, err := deSalting(password_hashed)
	if err != nil {
		return false
	}
	password_hashed = "$2a$10$" + password_sallt
	err1 := bcrypt.CompareHashAndPassword([]byte(password_hashed), []byte(password_raw))
	if err1 != nil {
		return false
	}
	return true
}
