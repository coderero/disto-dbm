package jwtcache

import (
	"context"
	"fmt"

	"coderero.dev/projects/go/gin/hello/cache"
)

// The line `var jwt_cache = cache.GetClient()` is initializing a variable `jwt_cache` with the value
// returned by the `GetClient()` function from the `cache` package. This line is likely setting up a
// connection to a cache server or creating a cache client object that will be used for storing and
// retrieving data from the cache.
var jwt_cache = cache.GetClient()

// The function RevokedToken adds a token to a list of revoked tokens in a cache.
func RevokeToken(token string) {
	err := jwt_cache.RPush(context.Background(), "revoked_tokens", token).Err()
	if err != nil {
		fmt.Println(err)
	}
}

// The IsTokenRevoked function checks if a given token is revoked by querying a cache.
func IsTokenRevoked(token string) bool {
	revoked_tokens, error := jwt_cache.LRem(context.Background(), "revoked_tokens", 1, token).Result()
	if error != nil {
		return false
	}

	return revoked_tokens == 1
}
