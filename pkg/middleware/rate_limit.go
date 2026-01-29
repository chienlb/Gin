package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"gin-demo/pkg/resilience"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitConfig defines rate limit configuration
type RateLimitConfig struct {
	RequestsPerSecond int
	Burst             int
}

// RateLimiterMiddleware creates a rate limiting middleware
func RateLimiterMiddleware(config RateLimitConfig) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(config.RequestsPerSecond), config.Burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimiter implements per-IP rate limiting
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewIPRateLimiter creates a new IP-based rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// GetLimiter returns the rate limiter for a given IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.limiters[ip] = limiter
	}

	return limiter
}

// IPRateLimitMiddleware creates a per-IP rate limiting middleware
func IPRateLimitMiddleware(requestsPerSecond int, burst int) gin.HandlerFunc {
	ipLimiter := NewIPRateLimiter(rate.Limit(requestsPerSecond), burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := ipLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Too many requests from your IP, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// TimeoutMiddleware adds timeout to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				c.JSON(http.StatusGatewayTimeout, gin.H{
					"status":  "error",
					"code":    "TIMEOUT",
					"message": "Request timeout",
				})
				c.Abort()
			}
		}
	}
}

// CircuitBreakerMiddleware adds circuit breaker protection
func CircuitBreakerMiddleware(maxFailures uint, resetTimeout time.Duration) gin.HandlerFunc {
	cb := resilience.NewCircuitBreaker(maxFailures, resetTimeout)

	return func(c *gin.Context) {
		err := cb.Execute(func() error {
			c.Next()

			// Check if request failed
			if c.Writer.Status() >= 500 {
				return http.ErrAbortHandler
			}
			return nil
		})

		if err != nil {
			if err.Error() == "circuit breaker is open" {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status":  "error",
					"code":    "SERVICE_UNAVAILABLE",
					"message": "Service temporarily unavailable",
				})
				c.Abort()
			}
		}
	}
}
