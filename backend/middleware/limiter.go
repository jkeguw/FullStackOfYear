package middleware

import (
	"project/backend/internal/errors"
	"project/backend/services/limiter"
	"github.com/gin-gonic/gin"
)

// RateLimit applies rate limiting to routes
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		allowed, err := limiter.CheckRateLimit(c, clientIP)
		if err != nil {
			c.AbortWithStatusJSON(500, errors.NewAppError(errors.InternalError, "Rate limit check failed"))
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(429, errors.NewAppError(errors.TooManyRequests, "Too many requests"))
			return
		}

		c.Next()
	}
}
