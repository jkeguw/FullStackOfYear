package user

import (
	"project/backend/internal/errors"
	userSvc "project/backend/services/user"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusUnauthorized, errors.NewAppError(errors.Unauthorized, "用户未认证"))
		return
	}

	result, err := h.service.GetUserProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateUserProfile 更新用户个人资料
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewAppError(errors.Unauthorized, "用户未认证"))
		return
	}

	// 获取当前用户资料
	user, err := h.service.GetUserProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的请求: "+err.Error()))
		return
	}

	// 更新用户资料
	err = h.service.UpdateUserProfile(c.Request.Context(), userID.(string), user)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	c.JSON(http.StatusOK, user)
}