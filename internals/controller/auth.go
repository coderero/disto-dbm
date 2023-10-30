package controller

import (
	"net/http"
	"strings"
	"time"

	jwtcache "coderero.dev/projects/go/gin/hello/cache/jwt_cache"
	"coderero.dev/projects/go/gin/hello/models"
	"coderero.dev/projects/go/gin/hello/pkg/security"
	"coderero.dev/projects/go/gin/hello/pkg/utils"
	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golodash/galidator"
)

type AuthController struct{}

var gal = galidator.New()

// The `SignUp` function is a method of the `AuthController` struct. It handles the registration
// process for a user.
func (*AuthController) Register(c *gin.Context) {
	// The code snippet is handling the registration process for a user.
	var signup types.SignUp

	if utils.CheckContentType(c, types.Application_x_www_form) {
		return
	}
	customizer := gal.Validator(signup, galidator.Messages{
		"required": "$field is required",
		"email":    "$field must be a valid email address",
		"min":      "$field is of wrong length or too short",
	})

	// Bind the form data to the Register struct
	if err := c.ShouldBindWith(&signup, binding.Form); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:     false,
			StatusCode: http.StatusBadRequest,
			Message:    "Fields are required",
			Data: map[string]any{
				"error": customizer.DecryptErrors(err),
			}})
		return
	}
	err3 := models.User{}.CheckForUser(signup.Username, signup.Email)
	if err3 {
		c.JSON(http.StatusConflict, types.Response{
			Status:     false,
			StatusCode: http.StatusConflict,
			Message:    "User already exists",
			Data:       map[string]any{},
		})
		return
	}

	hashPassword, err1 := security.HashPassword(signup.Password)
	if err1 != nil {
		panic(err1)
	}
	user := &models.User{
		Username:  signup.Username,
		Email:     signup.Email,
		Password:  hashPassword,
		FirstName: signup.FirstName,
		LastName:  signup.LastName,
		Age:       signup.Age,
	}

	registeredObj := user.Create()

	// The code snippet is generating access and refresh tokens for the registered user and setting them as
	// cookies in the response. It then returns a JSON response with the status, status code, message, and
	// the generated access and refresh tokens. This is typically done after a successful registration
	// process to provide the user with authentication tokens for subsequent requests.
	accessToken, refreshToken := security.GenerateAuthTokens(registeredObj)

	c.SetCookie("access_token", accessToken, 300, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 86400, "/", "localhost", false, true)

	c.JSON(http.StatusCreated, types.Response{
		Status:     true,
		StatusCode: http.StatusCreated,
		Message:    "Registration Successful",
		Data: map[string]any{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// The `Signin` function is a method of the `AuthController` struct. It handles the login process for a
// user.
func (*AuthController) Login(c *gin.Context) {
	// The code snippet is handling the login process for a user.
	var login types.Login
	if utils.CheckContentType(c, types.Application_x_www_form) {
		return
	}

	// The code snippet is retrieving the access token and refresh token from the request header and
	// cookies.
	token := c.Request.Header.Get("Authorization")
	raw_accessToken, _ := c.Request.Cookie("access_token")
	raw_refreshToken, _ := c.Request.Cookie("refresh_token")

	// The `revokeTokenIfPresent` function is used to check if an access token and refresh token are
	// present in the request. It takes in the access token, raw access token cookie, raw refresh token
	// cookie, and the current Gin context as parameters.
	revokeTokenIfPresent(token, raw_accessToken, raw_refreshToken, c)

	gal := galidator.New()
	customizer := gal.Validator(login, galidator.Messages{
		"required": "$field is required",
		"email":    "$field must be a valid email address",
		"min":      "$field is of wrong length or too short",
	})

	if err := c.ShouldBindWith(&login, binding.Form); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:     false,
			StatusCode: http.StatusBadRequest,
			Message:    "Fields are required",
			Data: map[string]any{
				"error": customizer.DecryptErrors(err),
			}})
		return
	}

	if loginValidation(c, login) {
		return
	}

	var user *models.User

	registeredObj := user.GetUserForLogin(login.Username, login.Email)

	// Check for password
	if !security.ComparePassword(login.Password, registeredObj.Password) {
		c.JSON(http.StatusUnauthorized, types.Response{
			Status:     false,
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid Credentials",
			Data:       map[string]any{},
		})
		return
	}

	if registeredObj.ID == 0 {
		c.JSON(http.StatusNotFound, types.Response{
			Status:     false,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
			Data:       map[string]any{},
		})
		return
	}

	// This code snippet is generating access and refresh tokens for a registered user and setting them as
	// cookies in the response. It then returns a JSON response with the status, status code, message, and
	// the generated access and refresh tokens. This is typically done after a successful login process to
	// provide the user with authentication tokens for subsequent requests.
	accessToken, refreshToken := security.GenerateAuthTokens(registeredObj)

	c.SetCookie("access_token", accessToken, 300, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 86400, "/", "localhost", false, true)

	c.JSON(http.StatusOK, types.Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Login Successful",
		Data: map[string]any{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// The `Logout` function is a method of the `AuthController` struct. It handles the logout process for
// a user.
func (*AuthController) Logout(c *gin.Context) {
	func() {
		// The code snippet is checking for the presence of an access token in the request header. It
		// retrieves the access token from the "Authorization" header by splitting the header value and
		// taking the second part.
		var accessToken, refreshToken string
		token := c.Request.Header.Get("Authorization")
		accessToken = strings.Split(token, " ")[1]
		if accessToken == "" {
			if accessToken == "" || refreshToken == "" {
				c.JSON(http.StatusBadRequest, types.Response{
					Status:     false,
					StatusCode: http.StatusBadRequest,
					Message:    "Refresh Token and Access Token is required",
					Data:       map[string]any{},
				})
				return
			}

			revoked := security.TokenRevoked(accessToken, refreshToken, c, false)
			if revoked {
				return
			}

			revoke(accessToken, refreshToken)

			c.JSON(http.StatusOK, types.Response{
				Status:     true,
				StatusCode: http.StatusOK,
				Message:    "Logout Successful",
				Data:       map[string]any{},
			})
			return
		}
	}()

	// This code snippet is handling the logout process for a user. It first tries to retrieve the access
	// token and refresh token from the request cookies. If both cookies are not found, it returns a JSON
	// response with a status code of 400 and a message indicating that no cookies were found.
	access, err := c.Request.Cookie("access_token")
	refresh, err1 := c.Request.Cookie("refresh_token")
	if err != nil && err1 != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:     false,
			StatusCode: http.StatusBadRequest,
			Message:    "No Cookie found",
			Data:       map[string]any{},
		})
		return
	}

	accessToken := access.Value
	refreshToken := refresh.Value

	revoked := security.TokenRevoked(accessToken, "", c, false)
	if revoked {
		return
	}

	revoke(accessToken, refreshToken)

	c.JSON(http.StatusOK, types.Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Logout Successful",
		Data:       map[string]any{},
	})
}

