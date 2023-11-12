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
	"github.com/golodash/galidator"
)

type AuthController struct{}

// The `gal` variable is an instance of the `galidator` package. It is used to validate the form data
// provided by the user during registration.
var gal = galidator.New()

// The `SignUp` function is a method of the `AuthController` struct. It handles the registration
// process for a user.
func (*AuthController) Register(c *gin.Context) {
	// The code snippet is handling the registration process for a user.
	var signup types.SignUp

	// The `CheckContentType` function is used to check if the content type of the request is
	if utils.CheckContentType(c, types.Application_json) {
		return
	}

	// The `tokenActionType` variable is used to check if the user wants to return the access token and
	tokenActionType := c.Query("return_token")

	// The `previousTokens` function is used to check if the access token and refresh token are present in
	// the request header or cookies. If they are present, they are revoked.
	previousTokens(c)

	// The `Validator` function is used to validate the form data. It takes in the form data and a custom
	// error message as parameters.
	customizer := gal.Validator(signup, galidator.Messages{
		"required": "$field is required",
		"email":    "$field must be a valid email address",
		"min":      "$field is of wrong length or too short",
	})

	// Bind the form data to the Register struct and check for errors in the process.
	if err := c.ShouldBindJSON(&signup); err != nil {
		if strings.Contains(err.Error(), "Key:") {

			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "validation error",
				},
				Data: map[string]any{
					"errors": customizer.DecryptErrors(err),
				},
			})
			return
		}
		if strings.Contains(err.Error(), "ReadBool:") || strings.Contains(err.Error(), "readUint64:") || strings.Contains(err.Error(), "ReadString:") {

			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "validation error",
				},
				Data: map[string]any{
					"error": "you have invalid fields in your json",
				},
			})
			return
		}
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "you have provided an invalid json",
			},
			Data: map[string]any{
				"error": err.Error()},
		})
		return
	}

	// The `loginValidation` function is used to check if the required fields for login are provided and
	// returns true if there are any errors.
	err3 := models.User{}.CheckForUser(signup.Username, signup.Email)
	if err3 {
		c.JSON(http.StatusConflict, types.Response{
			Status: types.Status{
				Code: http.StatusConflict,
				Msg:  "user already exists",
			},
		})
		return
	}

	// The `HashPassword` function is used to hash the password provided by the user during registration.
	hashPassword, err1 := security.HashPassword(signup.Password)
	if err1 != nil {
		panic(err1)
	}

	// The `Create` function is used to create a new user. It takes in the user object as a parameter.
	user := &models.User{
		Username:  signup.Username,
		Email:     signup.Email,
		Password:  hashPassword,
		FirstName: signup.FirstName,
		LastName:  signup.LastName,
		Age:       signup.Age,
	}

	// The `Create` function is used to create a new user. It takes in the user object as a parameter.
	registeredObj := user.Create()

	// The code snippet is generating access and refresh tokens for the registered user and setting them as
	// cookies in the response. It then returns a JSON response with the status, status code, message, and
	// the generated access and refresh tokens. This is typically done after a successful registration
	// process to provide the user with authentication tokens for subsequent requests.
	accessToken, refreshToken := security.GenerateAuthTokens(registeredObj)
	if tokenActionType == "true" {
		c.JSON(http.StatusCreated, types.Response{
			Status: types.Status{
				Code: http.StatusCreated,
				Msg:  "registration successful",
			},
			Data: map[string]any{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		})
		return
	}

	// The below code is setting two cookies named "__t" and "__rt" with the values of "accessToken" and
	// "refreshToken" respectively. The cookies are set to expire after a certain duration (300 seconds for
	// access token and 86400 seconds for refresh token).
	setTokenInCookies(c, accessToken, refreshToken)

	// The code snippet is returning a JSON response with the status, status code, and message. This is
	// typically done after a successful registration process.
	c.JSON(http.StatusCreated, types.Response{
		Status: types.Status{
			Code: http.StatusCreated,
			Msg:  "registration successful",
		},
	})
}

