package v1

import (
	"project/backend/handlers/review"
	"project/backend/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterReviewRoutes 注册评测相关路由
func (r *Router) RegisterReviewRoutes(router *gin.RouterGroup) {
	reviewHandler := review.NewHandler(r.reviewService)

	// 评测相关路由 - 部分公开
	reviewsGroup := router.Group("/reviews")
	{
		// 公开路由，无需认证
		reviewsGroup.GET("", reviewHandler.ListReviews)
		reviewsGroup.GET("/:id", reviewHandler.GetReview)
		
		// 需要认证的路由
		authReviewsGroup := reviewsGroup.Group("")
		authReviewsGroup.Use(middleware.Auth())
		{
			// 需要评测员或管理员权限的路由
			reviewerGroup := authReviewsGroup.Group("")
			reviewerGroup.Use(middleware.RequireRoles("reviewer", "admin"))
			{
				reviewerGroup.POST("", reviewHandler.CreateReview)
				reviewerGroup.PUT("/:id", reviewHandler.UpdateReview)
				reviewerGroup.DELETE("/:id", reviewHandler.DeleteReview)
			}
			
			// 获取用户评测统计
			authReviewsGroup.GET("/stats/:id", reviewHandler.GetUserReviewStats)
			authReviewsGroup.GET("/stats", reviewHandler.GetUserReviewStats)
		}
	}
	
	// 评测审核相关路由 - 仅限评测员和管理员
	reviewerGroup := router.Group("/reviewer")
	reviewerGroup.Use(middleware.Auth(), middleware.RequireRoles("reviewer", "admin"))
	{
		// 待审核评测列表
		reviewerGroup.GET("/reviews/pending", reviewHandler.GetPendingReviews)
		
		// 评测审核操作
		reviewerGroup.POST("/reviews/:id/approve", reviewHandler.ApproveReview)
		reviewerGroup.POST("/reviews/:id/reject", reviewHandler.RejectReview)
		reviewerGroup.POST("/reviews/:id/featured", reviewHandler.FeaturedReview)
	}
}