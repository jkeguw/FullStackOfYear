// api/v1/routes.go

package v1

import (
	authHandler "FullStackOfYear/backend/handlers/auth"
	measurementHandler "FullStackOfYear/backend/handlers/measurement" // 新增
	"FullStackOfYear/backend/middleware"
	authService "FullStackOfYear/backend/services/auth"
	"FullStackOfYear/backend/services/email"
	measurementService "FullStackOfYear/backend/services/measurement" // 新增
	"github.com/gin-gonic/gin"
)

type Router struct {
	authService        authService.Service
	emailService       *email.Service
	measurementService measurementService.Service // 新增
}

func NewRouter(authService authService.Service, emailService *email.Service, measurementService measurementService.Service) *Router {
	return &Router{
		authService:        authService,
		emailService:       emailService,
		measurementService: measurementService, // 新增
	}
}

func (r *Router) RegisterRoutes(router *gin.RouterGroup) {
	emailHandler := authHandler.NewEmailVerificationHandler(r.authService, r.emailService)
	// 新增测量处理器实例化
	mHandler := measurementHandler.NewHandler(r.measurementService)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.POST("/logout", middleware.Auth(), authHandler.Logout)

		oauthGroup := authGroup.Group("/oauth/google")
		{
			oauthGroup.GET("/login", authHandler.OAuthInstance.InitiateOAuth)
			oauthGroup.GET("/callback", authHandler.OAuthInstance.HandleCallback)
		}

		authGroup.GET("/verify-email", emailHandler.VerifyEmail)

		emailGroup := authGroup.Group("")
		emailGroup.Use(middleware.Auth())
		{
			emailGroup.POST("/send-verification", emailHandler.SendVerification)
			emailGroup.POST("/update-email", emailHandler.UpdateEmail)
		}
	}

	authenticated := router.Group("")
	authenticated.Use(middleware.Auth())
	{
		userGroup := authenticated.Group("/users")
		{
			userGroup.GET("/me", nil)
			userGroup.PUT("/me", nil)
		}

		reviewGroup := authenticated.Group("/reviews")
		reviewGroup.Use(middleware.RequireRoles("reviewer", "admin"))
		{
			reviewGroup.GET("", nil)
			reviewGroup.POST("", nil)
		}

		deviceGroup := authenticated.Group("/devices")
		{
			deviceGroup.GET("", nil)
			deviceGroup.POST("", nil)
		}

		// 新增测量工具相关路由
		measurementGroup := authenticated.Group("/measurements")
		{
			measurementGroup.POST("", mHandler.CreateMeasurement)
			measurementGroup.GET("", mHandler.ListMeasurements)
			measurementGroup.GET("/:id", mHandler.GetMeasurement)
			measurementGroup.PUT("/:id", mHandler.UpdateMeasurement)
			measurementGroup.DELETE("/:id", mHandler.DeleteMeasurement)
			measurementGroup.GET("/stats", mHandler.GetUserStats)
			measurementGroup.GET("/recommend", mHandler.GetRecommendations)
		}
	}
}