// The `Login` function is a method of the `AuthController` struct. It handles the login process for a
// user.
func (*AuthController) Login(c *gin.Context) {
	// The code snippet is handling the login process for a user.
	var login types.Login

	// The `CheckContentType` function is used to check if the content type of the request is
	if utils.CheckContentType(c, types.Application_json) {
		return
	}

	// The `tokenActionType` variable is used to check if the user wants to return the access token and
	// refresh token in the response body or as cookies.
	tokenActionType := c.Query("return_token")

	// The `previousTokens` function is used to check if the access token and refresh token are present in
	// the request header or cookies. If they are present, they are revoked.
	previousTokens(c)

	// The `Validator` function is used to validate the form data. It takes in the form data and a custom
	// error message as parameters.
	customizer := gal.Validator(login, galidator.Messages{
		"required": "$field is required",
		"email":    "$field must be a valid email address",
		"min":      "$field is of wrong length or too short",
	})

	// Bind the form data to the Login struct and check for errors in the process.
	if err := c.ShouldBindJSON(&login); err != nil {
		if strings.Contains(err.Error(), "Key:") {

			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "validation error",
				},
				Data: map[string]any{
					"errors": customizer.DecryptErrors(err),
				},
			})
			return
		}
		if strings.Contains(err.Error(), "ReadBool:") || strings.Contains(err.Error(), "readUint64:") || strings.Contains(err.Error(), "ReadString:") {

			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "validation error",
				},
				Data: map[string]any{
					"error": "you have invalid fields in your json",
				},
			})
			return
		}
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "you have provided an invalid json",
			}})
		return
	}

	// The `loginValidation` function is used to check if the required fields for login are provided and
	// returns true if there are any errors.
	if loginValidation(c, login) {
		return
	}

	// Initialize a new user object to store the user object returned by the `GetUserForLogin` function.
	var user *models.User

	// The `GetUserForLogin` function is used to get the user object for the user trying to login. It takes
	// in the username or email of the user as parameters.
	registeredObj := user.GetUserForLogin(login.Username, login.Email)
	if registeredObj.ID == 0 {
		c.JSON(http.StatusNotFound, types.Response{
			Status: types.Status{
				Code: http.StatusNotFound,
				Msg:  "user not found",
			},
		})
		return
	}

	// The `ComparePassword` function is used to compare the password provided by the user during login
	// with the hashed password stored in the database.
	if !security.ComparePassword(login.Password, registeredObj.Password) {
		c.JSON(http.StatusUnauthorized, types.Response{
			Status: types.Status{
				Code: http.StatusUnauthorized,
				Msg:  "invalid password",
			},
		})
		return
	}

	// This code snippet is generating access and refresh tokens for a registered user and setting them as
	// cookies in the response. It then returns a JSON response with the status, status code, message, and
	// the generated access and refresh tokens. This is typically done after a successful login process to
	// provide the user with authentication tokens for subsequent requests.
	accessToken, refreshToken := security.GenerateAuthTokens(registeredObj)

	// The code snippet is checking if the user wants to return the access token and refresh token in the
	// response body or as cookies. If the user wants to return the tokens in the response body, the code
	// snippet returns a JSON response with the status, status code, message, and the generated access and
	// refresh tokens.
	if tokenActionType == "true" {
		c.JSON(http.StatusOK, types.Response{
			Status: types.Status{
				Code: http.StatusOK,
				Msg:  "login successful",
			},
			Data: map[string]any{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		})
		return
	}

	// The below code is setting two cookies named "__t" and "__rt" with the values of "accessToken" and
	// "refreshToken" respectively. The cookies are set to expire after a certain duration (300 seconds for
	// access token and 86400 seconds for refresh token).
	setTokenInCookies(c, accessToken, refreshToken)

	// The code snippet is returning a JSON response with the status, status code, and message. This is
	// typically done after a successful login process.
	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "login successful",
		},
	})

}

// The `Logout` function is a method of the `AuthController` struct. It handles the logout process for
// a user.
func (*AuthController) Logout(c *gin.Context) {
	// The `token` variable is used to get the access token from the request header.
	token := c.Request.Header.Get("Authorization")

	// The `raw_accessToken` and `raw_refreshToken` variables are used to get the access token and refresh token
	// respectively from the request cookies.
	raw_accessToken, _ := c.Request.Cookie("__t")
	raw_refreshToken, _ := c.Request.Cookie("__rt")

	// Although the revokeTokenFunction below doing great job but add a extra validation check is good for
	// more information for the user to avoid panicking or 500 error.
	if token == "" && raw_accessToken == nil && raw_refreshToken == nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "no token provided (i.e. you are not logged in).",
			}})
		return
	}

	// The `revokeTokenIfPresent` function is used to check if an access token and refresh token are
	revokeTokenIfPresent(token, raw_accessToken, raw_refreshToken)

	// The code snippet is deleting the access token and refresh token cookies from the response.
	c.SetCookie("__t", "", -1, "/", "localhost", true, true)
	c.SetCookie("__rt", "", -1, "/", "localhost", true, true)

	// The code snippet is returning a JSON response with the status, status code, and message. This is
	// typically done after a successful logout process.
	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "logout Successful",
		},
	})

}

