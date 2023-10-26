package middleware

import (
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// access_token, err := c.Request.Cookie("access_token")
		// refresh_token, err1 := c.Request.Cookie("refresh_token")

		c.Next()
	}

}
