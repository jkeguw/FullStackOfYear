package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/backend/internal/errors"
	userSvc "project/backend/services/user"
)

type Handler struct {
	service userSvc.Service
}

func NewHandler(service userSvc.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetUserProfile 获取用户个人资料
func (h *Handler) GetUserProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		errMsg := "用户未认证"
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errors.Unauthorized,
			"message": errMsg,
			"data":    nil,
		})
		return
	}

	result, err := h.service.GetUserProfile(c.Request.Context(), userID.(string))
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		statusCode := errors.HTTPStatusFromError(err)
		errorCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if ok {
			errorCode = appErr.Code
			errorMsg = appErr.Message
		}

		c.JSON(statusCode, gin.H{
			"code":    errorCode,
			"message": errorMsg,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    result,
	})
}

// UpdateUserProfile 更新用户个人资料
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		errMsg := "用户未认证"
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errors.Unauthorized,
			"message": errMsg,
			"data":    nil,
		})
		return
	}

	// 获取当前用户资料
	user, err := h.service.GetUserProfile(c.Request.Context(), userID.(string))
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		statusCode := errors.HTTPStatusFromError(err)
		errorCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if ok {
			errorCode = appErr.Code
			errorMsg = appErr.Message
		}

		c.JSON(statusCode, gin.H{
			"code":    errorCode,
			"message": errorMsg,
			"data":    nil,
		})
		return
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(user); err != nil {
		errMsg := "无效的请求: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.BadRequest,
			"message": errMsg,
			"data":    nil,
		})
		return
	}

	// 更新用户资料
	err = h.service.UpdateUserProfile(c.Request.Context(), userID.(string), user)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		statusCode := errors.HTTPStatusFromError(err)
		errorCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if ok {
			errorCode = appErr.Code
			errorMsg = appErr.Message
		}

		c.JSON(statusCode, gin.H{
			"code":    errorCode,
			"message": errorMsg,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    user,
	})
}
