package middleware

import (
	"net/http"
	"os"

	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

// The `init` function loads environment variables from a `.env` file.
func init() {
	godotenv.Load()
}

// The function `parseCSRFMiddleware` is a helper function that wraps a given middleware function and
// handles CSRF protection for a Gin framework application.
func parseCSRFMiddleware(middleware func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// After passing through the Gorilla middleware, we reach this point.
			// Since we don't want to write anything or terminate the request,
			c.Request = r
		}))

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		h.ServeHTTP(c.Writer, c.Request)

		// Check if the status code is 400 or greater, or if the status code is 403.
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

// The function `CsrfCheck` returns a Gin middleware function that adds CSRF protection to the
// application.
func CsrfCheck() gin.HandlerFunc {
	// The `csrfMiddleware` variable is a middleware function that adds CSRF protection to the
	// application using the Gorilla CSRF middleware.
	csrfMiddleware := csrf.Protect(
		// The `[]byte(os.Getenv("CSRF_SECRET"))` is the secret key used to generate the CSRF token.
		[]byte(os.Getenv("CSRF_SECRET")),

		// The `csrf.Secure(true)` option ensures that the CSRF cookie is only sent over HTTPS.
		csrf.Secure(true),

		// The `csrf.CookieName("__csrf")` option sets the name of the CSRF cookie to "__csrf".
		csrf.CookieName("__csrf"),

		// The `csrf.Path("/")` option sets the path of the CSRF cookie to "/".
		csrf.Path("/"),

		// The `csrf.HttpOnly(true)` option sets the HttpOnly flag on the CSRF cookie.
		csrf.HttpOnly(true),

		// The `csrf.MaxAge(3600*12)` option sets the maximum age of the CSRF cookie to 12 hours.
		csrf.MaxAge(3600*12),

		// The `csrf.ErrorHandler()` option sets the error handler for the CSRF middleware.
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		})),
	)

	// The `parseCSRFMiddleware` function is called with the `csrfMiddleware` variable as an argument.
	return parseCSRFMiddleware(csrfMiddleware)
}
