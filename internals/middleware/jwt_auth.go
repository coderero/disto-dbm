package middleware

import (
	"net/http"
	"strings"

	"coderero.dev/projects/go/gin/hello/pkg/security"
	types "coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

// The JWTAuthMiddleWare function is a middleware that handles authentication using JSON Web Tokens
// (JWT) in a Go web application.
func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the token from the header
		token := c.Request.Header.Get("Authorization")

		// Check if the token is empty
		if token != "" {
			// Split the token to get the type of token
			typeOfToken := strings.Split(token, " ")

			// Check if the type of token is Bearer
			if typeOfToken[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, types.Response{
					Status:     false,
					StatusCode: http.StatusUnauthorized,
					Message:    "Unauthorized",
					Data: map[string]any{
						"error": "Check the Authorization header",
					},
				})
				c.Abort()
				return
			}
			// Verify the token
			jwtToken, err := security.VerifyToken(typeOfToken[1])

			if err != nil {
				InvalidToken(c)
				return
			}
			// Check if the token is valid
			if jwtToken.Valid {
				c.Next()
			} else {
				InvalidToken(c)
				return
			}
		}

		// If the token is empty, check if the cookie is present
		accessToken, _ := c.Request.Cookie("access_token")
		refreshToken, _ := c.Request.Cookie("refresh_token")

		// Check if the cookie is present
		if accessToken == nil && refreshToken == nil {
			c.JSON(http.StatusUnauthorized, types.Response{
				Status:     false,
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
				Data: map[string]any{
					"error": "You are not logged in",
				},
			})
			c.Abort()
			return
		}

		// Check if the token is revoked
		revoked := security.TokenRevoked(accessToken.Value, refreshToken.Value, c, true)
		if revoked {
			return
		}

		// Check if the token is expired
		if !security.IsTokenExpired(accessToken.Value) {
			c.Next()
		}

		// If the access token is expired but the refresh token is not expired, generate a new access token and set it as a cookie
		if security.IsTokenExpired(accessToken.Value) && !security.IsTokenExpired(refreshToken.Value) {
			newAccessToken := security.GenerateToken(accessToken.Value, security.AcessTokenExpireTime)
			c.SetCookie("access_token", newAccessToken, 3600, "/", "localhost", false, true)
			c.Next()
		}

		// If both the access token and refresh token are expired, return an error
		if security.IsTokenExpired(accessToken.Value) && security.IsTokenExpired(refreshToken.Value) {
			c.JSON(http.StatusUnauthorized, types.Response{
				Status:     false,
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
				Data: map[string]any{
					"error": "Access Token and Refresh Token has been Expired",
				},
			})
			c.Abort()
			return
		}

	}

}

// The InvalidToken function returns a JSON response indicating that the token is invalid and aborts
// the current request.
func InvalidToken(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, types.Response{
		Status:     false,
		StatusCode: http.StatusUnauthorized,
		Message:    "Unauthorized",
		Data: map[string]any{
			"error": "Invalid Token",
		},
	})
	c.Abort()
}
