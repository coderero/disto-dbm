package router

import (
	"coderero.dev/projects/go/gin/hello/internals/controller"
	"coderero.dev/projects/go/gin/hello/internals/middleware"
	"github.com/gin-gonic/gin"
)

// The function appRouter is used to register routes for the app group.
func appRouter(group *gin.RouterGroup) {
	// The `group.Use(middleware.JWTAuthMiddleWare())` function is used to register the `JWTAuthMiddleWare`.
	group.Use(middleware.JWTAuthMiddleWare())

	// `app := new(controller.AppController)` is creating a new instance of the `AppController` struct.
	app := new(controller.AppController)

	// The following code block registers app routes.
	{
		group.GET("/", app.Home)
	}
}
