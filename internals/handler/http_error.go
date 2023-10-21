package handler

import (
	"net/http"

	types "coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

// The function NoMethodHandler returns a Gin middleware function that handles requests with a 405
// status code and a JSON response indicating that the requested method is not allowed.
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, types.Response{
			Code:    http.StatusMethodNotAllowed,
			Status:  "Method Not Allowed",
			Message: "The requested method is not allowed",
			Details: []any{},
		})
		c.Next()
	}

}

// The function NoRouteHandler returns a Gin middleware that handles 404 errors by returning a JSON
// response with an error message.
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, types.Response{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "The requested resource was not found",
			Details: []any{},
		})
		c.Next()
	}
}

// The function InternalServerErrorHandler handles internal server errors by returning a JSON response
// with a status code of 500 and an error message.
func InternalServerErrorHandler(c *gin.Context, _ any) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, types.Response{
		Code:    http.StatusInternalServerError,
		Status:  "Internal Server Error",
		Message: "The server encountered an internal error",
		Details: []any{
			c.Errors,
		},
	})
	c.Next()
}
