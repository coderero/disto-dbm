package middleware

import (
	"net/http"
	"os"

	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func parseMiddleware(middleware func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// After passing through the Gorilla middleware, we reach this point.
			// Since we don't want to write anything or terminate the request,
			c.Request = r
		}))

		h.ServeHTTP(c.Writer, c.Request)

		if c.Writer.Status() > 399 || c.Writer.Status() == http.StatusForbidden {
			c.AbortWithStatusJSON(http.StatusForbidden, types.Response{
				Status: types.Status{
					Code: http.StatusForbidden,
					Msg:  "Forbidden",
				},
			})
			return
		}

		c.Next() // Proceed to the next middleware or request handler.
	}
}

func CsrfCheck() gin.HandlerFunc {

	csrfMiddleware := csrf.Protect(
		[]byte(os.Getenv("CSRF_SECRET")),
		csrf.Secure(true),
		csrf.CookieName("__csrf"),
		csrf.Path("/"),
		csrf.HttpOnly(false),
		csrf.MaxAge(3600*12),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		})),
	)

	return parseMiddleware(csrfMiddleware)
}
