package router

import (
	"coderero.dev/projects/go/gin/hello/internals/controller"
	"github.com/gin-gonic/gin"
)

// The function `csrfRouter` sets up a route for generating a CSRF token.
func csrfRouter(group *gin.RouterGroup) {
	// `csrf := new(controller.CSRFController)` is creating a new instance of the `CSRFController` struct.
	csrf := new(controller.CSRFController)

	// The following code block registers the CSRF route.
	{
		group.GET("/csrf", csrf.GenerateCsrfToken)
	}
}
