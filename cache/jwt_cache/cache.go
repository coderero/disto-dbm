package jwtcache

import (
	"context"

	"coderero.dev/projects/go/gin/hello/cache"
)

// RevokedToken function is used to revoke a token.

var jwt_cache = cache.GetClient()

func RevokedToken(token string) {
	jwt_cache.LPush(context.Background(), "revoked_tokens", token)
}

func IsTokenRevoked(token string) bool {
	revoked_tokens, error := jwt_cache.LRem(context.Background(), "revoked_tokens", 1, token).Result()
	if error != nil {
		return false
	}

	return revoked_tokens == 1
}
