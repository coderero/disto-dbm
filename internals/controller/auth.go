package controller

import (
	"net/http"

	"coderero.dev/projects/go/gin/hello/models"
	"coderero.dev/projects/go/gin/hello/pkg/security"
	"coderero.dev/projects/go/gin/hello/pkg/utils"
	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golodash/galidator"
)

type AuthController struct{}

// The GetUsers function returns a JSON response with a success message and an empty array of details.
func (*AuthController) SignUp(c *gin.Context) {
	// Create a new Register struct
	var signup types.SignUp

	if utils.CheckContentType(c, types.Application_x_www_form) {
		return
	}

	gal := galidator.New()
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
			Data: []any{
				customizer.DecryptErrors(err),
			}})
		return
	}

	hashPassword, err := security.HashPassword(signup.Password)
	if err != nil {
		panic(err)
	}
	user := &models.User{
		Username:  signup.Username,
		Email:     signup.Email,
		Password:  hashPassword,
		FirstName: signup.FirstName,
		LastName:  signup.LastName,
		Age:       signup.Age,
	}

	obj, err := user.Create()
	if err != nil {
		c.JSON(http.StatusConflict, types.Response{
			Status:     false,
			StatusCode: http.StatusConflict,
			Message:    "User already exists",
			Data:       []any{},
		})
		return
	}

	c.JSON(http.StatusCreated, types.Response{
		Status:     true,
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
		Data: []any{map[string]any{
			"user": obj,
		}},
	})

}

func loginValidation(c *gin.Context, register types.Login) bool {
	if len(c.Request.PostForm) > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Extra fields in form data"})
		return true
	}

	if (register.Username == "" || register.Password == "") || (register.Email != "" && register.Password != "") {
		c.JSON(http.StatusUnprocessableEntity, types.Response{
			Status:     false,
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Either Fields is required",
			Data: []any{
				"username",
				"email",
			},
		})
		return true
	}
	return false
}
