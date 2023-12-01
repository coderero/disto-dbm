package security

import (
	"coderero.dev/projects/go/gin/hello/pkg/utils"
)

// The function HashPassword takes a raw password as input and returns its hashed version using bcrypt
// algorithm.
func HashPassword(password_raw string) (string, error) {
	hashed_password, err := utils.CreatePassword(password_raw)
	if err != nil {
		return "", err
	}
	return hashed_password, nil
}

// The function `ComparePassword` compares a raw password with a hashed password and returns true if
// they match, and false otherwise.
func ComparePassword(password_raw string, password_hashed string) bool {
	err := utils.ComparePassword(password_raw, password_hashed)
	if err != nil {
		return false
	}
	return true
}
