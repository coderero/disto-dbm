package router

import (
	"coderero.dev/projects/go/gin/hello/internals/controller"
	"github.com/gin-gonic/gin"
)

func authRouter(group *gin.RouterGroup) {

	// `auth := new(controller.AuthController)` is creating a new instance of the `AuthController` struct
	// from the `controller` package. This instance is assigned to the variable `auth`.
	auth := new(controller.AuthController)

	{
		group.POST("/signup", auth.SignUp)
		group.POST("/signin", auth.Signin)
		group.GET("/logout", auth.Logout)
		group.POST("/refresh", auth.RefreshToken)
	}
}
