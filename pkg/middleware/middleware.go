package middleware

import (
	"time"

	"gin-demo/pkg/logger"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests and responses
func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Log request
		log.Debug("HTTP Request: " + c.Request.Method + " " + c.Request.RequestURI)

		// Process request
		c.Next()

		// Log response
		duration := time.Since(startTime)
		log.Info("HTTP Response: " + c.Request.Method + " " + c.Request.RequestURI +
			" | Status: " + string(rune(c.Writer.Status())) +
			" | Duration: " + duration.String())
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(log *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		log.Error("Panic recovered", nil)
		c.JSON(500, gin.H{
			"status":  500,
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Internal server error",
		})
	})
}

// RequestIDMiddleware adds request ID to context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("X-Request-ID")
		if requestID == "" {
			requestID = c.Request.Header.Get("X-Request-ID")
		}
		if requestID == "" {
			requestID = time.Now().Format("20060102150405") + "-" + c.ClientIP()
		}
		c.Set("request-id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
