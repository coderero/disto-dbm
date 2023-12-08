package router

import (
	"fmt"
	"os"
	"strings"
	"time"

	"coderero.dev/projects/go/gin/hello/internals/handler"
	"coderero.dev/projects/go/gin/hello/internals/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

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
	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		MaxAge:           12 * time.Hour,
	}))

	fmt.Println(strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","))

	// Sub-Routers
	sub := r.Group("/api/v1")

	// Add Global Middlewares
	sub.Use(middleware.CsrfCheck())

	// Route Handlers
	authRouter(sub)
	csrfRouter(sub)
	appRouter(sub)
	userRouter(sub)

	return r
}
