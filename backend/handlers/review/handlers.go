package review

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/backend/internal/errors"
	"project/backend/models"
	reviewService "project/backend/services/review"
	"project/backend/types/review"
)

type Handler struct {
	service reviewService.Service
}

func NewHandler(service reviewService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateReview 创建评测
func (h *Handler) CreateReview(c *gin.Context) {
	var request review.CreateReviewRequest
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

	// 检查用户角色 (需要评测员或管理员)
	role, roleExists := c.Get("userRole")
	if !roleExists || (role.(string) != "reviewer" && role.(string) != "admin") {
		c.JSON(http.StatusForbidden, errors.NewForbiddenError("只有评测员才能创建评测"))
		return
	}

	result, err := h.service.CreateReview(c.Request.Context(), userID.(string), request)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 构建响应
	response := mapReviewToResponse(result)
	c.JSON(http.StatusCreated, response)
}

// GetReview 获取评测详情
func (h *Handler) GetReview(c *gin.Context) {
	reviewID := c.Param("id")

	// 获取评测详情
	result, err := h.service.GetReview(c.Request.Context(), reviewID)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 增加查看次数
	_ = h.service.IncrementViewCount(c.Request.Context(), reviewID)

	// 构建响应
	response := mapReviewToResponse(result)
	c.JSON(http.StatusOK, response)
}

// UpdateReview 更新评测
func (h *Handler) UpdateReview(c *gin.Context) {
	reviewID := c.Param("id")

	var request review.UpdateReviewRequest
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

	result, err := h.service.UpdateReview(c.Request.Context(), userID.(string), reviewID, request)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 构建响应
	response := mapReviewToResponse(result)
	c.JSON(http.StatusOK, response)
}

// DeleteReview 删除评测
func (h *Handler) DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	err := h.service.DeleteReview(c.Request.Context(), userID.(string), reviewID)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListReviews 获取评测列表
func (h *Handler) ListReviews(c *gin.Context) {
	var request review.ReviewListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("无效的请求参数: "+err.Error()))
		return
	}

	result, err := h.service.ListReviews(c.Request.Context(), request)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserReviewStats 获取用户评测统计
func (h *Handler) GetUserReviewStats(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		// 如果未指定用户ID，使用当前用户ID
		currentUserID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
			return
		}
		userID = currentUserID.(string)
	}

	stats, err := h.service.GetUserReviewStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPendingReviews 获取待审核评测列表
func (h *Handler) GetPendingReviews(c *gin.Context) {
	// 从查询参数中获取筛选条件
	reviewType := c.Query("type")
	page := 1
	pageSize := 10

	result, err := h.service.GetPendingReviews(c.Request.Context(), reviewType, page, pageSize)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ApproveReview 批准评测
func (h *Handler) ApproveReview(c *gin.Context) {
	reviewID := c.Param("id")
	var request review.ReviewerActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("无效的请求: "+err.Error()))
		return
	}

	// 从上下文获取评审员ID
	reviewerID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	result, err := h.service.ApproveReview(c.Request.Context(), reviewerID.(string), reviewID, request.Notes)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 构建响应
	response := mapReviewToResponse(result)
	c.JSON(http.StatusOK, response)
}

// RejectReview 拒绝评测
func (h *Handler) RejectReview(c *gin.Context) {
	reviewID := c.Param("id")
	var request review.ReviewerActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("无效的请求: "+err.Error()))
		return
	}

	// 从上下文获取评审员ID
	reviewerID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	result, err := h.service.RejectReview(c.Request.Context(), reviewerID.(string), reviewID, request.Notes)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 构建响应
	response := mapReviewToResponse(result)
	c.JSON(http.StatusOK, response)
}

// FeaturedReview 设置为精选评测
func (h *Handler) FeaturedReview(c *gin.Context) {
	reviewID := c.Param("id")
	var request review.ReviewerActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("无效的请求: "+err.Error()))
		return
	}

	// 从上下文获取评审员ID
	reviewerID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("用户未认证"))
		return
	}

	result, err := h.service.FeaturedReview(c.Request.Context(), reviewerID.(string), reviewID, request.Rank)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), err)
		return
	}

	// 构建响应
	response := mapReviewToResponse(result)
	c.JSON(http.StatusOK, response)
}

// 辅助方法

// mapReviewToResponse 将Review模型转换为ReviewResponse
func mapReviewToResponse(r *models.Review) review.ReviewResponse {
	response := review.ReviewResponse{
		ID:             r.ID.Hex(),
		ExternalItemID: r.ExternalItemID.Hex(),
		ItemType:       r.ItemType,
		UserID:         r.UserID.Hex(),
		Content:        r.Content,
		Pros:           r.Pros,
		Cons:           r.Cons,
		Score:          r.Score,
		Usage:          r.Usage,
		Status:         string(r.Status),
		Type:           string(r.Type),
		ContentType:    string(r.ContentType),
		ViewCount:      r.ViewCount,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}

	if r.ReviewerID != nil {
		response.ReviewerID = r.ReviewerID.Hex()
	}

	response.ReviewerNotes = r.ReviewerNotes
	response.ReviewedAt = r.ReviewedAt
	response.PublishedAt = r.PublishedAt
	response.FeaturedRank = r.FeaturedRank

	return response
}
