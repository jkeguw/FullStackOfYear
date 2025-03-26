// handlers/measurement/handlers.go

package measurement

import (
	"FullStackOfYear/backend/internal/errors"
	measurementSvc "FullStackOfYear/backend/services/measurement"
	measurementTypes "FullStackOfYear/backend/types/measurement"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service measurementSvc.Service
}

func NewHandler(service measurementSvc.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateMeasurement 创建测量记录
func (h *Handler) CreateMeasurement(c *gin.Context) {
	var request measurementTypes.CreateMeasurementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("无效的请求: "+err.Error()))
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	result, err := h.service.CreateMeasurement(c.Request.Context(), userID.(string), request)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 转换为响应格式
	response := measurementTypes.MeasurementResponse{
		ID:        result.ID.Hex(),
		Palm:      result.Measurements.Palm,
		Length:    result.Measurements.Length,
		Unit:      result.Measurements.Unit,
		Quality:   &result.Quality,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetMeasurement 获取单条测量记录
func (h *Handler) GetMeasurement(c *gin.Context) {
	measurementID := c.Param("id")

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	result, err := h.service.GetMeasurement(c.Request.Context(), userID.(string), measurementID)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 转换为响应格式
	response := measurementTypes.MeasurementResponse{
		ID:        result.ID.Hex(),
		Palm:      result.Measurements.Palm,
		Length:    result.Measurements.Length,
		Unit:      result.Measurements.Unit,
		Quality:   &result.Quality,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// ListMeasurements 获取测量记录列表
func (h *Handler) ListMeasurements(c *gin.Context) {
	var request measurementTypes.MeasurementListRequest

	// 解析请求参数
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("无效的请求参数: "+err.Error()))
		return
	}

	// 设置默认值
	if request.Page == 0 {
		request.Page = 1
	}
	if request.PageSize == 0 {
		request.PageSize = 20
	}

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	result, err := h.service.ListMeasurements(c.Request.Context(), userID.(string), request)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateMeasurement 更新测量记录
func (h *Handler) UpdateMeasurement(c *gin.Context) {
	measurementID := c.Param("id")

	var request measurementTypes.UpdateMeasurementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("无效的请求: "+err.Error()))
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	result, err := h.service.UpdateMeasurement(c.Request.Context(), userID.(string), measurementID, request)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 转换为响应格式
	response := measurementTypes.MeasurementResponse{
		ID:        result.ID.Hex(),
		Palm:      result.Measurements.Palm,
		Length:    result.Measurements.Length,
		Unit:      result.Measurements.Unit,
		Quality:   &result.Quality,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteMeasurement 删除测量记录
func (h *Handler) DeleteMeasurement(c *gin.Context) {
	measurementID := c.Param("id")

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	err := h.service.DeleteMeasurement(c.Request.Context(), userID.(string), measurementID)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserStats 获取用户测量统计
func (h *Handler) GetUserStats(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	result, err := h.service.GetUserStats(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	response := measurementTypes.MeasurementStatsResponse{
		AveragePalm:      result.Averages.Palm,
		AverageLength:    result.Averages.Length,
		HandSize:         result.HandSize,
		MeasurementCount: result.MeasurementCount,
		LastMeasuredAt:   result.LastMeasuredAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetRecommendations 获取设备推荐
func (h *Handler) GetRecommendations(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	result, err := h.service.GetRecommendations(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	c.JSON(http.StatusOK, result)
}
