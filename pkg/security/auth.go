package security

import (
	"time"

	"coderero.dev/projects/go/gin/hello/models"
)

// The function GenerateAuthTokens generates access and refresh tokens for a user.
func GenerateAuthTokens(obj *models.User) (string, string) {
	// The subject of the token is the email of the user.
	sub := obj.Email

	// The current time is used to set the expiration time of the tokens.
	currentTime := time.Now()

	// The access token expires in 5 minutes and the refresh token expires in 24 hours from the
	// current time. Both tokens are signed with the secret key and given expiration times.
	accessToken := GenerateToken(sub, currentTime.Add(time.Minute*5))
	refreshToken := GenerateToken(sub, currentTime.Add(time.Hour*24))
	return accessToken, refreshToken
}
