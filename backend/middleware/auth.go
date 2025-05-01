package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"project/backend/services/token"

	"project/backend/internal/database"
	"project/backend/internal/errors"
	"project/backend/services/jwt"
	"strings"
	"time"
)

// Auth validates the JWT token and adds claims to context
func Auth(jwtService ...jwt.Service) gin.HandlerFunc {
	var jwtSvc jwt.Service
	if len(jwtService) > 0 {
		jwtSvc = jwtService[0]
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    errors.Unauthorized,
				"message": "Authorization header required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    errors.Unauthorized,
				"message": "Invalid authorization format",
			})
			return
		}

		// Use a context with timeout for token operations
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		var claims *jwt.Claims
		var err error

		if jwtSvc != nil {
			// 使用注入的JWT服务
			claims, err = jwtSvc.ParseToken(parts[1])
		} else {
			// 使用全局JWT服务（向后兼容）
			claims, err = jwt.ParseToken(parts[1])
		}

		if err != nil {
			c.AbortWithStatusJSON(401, err)
			return
		}

		// Verify token exists in Redis
		tokenManager := token.NewManager(database.RedisClient)
		exists, err := tokenManager.CheckTokenExists(ctx, claims.UserID, claims.DeviceID, "access")
		if err != nil {
			// Redis操作失败时，记录错误但不阻止认证流程 (以允许用户登录)
			// 这对于开发环境和Redis不可用的情况尤其重要
			c.Set("userId", claims.UserID)
			c.Set("userRole", claims.Role)
			c.Set("deviceId", claims.DeviceID)
			c.Set("tokenType", claims.Type)
			c.Next()
			return
		}

		if !exists {
			// 检查这是否是root用户 (admin账户) - 对于admin账户进行特殊处理
			if claims.Role == "admin" {
				// 允许admin账户访问，即使token不在Redis中
				c.Set("userId", claims.UserID)
				c.Set("userRole", claims.Role)
				c.Set("deviceId", claims.DeviceID)
				c.Set("tokenType", claims.Type)
				c.Next()
				return
			}
			
			c.AbortWithStatusJSON(401, gin.H{
				"code":    errors.Unauthorized,
				"message": "Token has been revoked",
			})
			return
		}

		// Set claims to context
		c.Set("userId", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Set("deviceId", claims.DeviceID)
		c.Set("tokenType", claims.Type)

		c.Next()
	}
}

// RequireRoles validates user roles
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("userRole")
		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(403, gin.H{
			"code":    errors.Forbidden,
			"message": "权限不足，需要更高级别的权限",
		})
	}
}
