package auth

import (
	"github.com/gin-gonic/gin"
	"project/backend/internal/errors"
)

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// DummyRefreshTokenHandler 处理令牌刷新请求
func DummyRefreshTokenHandler(c *gin.Context) {
	// 这是一个空的实现，实际上应该处理令牌刷新请求
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "Token refresh is not implemented yet"))
}