// The `RefreshToken` function is a method of the `AuthController` struct. It handles the process of
// refreshing an access token using a refresh token.
func (*AuthController) RefreshToken(c *gin.Context) {
	// The `tokens` variable is Initialize a new `RefreshToken` struct to store the refresh token and access
	// token provided by the user.
	var tokens *types.RefreshToken

	// The `CheckContentType` function is used to check if the content type of the request is
	if utils.CheckContentType(c, types.Application_json) {
		return
	}

	// The `Validator` function is used to validate the form data. It takes in the form data and a custom
	// error message as parameters.
	customizer := gal.Validator(tokens, galidator.Messages{
		"required": "$field is required",
	})

	// Bind the form data to the RefreshToken struct and check for errors in the process.
	if err := c.ShouldBindJSON(&tokens); err != nil {
		if strings.Contains(err.Error(), "Key:") {

			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "validation error",
				},
				Data: map[string]any{
					"errors": customizer.DecryptErrors(err),
				},
			})
			return
		}
		if strings.Contains(err.Error(), "ReadBool:") || strings.Contains(err.Error(), "readUint64:") || strings.Contains(err.Error(), "ReadString:") {

			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "validation error",
				},
				Data: map[string]any{
					"error": "you have invalid fields in your json",
				},
			})
			return
		}
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "you have provided an invalid json",
			},
		})
		return
	}

	// The `refreshToken` and `accessToken` variables are used to get the refresh token and access token
	// respectively from the request body.
	refreshToken, accessToken := tokens.RefreshToken, tokens.AcessToken

	// The `tokenRevoked` variable is used to check if the access token and refresh token are revoked.
	revoked := security.IsTokenRevoked(accessToken, refreshToken, c, true)
	if revoked {
		return
	}

	// The `RevokeToken` function is used to revoke the access token.
	cache.RevokeToken(accessToken)

	// The `IsTokenExpired` function is used to check if the refresh token is expired.
	if security.IsTokenExpired(refreshToken) {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "token expired",
			},
		})
		return
	}

	// The `VerifyToken` function is used to verify the refresh token and return the claims.
	claims, err := security.VerifyToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, types.Response{
			Status: types.Status{
				Code: http.StatusUnauthorized,
				Msg:  "unauthorized",
			},
		})
		return
	}

	// The `sub` variable is used to get the subject from the claims.
	sub, _ := claims.Claims.GetSubject()

	// The `GenerateToken` function is used to generate a new access token for the user.
	accessToken = security.GenerateToken(sub, time.Now().Add(time.Minute*5))

	// The code snippet is returning the new access token in the response body.
	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "token refreshed",
		},
		Data: map[string]any{
			"access_token": accessToken,
		},
	})
}

func (*AuthController) IsLoggedIn(c *gin.Context) {
	// The `token` variable is used to get the access token from the request header.
	token := c.Request.Header.Get("Authorization")

	// The `raw_accessToken` and `raw_refreshToken` variables are used to get the access token and refresh token
	// respectively from the request cookies.
	raw_accessToken, _ := c.Request.Cookie("__t")
	raw_refreshToken, _ := c.Request.Cookie("__rt")

	// Check if the token is empty and if the access token and refresh token are nil.
	if token == "" && raw_accessToken == nil && raw_refreshToken == nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "no token provided (i.e. you are not logged in).",
			}})
		return
	}

	// The `revokedOrExpired` variable is used to check if the access token and refresh token are revoked or
	// expired.
	revokedOrExpired := (cache.IsTokenRevoked(token) && cache.IsTokenRevoked(raw_accessToken.Value) && cache.IsTokenRevoked(raw_refreshToken.Value)) || (security.IsTokenExpired(token) && security.IsTokenExpired(raw_accessToken.Value) && security.IsTokenExpired(raw_refreshToken.Value))
	if revokedOrExpired {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "token expired or revoked",
			},
		})
		return
	}

	// The code snippet is returning a JSON response with the status, status code, and message. This is
	// typically done after a successful login process.
	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "logged In",
		},
	})
}

