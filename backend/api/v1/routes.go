package v1

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	// user
	user := router.Group("/users")
	{
		user.POST("/register", nil) // TODO
		user.POST("/login", nil)    // TODO
	}

	// comment
	review := router.Group("/reviews")
	{
		review.GET("", nil)  // TODO
		review.POST("", nil) // TODO
	}

	// devices
	device := router.Group("/devices")
	{
		device.GET("", nil)  // TODO
		device.POST("", nil) // TODO
	}
}
