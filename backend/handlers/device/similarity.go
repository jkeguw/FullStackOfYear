package device

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"project/backend/internal/errors"
)

// CompareMice 比较多个鼠标
func (h *Handler) CompareMice(c *gin.Context) {
	// 从查询参数获取鼠标ID
	idsQuery := c.Query("ids")
	if idsQuery == "" {
		errors.HandleError(c, errors.NewBadRequestError("缺少鼠标ID参数"))
		return
	}

	// 解析鼠标ID
	mouseIDs := strings.Split(idsQuery, ",")
	if len(mouseIDs) < 1 {
		errors.HandleError(c, errors.NewBadRequestError("至少需要一个鼠标ID"))
		return
	}

	// 调用服务
	result, err := h.deviceService.CompareMice(c.Request.Context(), mouseIDs)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    result,
	})
}

// FindSimilarMice 寻找相似鼠标
func (h *Handler) FindSimilarMice(c *gin.Context) {
	// 从路径参数获取鼠标ID
	mouseID := c.Param("id")
	if mouseID == "" {
		errors.HandleError(c, errors.NewBadRequestError("缺少鼠标ID"))
		return
	}

	// 获取限制参数
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	// 调用服务
	similarMice, err := h.deviceService.FindSimilarMice(c.Request.Context(), mouseID, limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    similarMice,
	})
}