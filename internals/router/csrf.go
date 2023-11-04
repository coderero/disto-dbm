package router

import (
	"coderero.dev/projects/go/gin/hello/internals/controller"
	"github.com/gin-gonic/gin"
)

func csrfRouter(group *gin.RouterGroup) {
	csrf := new(controller.CSRFController)

	{
		group.GET("/csrf", csrf.GenerateCsrfToken)
	}
}
