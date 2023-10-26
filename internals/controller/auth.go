package controller

import (
	"net/http"
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

// The GetUsers function returns a JSON response with a success message and an empty array of details.
func (*AuthController) SignUp(c *gin.Context) {
	// Create a new Register struct
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

	haveQueryParam := getTokenInbody(registeredObj, c)
	if haveQueryParam {
		return
	}

	setAuthCookies(registeredObj, c)

	c.JSON(http.StatusCreated, types.Response{
		Status:     true,
		StatusCode: http.StatusCreated,
		Message:    "Registration Successful",
		Data:       map[string]any{},
	})
}

func (*AuthController) Login(c *gin.Context) {
	var login types.Login
	if utils.CheckContentType(c, types.Application_x_www_form) {
		return
	}

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

	registeredObj := user.GetUserForLogin(login.Username, login.Email, login.Password)

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

	haveQueryParam := getTokenInbody(registeredObj, c)
	if haveQueryParam {
		return
	}
	accessToken, refreshToken := generateAuthTokes(registeredObj)

	c.SetCookie("access_token", accessToken, 300, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 86400, "/", "localhost", false, true)

	c.JSON(http.StatusOK, types.Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Login Successful",
		Data:       map[string]any{},
	})
}

func (*AuthController) Logout(c *gin.Context) {
	var refreshToken string
	var accessToken string
	_, token := c.GetQuery("token")
	if token {
		accessToken = c.Request.Header.Get("Authorization")
		refreshToken = c.Query("refresh_token")

		if accessToken == "" && refreshToken == "" {
			c.JSON(http.StatusBadRequest, types.Response{
				Status:     false,
				StatusCode: http.StatusBadRequest,
				Message:    "Refresh Token and Access Token is required",
				Data:       map[string]any{},
			})
			return
		}

		revoked := tokenRevoked(accessToken, refreshToken, c)
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

	accessToken = access.Value
	refreshToken = refresh.Value

	revoked := tokenRevoked(accessToken, refreshToken, c)
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

	revoked := tokenRevoked(accessToken, refreshToken, c)
	if revoked {
		return
	}

	jwtcache.RevokedToken(accessToken)

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

func revoke(accessToken string, refreshToken string) {
	jwtcache.RevokedToken(accessToken)
	jwtcache.RevokedToken(refreshToken)
}

func setAuthCookies(registeredObj *models.User, c *gin.Context) {
	accessToken, refreshToken := generateAuthTokes(registeredObj)

	c.SetCookie("access_token", accessToken, 300, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 86400, "/", "localhost", false, true)
}

func generateAuthTokes(obj *models.User) (string, string) {
	sub := obj.Email
	currentTime := time.Now()
	accessToken := security.GenerateToken(sub, currentTime.Add(time.Minute*5))
	refreshToken := security.GenerateToken(sub, currentTime.Add(time.Hour*24))
	return accessToken, refreshToken
}

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

func getTokenInbody(registeredObj *models.User, c *gin.Context) bool {
	access, refresh := generateAuthTokes(registeredObj)
	_, present := c.GetQuery("token")
	if present {
		c.JSON(http.StatusCreated, types.Response{
			Status:     true,
			StatusCode: http.StatusCreated,
			Message:    "Registration Successful",
			Data: map[string]any{
				"access_token":  access,
				"refresh_token": refresh,
			}})
		return true
	}
	return false
}

func tokenRevoked(accessToken string, refreshToken string, c *gin.Context) bool {
	if jwtcache.IsTokenRevoked(accessToken) && jwtcache.IsTokenRevoked(refreshToken) {
		c.JSON(http.StatusUnauthorized, types.Response{
			Status:     false,
			StatusCode: http.StatusUnauthorized,
			Message:    "Token's Have been Revoked",
			Data:       map[string]any{},
		})
		return true
	}
	return false
}
