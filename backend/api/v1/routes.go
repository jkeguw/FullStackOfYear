package v1

import (
	"FullStackOfYear/backend/handlers/auth"
	"FullStackOfYear/backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	// auth route
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/refresh", auth.RefreshToken)
		authGroup.POST("/logout", middleware.Auth(), auth.Logout)

		// OAuth routes
		oauthGroup := authGroup.Group("/oauth/google")
		{
			oauthGroup.GET("/login", auth.OAuthInstance.InitiateOAuth)
			oauthGroup.GET("/callback", auth.OAuthInstance.HandleCallback)
		}
	}

	// route need auth
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
