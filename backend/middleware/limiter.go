package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"sync"
)

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	var clients sync.Map

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if _, ok := clients.Load(ip); !ok {
			clients.Store(ip, rate.NewLimiter(r, b))
		}
		if !limiter.Allow() {
			c.AbortWithStatus(429)
			return
		}
		c.Next()
	}
}
