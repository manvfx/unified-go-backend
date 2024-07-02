package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Define a struct to hold rate limiters for different users
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// AddLimiter adds a rate limiter for a given key
func (r *RateLimiter) AddLimiter(key string, limiter *rate.Limiter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.limiters[key] = limiter
}

// GetLimiter gets the rate limiter for a given key
func (r *RateLimiter) GetLimiter(key string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if limiter, exists := r.limiters[key]; exists {
		return limiter
	}
	limiter := rate.NewLimiter(1, 5) // 1 request per second, burst size of 5
	r.limiters[key] = limiter
	return limiter
}

// RateLimiterMiddleware creates a new rate limiter middleware
func RateLimiterMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use client IP as the key
		key := c.ClientIP()

		limiter := rl.GetLimiter(key)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
