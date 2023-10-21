package controller

import (
	"net/http"

	"coderero.dev/projects/go/gin/hello/pkg/utils"
	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct{}

// The GetUsers function returns a JSON response with a success message and an empty array of details.
func (*AuthController) Register(c *gin.Context) {
	if utils.CheckContentType(c, types.Application_x_www_form) {
		return
	}

	// Create a new Register struct
	var register types.Login

	//Check if request form contains values more than the struct fields

	// Bind the form data to the Register struct
	if err := c.ShouldBindWith(&register, binding.Form); err != nil {
		return
	}
}

func LoginValidation(c *gin.Context, register types.Login) bool {
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
