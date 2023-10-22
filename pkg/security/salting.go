package security

import (
	"encoding/base64"
	"os"

	"github.com/joho/godotenv"
)

var salt []byte

func init() {
	godotenv.Load()
	salt = []byte(os.Getenv("ENCRYPTION_KEY"))
}

func salting(password_raw string) string {

	password := []byte(password_raw)

	for i := 0; i < len(password); i++ {
		password[i] = password[i] ^ salt[i%len(salt)]
	}

	return base64.StdEncoding.EncodeToString(password)
}

func deSalting(password_hashed string) (string, error) {

	password, err := base64.StdEncoding.DecodeString(password_hashed)

	if err != nil {
		return "", err
	}

	for i := 0; i < len(password); i++ {
		password[i] = password[i] ^ salt[i%len(salt)]
	}

	return string(password), nil
}
