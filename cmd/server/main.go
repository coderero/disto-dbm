package main

import (
	"os"

	"coderero.dev/projects/go/gin/hello/internals/router"
	"coderero.dev/projects/go/gin/hello/models"
	"coderero.dev/projects/go/gin/hello/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// The `init` function checks if the "private.key" and "public.pem" files exist in the "./certs/jwt"
// directory and panics if they don't.
func init() {
	godotenv.Load()

	func() {
		files, err := os.ReadDir("./certs/jwt")
		if err != nil {
			panic(err)
		}

		if utils.ContainsFile(files, "private.key") && utils.ContainsFile(files, "public.pem") {
			return
		}
	}()
}

// The main function sets the mode for the Gin framework and starts the server on port 8000.
func main() {
	gin.SetMode(os.Getenv("MODE"))
	r := router.AuthRouter()

	var port string = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	db := &models.User{}
	db.GetUserById(1)
	r.Run(":" + port)
}
