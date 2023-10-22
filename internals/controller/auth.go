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
func (*AuthController) Register(c *gin.Context) {
	// Create a new Register struct
	var register types.Register

	if utils.CheckContentType(c, types.Application_x_www_form) {
		return
	}

	gal := galidator.New()
	customizer := gal.Validator(register, galidator.Messages{
		"required": "$field is required",
		"email":    "$field must be a valid email address",
		"min":      "$field is of wrong length or too short",
	})

	// Bind the form data to the Register struct
	if err := c.ShouldBindWith(&register, binding.Form); err != nil {

		c.JSON(http.StatusBadRequest, types.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Fields are required",
			Details: []any{
				customizer.DecryptErrors(err),
			}})
		return
	}

	hashPassword, err := security.HashPassword(register.Password)
	if err != nil {
		panic(err)
	}
	user := &models.User{
		Username:  register.Username,
		Email:     register.Email,
		Password:  hashPassword,
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Age:       register.Age,
	}

	obj := user.Create()

	c.JSON(http.StatusCreated, types.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "User created successfully",
		Details: []any{
			obj,
		},
	})

}

func loginValidation(c *gin.Context, register types.Login) bool {
	if len(c.Request.PostForm) > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Extra fields in form data"})
		return true
	}

	if (register.Username == "" || register.Password == "") || (register.Email != "" && register.Password != "") {
		c.JSON(http.StatusUnprocessableEntity, types.Response{
			Code:    http.StatusUnprocessableEntity,
			Status:  "Unprocessable Entity",
			Message: "Either Fields is required",
			Details: []any{
				"username",
				"email",
			},
		})
		return true
	}
	return false
}
