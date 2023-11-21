package router

import (
	"coderero.dev/projects/go/gin/hello/internals/controller"
	"coderero.dev/projects/go/gin/hello/internals/middleware"
	"github.com/gin-gonic/gin"
)

// The function appRouter is used to register routes for the app group.
func userRouter(group *gin.RouterGroup) {
	// The `group.Use(middleware.JWTAuthMiddleWare())` function is used to register the `JWTAuthMiddleWare`.
	group.Use(middleware.JWTAuthMiddleWare())
	user := new(controller.UserController)

	// The following code block registers app routes.
	{
		group.GET("/user", user.Get)
		group.PATCH("/user", user.Update)
		group.DELETE("/user", user.Delete)
	}
}
