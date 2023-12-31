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
			Status: types.Status{
				Code: http.StatusMethodNotAllowed,
				Msg:  "method not allowed",
			},
		})
		c.Next()
	}

}

// The function NoRouteHandler returns a Gin middleware that handles 404 errors by returning a JSON
// response with an error message.
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, types.Response{
			Status: types.Status{
				Code: http.StatusNotFound,
				Msg:  "not found",
			},
		})
		c.Next()
	}
}

// The function InternalServerErrorHandler handles internal server errors by returning a JSON response
// with a status code of 500 and an error message.
func InternalServerErrorHandler(c *gin.Context, _ any) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, types.Response{
		Status: types.Status{
			Code: http.StatusInternalServerError,
			Msg:  "internal server error",
		},
	})
	c.Next()
}