// The `previousTokens` function is used to check if the access token and refresh token are present in
// the request header or cookies. If they are present, they are revoked.
func previousTokens(c *gin.Context) {
	// The `token` variable is used to get the access token from the request header.
	token := c.Request.Header.Get("Authorization")

	// The `raw_accessToken` and `raw_refreshToken` variables are used to get the access token and refresh token
	// respectively from the request cookies.
	raw_accessToken, _ := c.Request.Cookie("__t")
	raw_refreshToken, _ := c.Request.Cookie("__rt")

	// The `revokeTokenIfPresent` function is used to check if an access token and refresh token are
	// present in the request header or cookies. If they are present, they are revoked.
	revokeTokenIfPresent(token, raw_accessToken, raw_refreshToken)
}

func revokeTokenIfPresent(token string, raw_accessToken, raw_refreshToken *http.Cookie) {
	// The `accessToken` and `refreshToken` variables are used to get the access token and refresh token
	var accessToken, refreshToken string

	// The code below is checking the values of different variables and assigning them to other variables.
	if token != "" {
		accessToken = strings.Split(token, " ")[1]
	}
	if raw_accessToken != nil {
		accessToken = raw_accessToken.Value
	}
	if raw_refreshToken != nil {
		refreshToken = raw_refreshToken.Value
	}

	// The code below is checking if both the `accessToken` and `refreshToken` are not empty. If they are
	// not empty, it then checks if either both tokens are revoked or both tokens are expired. If either of
	// these conditions is true, the code returns and does not proceed further.
	if accessToken != "" && refreshToken != "" {
		revoked := (cache.IsTokenRevoked(accessToken) && cache.IsTokenRevoked(refreshToken)) || (security.IsTokenExpired(accessToken) && security.IsTokenExpired(refreshToken))
		if revoked {
			return
		}
	}

	// The code below is checking if the `accessToken` is not empty. If it is not empty, it then checks if
	// the token is revoked or expired using the `cache.IsTokenRevoked()` and `security.IsTokenExpired()`
	// functions respectively. If the token is either revoked or expired, the code returns. Otherwise, it
	// calls the `cache.RevokeToken()` function to revoke the token.
	if accessToken != "" {
		if cache.IsTokenRevoked(accessToken) || security.IsTokenExpired(accessToken) {
			return
		}
		cache.RevokeToken(accessToken)

	}

	// The code below is checking if the `refreshToken` is not empty. If it is not empty, it then checks if
	// the token is revoked or expired using the `cache.IsTokenRevoked()` and `security.IsTokenExpired()`
	// functions respectively. If the token is revoked or expired, the code returns without performing any
	// further actions. If the token is valid, it then revokes the token by calling the
	// `cache.RevokeToken()` function.
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
	// The code below is revoking the access token.
	cache.RevokeToken(accessToken)

	// The code below is checking if the refresh token is not empty. If it is not empty, it then revokes
	if refreshToken != "" {
		cache.RevokeToken(refreshToken)
	}
}

// The loginValidation function checks if the required fields for login are provided and returns true
// if there are any errors.
func loginValidation(c *gin.Context, login types.Login) bool {

	// The code below is checking if both the email and username fields in the login object are empty. If
	// they are empty, it returns a JSON response with a status code of 422 (Unprocessable Entity) and a
	// message indicating that either the username or email is required.
	if login.Email == "" && login.Username == "" {
		c.JSON(http.StatusUnprocessableEntity, types.Response{
			Status: types.Status{
				Code: http.StatusUnprocessableEntity,
				Msg:  "username or email is required",
			},
		})
		return true
	}

	// The code below is checking if both the email and username fields in the login object are not empty.
	// If both fields are not empty, it returns a JSON response with a status code of 422 (Unprocessable
	// Entity) and a message stating that only the username or email is required.
	if login.Email != "" && login.Username != "" {
		c.JSON(http.StatusUnprocessableEntity, types.Response{
			Status: types.Status{
				Code: http.StatusUnprocessableEntity,
				Msg:  "only username or email is required",
			},
		})
		return true
	}

	return false
}

// The `setTokenInCookies` function is used to set the access token and refresh token as cookies in the
// response.
func setTokenInCookies(c *gin.Context, accessToken string, refreshToken string) {
	c.SetCookie("__t", accessToken, 300, "/", "localhost", true, true)
	c.SetCookie("__rt", refreshToken, 86400, "/", "localhost", true, true)
}
