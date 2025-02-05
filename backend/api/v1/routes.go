package v1

import (
	authHandler "FullStackOfYear/backend/handlers/auth"
	"FullStackOfYear/backend/middleware"
	authService "FullStackOfYear/backend/services/auth"
	"FullStackOfYear/backend/services/email"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authService  authService.Service
	emailService *email.Service
}

func NewRouter(authService authService.Service, emailService *email.Service) *Router {
	return &Router{
		authService:  authService,
		emailService: emailService,
	}
}

func (r *Router) RegisterRoutes(router *gin.RouterGroup) {
	emailHandler := authHandler.NewEmailVerificationHandler(r.authService, r.emailService)

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
	}
}
