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

		// `token := c.Request.Header.Get("Authorization")` is retrieving the value of the "Authorization"
		// header from the HTTP request. The "Authorization" header is commonly used to send authentication
		// credentials, such as a token, with the request. In this case, the code is retrieving the token
		// from the header and assigning it to the `token` variable for further processing.
		token := c.Request.Header.Get("Authorization")

		// This code block is checking if the `token` variable is not empty. If the token is not empty, it
		// splits the token to get the type of token. It then checks if the type of token is "Bearer". If the
		// type of token is not "Bearer", it returns a JSON response indicating that the Authorization header
		// is incorrect and aborts the current request.
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
			if security.IsTokenRevoked(typeOfToken[1], "", c, false) {
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

		// The code `var ( accessToken string refreshToken string )` is declaring two variables,
		// `accessToken` and `refreshToken`, of type string. These variables are used to store the values of
		// access token and refresh token, respectively.
		var (
			accessToken  string
			refreshToken string
		)

		// These lines of code are retrieving the values of two cookies named "__t" and "__rt" from the HTTP
		// request. The cookies are accessed using the `c.Request.Cookie()` method, which returns a cookie and
		// an error.The retrieved cookie values are then assigned to the variables `raw_accessToken` and
		// `raw_refreshToken` respectively.
		raw_accessToken, _ := c.Request.Cookie("__t")
		raw_refreshToken, _ := c.Request.Cookie("__rt")

		// The code block is checking if the `raw_accessToken` and `raw_refreshToken` variables are not nil.
		// If they are not nil, it means that the corresponding cookies "__t" and "__rt" exist in the HTTP
		// request. The values of these cookies are then assigned to the `accessToken` and `refreshToken`
		// variables, respectively.
		if raw_accessToken != nil {
			accessToken = raw_accessToken.Value
		}
		if raw_refreshToken != nil {
			refreshToken = raw_refreshToken.Value
		}

		// The code block is checking if both the `accessToken` and `refreshToken` variables are empty. If
		// they are empty, it means that the user is not authenticated and does not have valid tokens. In
		// this case, the code returns a JSON response with a status code of 401 (Unauthorized) and a message
		// indicating that the user is unauthorized. It then aborts the current request and returns from the
		// middleware function.
		if accessToken == "" && refreshToken == "" {
			c.JSON(http.StatusUnauthorized, types.Response{
				Status: types.Status{
					Code: http.StatusUnauthorized,
					Msg:  "unauthorized",
				},
			})
			c.Abort()
			return
		}

		// The code block is calling the `checkTokenRevoketion` function with the `accessToken`,
		// `refreshToken`, and `c` (gin.Context) as arguments. The function checks if the tokens have been
		// revoked based on the provided access token and refresh token.
		haveErr := checkTokenRevoketion(accessToken, refreshToken, c)
		if haveErr {
			c.JSON(http.StatusUnauthorized, types.Response{
				Status: types.Status{
					Code: http.StatusUnauthorized,
					Msg:  "unauthorized",
				},
			})
			c.Abort()
			return
		}

		// The code block is checking if the access token is not expired. If the access token is not expired,
		// it calls the `c.Next()` function to pass the request to the next middleware function.
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
					Msg:  "unauthorized",
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
		revoked := security.IsTokenRevoked(accessToken, refreshToken, c, true)
		if revoked {
			return true
		}
	}
	if refreshToken != "" && accessToken == "" {
		revoked := security.IsTokenRevoked("", refreshToken, c, false)
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
			Msg:  "invalid token",
		},
	})
	c.Abort()
}
