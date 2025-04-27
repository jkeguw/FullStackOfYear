package review

import (
	"context"
	"project/backend/models"
	"project/backend/types/review"
)

// Service 定义评测服务接口
type Service interface {
	// 基本评测操作
	CreateReview(ctx context.Context, userID string, request review.CreateReviewRequest) (*models.Review, error)
	GetReview(ctx context.Context, reviewID string) (*models.Review, error)
	UpdateReview(ctx context.Context, userID, reviewID string, request review.UpdateReviewRequest) (*models.Review, error)
	DeleteReview(ctx context.Context, userID, reviewID string) error
	ListReviews(ctx context.Context, request review.ReviewListRequest) (*review.ReviewListResponse, error)
	
	// 评测员操作
	ApproveReview(ctx context.Context, reviewerID, reviewID string, notes string) (*models.Review, error)
	RejectReview(ctx context.Context, reviewerID, reviewID string, notes string) (*models.Review, error)
	FeaturedReview(ctx context.Context, reviewerID, reviewID string, rank int) (*models.Review, error)
	GetPendingReviews(ctx context.Context, reviewType string, page, pageSize int) (*review.ReviewListResponse, error)
	
	// 统计相关
	GetUserReviewStats(ctx context.Context, userID string) (*review.UserReviewStats, error)
	IncrementViewCount(ctx context.Context, reviewID string) error
}