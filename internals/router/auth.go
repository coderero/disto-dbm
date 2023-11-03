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
		group.POST("/register", auth.Register)
		group.POST("/login", auth.Login)
		group.POST("/logout", auth.Logout)
		group.POST("/refresh", auth.RefreshToken)
		group.GET("/isloggedin", auth.IsLoggedIn)
	}
}
