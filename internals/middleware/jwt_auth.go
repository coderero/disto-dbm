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
					Status: types.Status{
						Code: http.StatusUnauthorized,
						Msg:  "Check the Authorization header",
					},
				})
				c.Abort()
				return
			}
			// Verify the token
			if security.TokenRevoked(typeOfToken[1], "", c, false) {
				InvalidToken(c)
				return
			}

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
		var accessToken string
		var refreshToken string

		raw_accessToken, _ := c.Request.Cookie("__t")
		raw_refreshToken, _ := c.Request.Cookie("__rt")
		if raw_accessToken != nil {
			accessToken = raw_accessToken.Value
		}
		if raw_refreshToken != nil {
			refreshToken = raw_refreshToken.Value
		}

		// Check if the cookie is present
		if accessToken == "" && refreshToken == "" {
			c.JSON(http.StatusUnauthorized, types.Response{
				Status: types.Status{
					Code: http.StatusUnauthorized,
					Msg:  "Unauthorized",
				},
			})
			c.Abort()
			return
		}

		// Check if the token is revoked
		haveErr := checkTokenRevoketion(accessToken, refreshToken, c)
		if haveErr {
			c.JSON(http.StatusUnauthorized, types.Response{
				Status: types.Status{
					Code: http.StatusUnauthorized,
					Msg:  "Unauthorized",
				},
			})
			c.Abort()
			return
		}

		// Check if the token is expired
		if !security.IsTokenExpired(accessToken) {
			c.Next()
		}

		// If the access token is expired but the refresh token is not expired, generate a new access token and set it as a cookie
		if security.IsTokenExpired(accessToken) && !security.IsTokenExpired(refreshToken) {
			newAccessToken := security.GenerateToken(accessToken, security.AcessTokenExpireTime)
			c.SetCookie("__t", newAccessToken, 3600, "/", "localhost", false, true)
			c.Next()
		}

		// If both the access token and refresh token are expired, return an error
		if security.IsTokenExpired(accessToken) && security.IsTokenExpired(refreshToken) {
			c.JSON(http.StatusUnauthorized, types.Response{
				Status: types.Status{
					Code: http.StatusUnauthorized,
					Msg:  "Unauthorized",
				},
			})
			c.Abort()
			return
		}

	}

}

// The function checks if a token has been revoked based on the provided access token and refresh
// token.
func checkTokenRevoketion(accessToken string, refreshToken string, c *gin.Context) bool {
	if accessToken != "" && refreshToken != "" {
		revoked := security.TokenRevoked(accessToken, refreshToken, c, true)
		if revoked {
			return true
		}
	}
	if refreshToken != "" && accessToken == "" {
		revoked := security.TokenRevoked("", refreshToken, c, false)
		if revoked {
			return true
		}
	}
	return false
}

// The InvalidToken function returns a JSON response indicating that the token is invalid and aborts
// the current request.
func InvalidToken(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, types.Response{
		Status: types.Status{
			Code: http.StatusUnauthorized,
			Msg:  "Invalid Token",
		},
	})
	c.Abort()
}
