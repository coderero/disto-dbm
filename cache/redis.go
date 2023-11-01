package cache

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// `var client *redis.Client` is declaring a variable named `client` of type `*redis.Client`. This
// variable will be used to store a pointer to a Redis client object.
var client *redis.Client

func init() {
	// `godotenv.Load()` is a function from the `godotenv` package that loads environment variables from a
	// `.env` file into the current environment. It allows you to store sensitive information like Redis
	// host, port, and password in a separate file instead of hardcoding them in your code.
	godotenv.Load()
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})
}

// The function returns a Redis client.
func getClient() *redis.Client {
	return client
}
