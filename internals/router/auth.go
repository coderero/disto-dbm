package router

import (
	"time"

	"coderero.dev/projects/go/gin/hello/internals/controller"
	"coderero.dev/projects/go/gin/hello/internals/handler"
	"coderero.dev/projects/go/gin/hello/internals/middleware"
	"github.com/gin-gonic/gin"
)

// The AuthRouter function sets up a Gin router with various middleware and routes for handling
// authentication-related requests.
func AuthRouter() *gin.Engine {
	router := gin.Default()

	router.HandleMethodNotAllowed = true

	router.Use(gin.CustomRecovery(handler.InternalServerErrorHandler))

	router.NoMethod(handler.NoMethodHandler())

	router.NoRoute(handler.NoRouteHandler())

	router.Use(middleware.RateLimitHandler(1000, time.Minute))

	router.Use(middleware.RateLimitHandler(100, time.Second))

	api := router.Group("/api")
	api.GET("/", controller.GetUsers)

	return router
}
