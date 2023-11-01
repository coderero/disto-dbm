package middleware

import (
	"net/http"
	"sync"
	"time"

	types "coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

// The RateLimitHandler function is a middleware that limits the number of requests from a specific IP
// address within a given duration.
func RateLimitHandler(limit int, duration time.Duration) gin.HandlerFunc {
	// Map to store the number of requests from a specific IP address
	currentIp := make(map[string]int)

	// Mutex to lock the map
	var mutex sync.Mutex = sync.Mutex{}

	// Return the middleware function
	return func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()

		// Get the IP address of the client
		ip := c.ClientIP()

		// Get the number of requests from this IP address
		count, ok := currentIp[ip]
		if !ok {
			// If the IP address is not in the map, add it with a count of 1
			currentIp[ip] = 1
		} else {
			// If the IP address is in the map, exceed the limit
			if count >= limit {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, types.Response{
					Status:     false,
					StatusCode: http.StatusTooManyRequests,
					Message:    "You have exceeded your request limit",
				})
				c.Next()
				return
			}
			currentIp[ip] = count + 1
		}

		// Delete the IP address from the map after the duration
		go func() {
			<-time.After(duration)
			mutex.Lock()
			defer mutex.Unlock()
			delete(currentIp, ip)
		}()
		c.Next()
	}

}