// The `RefreshToken` function is a method of the `AuthController` struct. It handles the process of
// refreshing an access token using a refresh token.
func (*AuthController) RefreshToken(c *gin.Context) {
	var tokens *types.RefreshToken

	customizer := gal.Validator(tokens, galidator.Messages{
		"required": "$field is required",
	})

	if err := c.ShouldBindJSON(&tokens); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:     false,
			StatusCode: http.StatusBadRequest,
			Message:    "Fields are required",
			Data: map[string]any{
				"error": customizer.DecryptErrors(err),
			}})
		return
	}

	refreshToken, accessToken := tokens.RefreshToken, tokens.AccessToken

	revoked := security.TokenRevoked(accessToken, refreshToken, c, true)
	if revoked {
		return
	}

	jwtcache.RevokeToken(accessToken)

	if security.IsTokenExpired(refreshToken) {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:     false,
			StatusCode: http.StatusBadRequest,
			Message:    "Unauthorized",
			Data: map[string]any{
				"error": "Refresh Token Expired",
			},
		})
		return
	}

	claims, err := security.VerifyToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, types.Response{
			Status:     false,
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
			Data:       map[string]any{},
		})
		return
	}
	sub, _ := claims.Claims.GetSubject()
	accessToken = security.GenerateToken(sub, time.Now().Add(time.Minute*5))

	c.JSON(http.StatusOK, types.Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Token Refreshed",
		Data: map[string]any{
			"access_token": accessToken,
		},
	})
}

func revokeTokenIfPresent(token string, raw_accessToken, raw_refreshToken *http.Cookie, c *gin.Context) {
	var accessToken, refreshToken string
	if token != "" {
		accessToken = strings.Split(token, " ")[1]
	}
	if raw_accessToken != nil {
		accessToken = raw_accessToken.Value
	}
	if raw_refreshToken != nil {
		refreshToken = raw_refreshToken.Value
	}
	if accessToken != "" && refreshToken != "" {
		revoked := (jwtcache.IsTokenRevoked(accessToken) && jwtcache.IsTokenRevoked(refreshToken)) || (security.IsTokenExpired(accessToken) && security.IsTokenExpired(refreshToken))
		if revoked {
			return
		}
	}
	if accessToken != "" {
		if jwtcache.IsTokenRevoked(accessToken) || security.IsTokenExpired(accessToken) {
			return
		}
		jwtcache.RevokeToken(accessToken)

	}
	if refreshToken != "" {
		if jwtcache.IsTokenRevoked(refreshToken) || security.IsTokenExpired(refreshToken) {
			return
		}
		jwtcache.RevokeToken(refreshToken)
	}

}

// The function "revoke" revokes the access token and refresh token by adding them to the revoked token
// cache.
func revoke(accessToken string, refreshToken string) {
	jwtcache.RevokeToken(accessToken)
	jwtcache.RevokeToken(refreshToken)
}

// The loginValidation function checks if the required fields for login are provided and returns true
// if there are any errors.
func loginValidation(c *gin.Context, register types.Login) bool {
	if register.Email != "" && register.Username != "" {
		c.JSON(http.StatusUnprocessableEntity, types.Response{
			Status:     false,
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Either Fields is required",
			Data: map[string]any{
				"error": "either username or email is required",
			},
		})
		return true
	}
	if len(c.Request.PostForm) > 3 {
		c.JSON(http.StatusUnprocessableEntity, types.Response{
			Status:     false,
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Only Username or Email and Password is required",
			Data:       map[string]any{},
		})
		return true
	}

	return false
}
