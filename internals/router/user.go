package router

import (
	"coderero.dev/projects/go/gin/hello/internals/controller"
	"coderero.dev/projects/go/gin/hello/internals/middleware"
	"github.com/gin-gonic/gin"
)

func userRouter(group *gin.RouterGroup) {
	// `auth := new(controller.AuthController)` is creating a new instance of the `AuthController` struct
	// from the `controller` package. This instance is assigned to the variable `auth`.
	// `user := new(controller.UserController)` is creating a new instance of the `UserController` struct
	// from the `controller` package. This instance is assigned to the variable `user`.
	user := new(controller.UserController)

	group.Use(middleware.JWTAuthMiddleWare())

	// The following code block registers auth routes.
	{
		group.PATCH("/user", user.UpdateUser)
		group.DELETE("/user", user.DeleteUser)
	}
}
