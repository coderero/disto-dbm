package controller

import (
	"net/http"

	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

type AppController struct{}

// The `func (*AppController) Home(c *gin.Context)` function is a method of the `AppController` struct.
// It is used as a handler function for the `/` route.
func (*AppController) Home(c *gin.Context) {
	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "ok",
		},
		Data: []map[string]any{
			{"message": "Hello, World!"},
		}})
}
