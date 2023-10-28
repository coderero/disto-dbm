package controller

import (
	"net/http"

	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

type AppController struct{}

func (*AppController) Home(c *gin.Context) {
	c.JSON(http.StatusOK, types.Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Hello World!",
		Data:       nil,
	})
}
