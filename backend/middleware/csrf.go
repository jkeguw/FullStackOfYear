package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		c.Header("X-CSRF-Token", token)
		c.Next()
	}
}
