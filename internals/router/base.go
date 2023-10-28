package router

import (
	"time"

	"coderero.dev/projects/go/gin/hello/internals/handler"
	"coderero.dev/projects/go/gin/hello/internals/middleware"
	"github.com/gin-gonic/gin"
)

// The function `Router` returns a Gin router with a sub-router for handling authentication routes.
func Router() *gin.Engine {
	r := gin.Default()
	r.HandleMethodNotAllowed = true

	// Error Handlers
	r.Use(gin.CustomRecovery(handler.InternalServerErrorHandler))
	r.NoMethod(handler.NoMethodHandler())
	r.NoRoute(handler.NoRouteHandler())

	// Middleware's
	r.Use(middleware.RateLimitHandler(1000, time.Minute))
	r.Use(middleware.RateLimitHandler(100, time.Second))

	// Sub-Routers
	sub := r.Group("/api/v1")

	// Route Handlers
	authRouter(sub)
	appRouter(sub)

	return r
}
