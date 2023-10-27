package security

import (
	"time"

	"coderero.dev/projects/go/gin/hello/models"
)

// The function GenerateAuthTokens generates access and refresh tokens for a user.
func GenerateAuthTokens(obj *models.User) (string, string) {
	sub := obj.Email
	currentTime := time.Now()
	accessToken := GenerateToken(sub, currentTime.Add(time.Minute*5))
	refreshToken := GenerateToken(sub, currentTime.Add(time.Hour*24))
	return accessToken, refreshToken
}
