package utils

import (
	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

func CheckContentType(ctx *gin.Context, t string) bool {
	if ctx.Request.Header.Get("Content-Type") != t {
		ctx.JSON(422, types.Response{
			Code:    422,
			Status:  "Unprocessable Entity",
			Message: "The request body must be of type 'application/x-www-form-urlencoded'",
			Details: []any{},
		})
		return true
	}
	return false
}
