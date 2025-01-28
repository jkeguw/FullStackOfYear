package middleware

import (
	"FullStackOfYear/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, errors.NewAppError(errors.Unauthorized, "未授权访问"))
			return
		}
		// TODO: JWT验证逻辑
		c.Next()
	}
}
