package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nirmal/crm/utils"
	"golang.org/x/time/rate"
)

// Limiter is a struct that holds the rate limiter and its associated mutex.
var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

// GetVisitor retrieves the rate limiter for a given IP address.
func GetVisitor(ip string) *rate.Limiter {
	mu.Lock()
    defer mu.Unlock()

    if limiter, exists := visitors[ip]; exists {
        return limiter
    }

    limiter := rate.NewLimiter(utils.RATE_LIMIT, utils.BURST_LIMIT) // 1 request per second, with a burst size of 5. Can be adjusted as needed.
    visitors[ip] = limiter
    return limiter
}

// RateLimiterMiddleware is a Gin middleware to limit the number of requests.
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := GetVisitor(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
