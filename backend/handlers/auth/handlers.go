package auth

import (
	"github.com/gin-gonic/gin"
	"project/backend/internal/errors"
)

// Handlers contains all auth related handlers
type Handlers struct {
	OAuth *OAuthHandler
	// TODO: 其他 handlers
}

// DummyRegister 处理用户注册
func DummyRegister(c *gin.Context) {
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "This endpoint is not fully implemented"))
}

// DummyLogin 处理用户登录
func DummyLogin(c *gin.Context) {
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "This endpoint is not fully implemented"))
}

// DummyRefreshToken 处理令牌刷新
func DummyRefreshToken(c *gin.Context) {
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "This endpoint is not fully implemented"))
}

// DummyLogout 处理用户登出
func DummyLogout(c *gin.Context) {
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "This endpoint is not fully implemented"))
}