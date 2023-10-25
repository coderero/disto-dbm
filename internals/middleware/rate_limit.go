package middleware

import (
	"net/http"
	"sync"
	"time"

	types "coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

// The RateLimitHandler function is a middleware that limits the number of requests from a specific IP
// address within a given time duration.
func RateLimitHandler(limit int, duration time.Duration) gin.HandlerFunc {
	currentIp := make(map[string]int)

	var mutex sync.Mutex = sync.Mutex{}

	return func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()

		ip := c.ClientIP()
		count, ok := currentIp[ip]
		if !ok {
			currentIp[ip] = 1
		} else {
			if count >= limit {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, types.Response{
					Status:     false,
					StatusCode: http.StatusTooManyRequests,
					Message:    "You have exceeded your request limit",
					Data:       map[string]any{},
				})
				c.Next()
				return
			}
			currentIp[ip] = count + 1
		}

		go func() {
			<-time.After(duration)
			mutex.Lock()
			defer mutex.Unlock()
			delete(currentIp, ip)
		}()
		c.Next()
	}

}
