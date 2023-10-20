package controller

import (
	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

// The GetUsers function returns a JSON response with a success message and an empty array of details.
func GetUsers(c *gin.Context) {
	c.JSON(200, types.SuccessResponse{Success: types.ResponseSkeleton{
		Code:    200,
		Status:  "OK",
		Message: "Successfully fetched users",
		Details: []any{},
	}})
}
