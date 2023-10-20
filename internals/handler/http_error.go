package handler

import (
	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

// The function NoMethodHandler returns a Gin middleware function that handles requests with a 405
// status code and a JSON response indicating that the requested method is not allowed.
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(405, types.ErrorResponse{Error: types.ResponseSkeleton{
			Code:    405,
			Status:  "Method Not Allowed",
			Message: "The requested method is not allowed",
			Details: []any{},
		}})
		c.Next()
	}

}

// The function NoRouteHandler returns a Gin middleware that handles 404 errors by returning a JSON
// response with an error message.
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, types.ErrorResponse{Error: types.ResponseSkeleton{
			Code:    404,
			Status:  "Not Found",
			Message: "The requested resource was not found",
			Details: []any{},
		}})
		c.Next()
	}
}

// The function InternalServerErrorHandler handles internal server errors by returning a JSON response
// with a status code of 500 and an error message.
func InternalServerErrorHandler(c *gin.Context, _ any) {
	c.AbortWithStatusJSON(500, types.ErrorResponse{Error: types.ResponseSkeleton{
		Code:    500,
		Status:  "Internal Server Error",
		Message: "The server encountered an internal error",
		Details: []any{},
	}})
	c.Next()
}
