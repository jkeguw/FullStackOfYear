package middleware

import (
	"github.com/gin-gonic/gin"
	"project/backend/config"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		config.Logger.Sugar().Infow("request completed",
			"path", path,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", time.Since(start),
		)
	}
}
