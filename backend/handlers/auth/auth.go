package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/backend/internal/errors"
	"project/backend/services/auth"
	authtypes "project/backend/types/auth"
	"time"
)

// Register handles user registration
func Register(c *gin.Context) {
	var regReq authtypes.RegisterRequest
	if err := c.ShouldBindJSON(&regReq); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "Invalid request format"))
		return
	}

	// 获取认证服务
	authService, exists := c.MustGet("authService").(auth.Service)
	if !exists {
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, "Auth service not available"))
		return
	}

	// 设置请求上下文和超时
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// 注册用户
	user, err := authService.Register(ctx, &regReq)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, "Registration failed"))
		}
		return
	}

	// 返回注册响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Registration successful",
		"data":    user,
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var loginReq authtypes.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "Invalid request format"))
		return
	}

	// 设置登录类型（默认为邮箱登录）
	loginReq.LoginType = authtypes.EmailLogin

	// 获取客户端 IP
	loginReq.IP = c.ClientIP()

	// 获取认证服务
	authService, exists := c.MustGet("authService").(auth.Service)
	if !exists {
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, "Auth service not available"))
		return
	}

	// 设置请求上下文和超时
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// 调用登录服务
	resp, err := authService.Login(ctx, &loginReq)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, "Login failed"))
		}
		return
	}

	// 返回登录响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Login successful",
		"data":    resp,
	})
}

// RefreshToken handles token refresh requests
func RefreshToken(c *gin.Context) {
	var refreshReq struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&refreshReq); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "Invalid request format"))
		return
	}

	// 获取认证服务
	authService, exists := c.MustGet("authService").(auth.Service)
	if !exists {
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, "Auth service not available"))
		return
	}

	// 设置请求上下文和超时
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// 刷新令牌
	accessToken, refreshToken, err := authService.RefreshToken(ctx, refreshReq.RefreshToken)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, "Token refresh failed"))
		}
		return
	}

	// 返回新的令牌
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Token refreshed successfully",
		"data": gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"expiresIn":    3600,
			"tokenType":    "Bearer",
		},
	})
}

// Logout handles user logout
func Logout(c *gin.Context) {
	var logoutReq struct {
		DeviceID string `json:"deviceId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&logoutReq); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "Invalid request format"))
		return
	}

	// 获取认证服务
	authService, exists := c.MustGet("authService").(auth.Service)
	if !exists {
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, "Auth service not available"))
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewAppError(errors.Unauthorized, "User not authenticated"))
		return
	}

	// 设置请求上下文和超时
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// 撤销令牌
	err := authService.RevokeTokens(ctx, userID.(string), logoutReq.DeviceID)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, "Logout failed"))
		}
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Logged out successfully",
		"data":    nil,
	})
}
