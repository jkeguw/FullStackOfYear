package middleware

import (
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/services/auth"
	"github.com/gin-gonic/gin"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, errors.NewAppError(errors.Unauthorized, "Authorization header required"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(401, errors.NewAppError(errors.Unauthorized, "Invalid authorization format"))
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, err)
			return
		}

		// Store user information in context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("device", claims.Device)

		c.Next()
	}
}

// RequireRoles Role Verification Middleware
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(403, errors.NewAppError(errors.Forbidden, "Insufficient permissions"))
	}
}
