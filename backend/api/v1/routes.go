package v1

import (
	authHandler "FullStackOfYear/backend/handlers/auth" // 重命名为 authHandler
	"FullStackOfYear/backend/middleware"
	"FullStackOfYear/backend/services/auth" // auth service
	"FullStackOfYear/backend/services/email"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authService  *auth.AuthService
	emailService *email.Service
}

func NewRouter(authService *auth.AuthService, emailService *email.Service) *Router {
	return &Router{
		authService:  authService,
		emailService: emailService,
	}
}

func (r *Router) RegisterRoutes(router *gin.RouterGroup) {
	// init EmailVerificationHandler
	emailHandler := authHandler.NewEmailVerificationHandler(r.authService, r.emailService)

	// auth route
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.POST("/logout", middleware.Auth(), authHandler.Logout)

		// OAuth routes
		oauthGroup := authGroup.Group("/oauth/google")
		{
			oauthGroup.GET("/login", authHandler.OAuthInstance.InitiateOAuth)
			oauthGroup.GET("/callback", authHandler.OAuthInstance.HandleCallback)
		}

		// Email verification routes - 新增
		authGroup.GET("/verify-email", emailHandler.VerifyEmail) // 公开接口

		// Protected email routes - 新增
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
		// user related
		user := authenticated.Group("/users")
		{
			user.GET("/me", nil) // TODO
			user.PUT("/me", nil) // TODO
		}

		// comment related
		review := authenticated.Group("/reviews")
		review.Use(middleware.RequireRoles("reviewer", "admin"))
		{
			review.GET("", nil)  // TODO
			review.POST("", nil) // TODO
		}

		// devices related
		device := authenticated.Group("/devices")
		{
			device.GET("", nil)  // TODO
			device.POST("", nil) // TODO
		}
	}
}
