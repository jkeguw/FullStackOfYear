package v1

import (
	deviceHandler "project/backend/handlers/device"
	"project/backend/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterDeviceRoutes 注册设备相关路由
func (r *Router) RegisterDeviceRoutes(router *gin.RouterGroup) {
	dHandler := deviceHandler.NewHandler(r.deviceService)

	// 设备相关路由 - 部分公开
	deviceGroup := router.Group("/devices")
	{
		// 公开路由，无需认证
		deviceGroup.GET("", dHandler.ListDevices)
		deviceGroup.GET("/:id", dHandler.GetMouseDevice)
		
		// 鼠标比较和相似度查询
		deviceGroup.GET("/mice/compare", dHandler.CompareMice)
		deviceGroup.GET("/mice/:id/similar", dHandler.FindSimilarMice)
		
		// 需要认证的路由
		authDeviceGroup := deviceGroup.Group("")
		authDeviceGroup.Use(middleware.Auth())
		{
			// 需要管理员权限
			adminDeviceGroup := authDeviceGroup.Group("")
			adminDeviceGroup.Use(middleware.RequireRoles("admin"))
			{
				adminDeviceGroup.POST("/mouse", dHandler.CreateMouseDevice)
				adminDeviceGroup.PUT("/mouse/:id", dHandler.UpdateMouseDevice)
				adminDeviceGroup.DELETE("/:id", dHandler.DeleteDevice)
			}
		}
	}

	// 个人设备管理功能已移除
	// 以下代码保留以便将来可能的恢复
	/*
	// 用户设备配置相关路由
	userDeviceGroup := router.Group("/user-devices")
	{
		// 公开路由，查看公开的用户设备配置
		userDeviceGroup.GET("/public", dHandler.ListPublicUserDevices)
		
		// 需要认证的路由
		authUserDeviceGroup := userDeviceGroup.Group("")
		authUserDeviceGroup.Use(middleware.Auth())
		{
			authUserDeviceGroup.GET("", dHandler.ListUserDevices)
			authUserDeviceGroup.GET("/:id", dHandler.GetUserDevice)
			authUserDeviceGroup.POST("", dHandler.CreateUserDevice)
			authUserDeviceGroup.PUT("/:id", dHandler.UpdateUserDevice)
			authUserDeviceGroup.DELETE("/:id", dHandler.DeleteUserDevice)
		}
	}
	*/
}