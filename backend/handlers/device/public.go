package device

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"project/backend/internal/errors"
)

// ListPublicUserDevices 获取公开的用户设备配置
func (h *Handler) ListPublicUserDevices(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 调用服务
	result, err := h.deviceService.GetPublicUserDevices(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), gin.H{"error": err.Error()})
		return
	}

	// 响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    result,
	})
}