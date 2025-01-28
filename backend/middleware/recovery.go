package middleware

import (
	"FullStackOfYear/backend/config"
	"FullStackOfYear/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(nil,
		func(c *gin.Context, err interface{}) {
			config.Logger.Sugar().Errorw("panic recovered",
				"error", err,
				"path", c.Request.URL.Path,
			)
			c.JSON(500, errors.NewAppError(errors.InternalError, "服务器内部错误"))
		})
}
