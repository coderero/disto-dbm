package controller

import (
	"net/http"

	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

type CSRFController struct{}

func (CSRFController) GenerateCsrfToken(c *gin.Context) {
	token := csrf.Token(c.Request)
	c.Header("X-CSRF-Token", token)
	c.JSON(http.StatusOK, types.Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "CSRF Token Generated",
		Data: map[string]interface{}{
			"csrf_token": token,
		}})
}
