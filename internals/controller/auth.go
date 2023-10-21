package controller

import (
	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

// The GetUsers function returns a JSON response with a success message and an empty array of details.
func (auth *AuthController) GetUsers(c *gin.Context) {
	c.JSON(200, types.Response{
		Code:    200,
		Status:  "OK",
		Message: "Successfully retrieved users",
		Details: []any{},
	})
	c.Next()
}
