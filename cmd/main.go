package main

import (
	"os"

	"coderero.dev/projects/go/gin/hello/internals/router"
	"coderero.dev/projects/go/gin/hello/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// The line `var port string` is declaring a variable named `port` of type `string`. This variable will
// be used to store the port number on which the server will listen for incoming requests.
var port string

// The `init` function checks if the "private.key" and "public.pem" files exist in the "./certs/jwt"
// directory and panics if they don't.
func init() {
	godotenv.Load()

	// Fuction to check for rsa key pair files in ./certs
	func() {
		files, err := os.ReadDir("./certs")
		if err != nil {
			panic(err)
		}

		if utils.ContainsFile(files, "private.key") && utils.ContainsFile(files, "public.pem") {
			return
		}
	}()

	//Function to get the port and mode from the environment variables
	func() {
		// gin.SetMode(os.Getenv("GIN_MODE"))
		port = os.Getenv("PORT")
		if port == "" {
			port = "8000"
		}
	}()

}

// The main function sets the mode for the Gin framework and starts the server on port 8000.
func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))

	// The code `r := router.Router()` creates a new instance of a Gin router. The `Router()` function is a
	// custom function defined in the `router` package that returns a new Gin router.
	r := router.Router()

	// The following code starts the server on port 8000.
	r.Run(":" + port)
}
