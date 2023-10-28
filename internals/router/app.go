package router

import (
	"coderero.dev/projects/go/gin/hello/internals/controller"
	"coderero.dev/projects/go/gin/hello/internals/middleware"
	"github.com/gin-gonic/gin"
)

func appRouter(group *gin.RouterGroup) {
	group.Use(middleware.JWTAuthMiddleWare())
	app := new(controller.AppController)

	{
		group.GET("", app.Home)
	}
}
