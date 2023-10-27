package utils

import (
	"net/http"

	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

// The function checks if the content type of the request is the expected type and returns true if it
// is not.
func CheckContentType(ctx *gin.Context, t string) bool {
	if ctx.Request.Header.Get("Content-Type") != t {
		ctx.JSON(422, types.Response{
			Status:     false,
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "The request body must be of type 'application/x-www-form-urlencoded'",
			Data:       nil,
		})
		return true
	}
	return false
}
