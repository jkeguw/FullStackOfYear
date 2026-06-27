package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that restricts cross-origin requests to allowed origins.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := getAllowedOrigins()

		allowed := false
		for _, o := range allowedOrigins {
			if strings.TrimSpace(o) == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-CSRF-Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func getAllowedOrigins() []string {
	env := os.Getenv("ALLOWED_ORIGINS")
	if env != "" {
		return strings.Split(env, ",")
	}
	return []string{
		"http://localhost",
		"http://localhost:80",
		"http://localhost:3000",
		"http://localhost:5173",
	}
}
