package middleware

import (
	"FullStackOfYear/backend/internal/database"
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/services/jwt"
	"FullStackOfYear/backend/services/token"
	"github.com/gin-gonic/gin"
	"strings"
)

// Auth validates the JWT token and adds claims to context
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

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, err)
			return
		}

		// Verify token exists in Redis
		tokenManager := token.NewManager(database.RedisClient)
		exists, err := tokenManager.CheckTokenExists(c, claims.UserID, claims.DeviceID, "access")
		if err != nil {
			c.AbortWithStatusJSON(500, errors.NewAppError(errors.InternalError, "Token verification failed"))
			return
		}

		if !exists {
			c.AbortWithStatusJSON(401, errors.NewAppError(errors.Unauthorized, "Token has been revoked"))
			return
		}

		// Set claims to context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("deviceId", claims.DeviceID)

		c.Next()
	}
}

// RequireRoles validates user roles
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
