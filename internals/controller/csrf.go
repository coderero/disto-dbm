package controller

import (
	"net/http"

	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

type CSRFController struct{}

// The `GenerateCsrfToken` function is a method of the `CSRFController` struct. It is used as a handler
// function for the `/csrf` route, which is used to generate a CSRF token for the client and send it
// back as a response header.
func (CSRFController) GenerateCsrfToken(c *gin.Context) {
	token := csrf.Token(c.Request)
	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "ok",
		},
		Data: []map[string]interface{}{
			{"csrf_token": token},
		}})
}
