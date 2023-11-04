package controller

import (
	"net/http"
	"strings"
	"time"

	"coderero.dev/projects/go/gin/hello/cache"
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

	tokenActionType := c.Query("return_token")

	// The `previousTokens` function is used to check if the access token and refresh token are present in
	// the request header or cookies. If they are present, they are revoked.
	previousTokens(c)

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
	if tokenActionType == "true" {
		c.JSON(http.StatusCreated, types.Response{
			Status:     true,
			StatusCode: http.StatusCreated,
			Message:    "Login Successful",
			Data: map[string]any{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		})
		return
	}
	c.SetCookie("__t", accessToken, 300, "/", "localhost", true, true)
	c.SetCookie("__rt", refreshToken, 86400, "/", "localhost", true, true)

	c.JSON(http.StatusCreated, types.Response{
		Status:     true,
		StatusCode: http.StatusCreated,
		Message:    "Login Successful",
	})
}

// The `Login` function is a method of the `AuthController` struct. It handles the login process for a
// user.
func (*AuthController) Login(c *gin.Context) {
	// The code snippet is handling the login process for a user.
	var login types.Login
	if utils.CheckContentType(c, types.Application_x_www_form) {
		return
	}

	tokenActionType := c.Query("return_token")

	// The `previousTokens` function is used to check if the access token and refresh token are present in
	// the request header or cookies. If they are present, they are revoked.
	previousTokens(c)

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
		})
		return
	}

	if registeredObj.ID == 0 {
		c.JSON(http.StatusNotFound, types.Response{
			Status:     false,
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		})
		return
	}

	// This code snippet is generating access and refresh tokens for a registered user and setting them as
	// cookies in the response. It then returns a JSON response with the status, status code, message, and
	// the generated access and refresh tokens. This is typically done after a successful login process to
	// provide the user with authentication tokens for subsequent requests.
	accessToken, refreshToken := security.GenerateAuthTokens(registeredObj)
	if tokenActionType == "true" {

		c.JSON(http.StatusOK, types.Response{
			Status:     true,
			StatusCode: http.StatusOK,
			Message:    "Login Successful",
			Data: map[string]any{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		})
		return
	}
	c.SetCookie("__t", accessToken, 300, "/", "localhost", true, true)
	c.SetCookie("__rt", refreshToken, 86400, "/", "localhost", true, true)

	c.JSON(http.StatusOK, types.Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Login Successful",
	})

}

// The `Logout` function is a method of the `AuthController` struct. It handles the logout process for
// a user.
func (*AuthController) Logout(c *gin.Context) {
	// The code snippet is handling the logout process for a user.
	token := c.Request.Header.Get("Authorization")
	raw_accessToken, _ := c.Request.Cookie("__t")
	raw_refreshToken, _ := c.Request.Cookie("__rt")

	// Although the revokeTokenFunction below doing great job but add a extra validation check is good for
	// more information for the user to avoid panicking or 500 error.
	if token == "" && raw_accessToken == nil && raw_refreshToken == nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:     false,
			StatusCode: http.StatusBadRequest,
			Message:    "Bad Request",
			Data: map[string]any{
				"error": "No token provided (i.e. you are not logged in).",
			}})
		return
	}

	// The `revokeTokenIfPresent` function is used to check if an access token and refresh token are
	revokeTokenIfPresent(token, raw_accessToken, raw_refreshToken)

	// The `revoke` function is used to revoke the access token and refresh token. It takes in the access
	// token and refresh token as parameters.

	// The code snippet is deleting the access token and refresh token cookies from the response.
	c.SetCookie("__t", "", -1, "/", "localhost", true, true)
	c.SetCookie("__rt", "", -1, "/", "localhost", true, true)

	c.JSON(http.StatusOK, types.Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Logout Successful",
	})

}

// The `RefreshToken` function is a method of the `AuthController` struct. It handles the process of
// refreshing an access token using a refresh token.
func (*AuthController) RefreshToken(c *gin.Context) {
	var tokens *types.RefreshToken

	if utils.CheckContentType(c, types.Application_json) {
		return
	}

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

	refreshToken, accessToken := tokens.RefreshToken, tokens.AcessToken

	revoked := security.TokenRevoked(accessToken, refreshToken, c, true)
	if revoked {
		return
	}

	cache.RevokeToken(accessToken)

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

func (*AuthController) IsLoggedIn(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	raw_accessToken, _ := c.Request.Cookie("__t")
	raw_refreshToken, _ := c.Request.Cookie("__rt")

	if token == "" && raw_accessToken == nil && raw_refreshToken == nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:     false,
			StatusCode: http.StatusBadRequest,
			Message:    "Bad Request",
			Data: map[string]any{
				"error": "No token provided (i.e. you are not logged in).",
			}})
		return
	}

	revokedOrExpired := (cache.IsTokenRevoked(token) && cache.IsTokenRevoked(raw_accessToken.Value) && cache.IsTokenRevoked(raw_refreshToken.Value)) || (security.IsTokenExpired(token) && security.IsTokenExpired(raw_accessToken.Value) && security.IsTokenExpired(raw_refreshToken.Value))
	if revokedOrExpired {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:     false,
			StatusCode: http.StatusBadRequest,
			Message:    "Bad Request",
			Data: map[string]any{
				"error": "Token's have been revoked or expired.",
			}})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "You are Logged in",
		Data: map[string]any{
			"details": true,
		},
	})
}

func previousTokens(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	raw_accessToken, _ := c.Request.Cookie("__t")
	raw_refreshToken, _ := c.Request.Cookie("__rt")

	revokeTokenIfPresent(token, raw_accessToken, raw_refreshToken)
}

func revokeTokenIfPresent(token string, raw_accessToken, raw_refreshToken *http.Cookie) {
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
		revoked := (cache.IsTokenRevoked(accessToken) && cache.IsTokenRevoked(refreshToken)) || (security.IsTokenExpired(accessToken) && security.IsTokenExpired(refreshToken))
		if revoked {
			return
		}
	}
	if accessToken != "" {
		if cache.IsTokenRevoked(accessToken) || security.IsTokenExpired(accessToken) {
			return
		}
		cache.RevokeToken(accessToken)

	}
	if refreshToken != "" {
		if cache.IsTokenRevoked(refreshToken) || security.IsTokenExpired(refreshToken) {
			return
		}
		cache.RevokeToken(refreshToken)
	}

}

// The function "revoke" revokes the access token and refresh token by adding them to the revoked token
// cache.
func revoke(accessToken string, refreshToken string) {
	cache.RevokeToken(accessToken)
	if refreshToken != "" {
		cache.RevokeToken(refreshToken)
	}
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
		})
		return true
	}

	return false
}
