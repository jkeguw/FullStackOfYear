package auth

import (
	"github.com/gin-gonic/gin"
	"project/backend/internal/errors"
)

// OAuthHandler 处理OAuth认证
type OAuthHandler struct{}

// HandleOAuthLogin 处理OAuth登录请求
func (h *OAuthHandler) HandleOAuthLogin(c *gin.Context) {
	// 暂不实现
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "OAuth login not implemented"))
}

// HandleOAuthCallback 处理OAuth回调
func (h *OAuthHandler) HandleOAuthCallback(c *gin.Context) {
	// 暂不实现
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "OAuth callback not implemented"))
}