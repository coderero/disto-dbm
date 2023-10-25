package router

import (
	"coderero.dev/projects/go/gin/hello/internals/controller"
	"github.com/gin-gonic/gin"
)

func authRouter(group *gin.RouterGroup) {

	auth := new(controller.AuthController)

	{
		group.POST("/signup", auth.SignUp)
		group.POST("/login", auth.Login)
	}
}
