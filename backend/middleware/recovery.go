package middleware

import (
	"github.com/gin-gonic/gin"
	"project/backend/config"
	"project/backend/internal/errors"
)

func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(nil,
		func(c *gin.Context, err interface{}) {
			config.Logger.Sugar().Errorw("panic recovered",
				"error", err,
				"path", c.Request.URL.Path,
			)
			// 使用统一的响应格式
			c.JSON(500, gin.H{
				"code":    errors.InternalError,
				"message": "服务器内部错误",
			})
		})
}
