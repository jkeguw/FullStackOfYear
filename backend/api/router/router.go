package router

import (
	"github.com/gin-gonic/gin"
	"project/backend/api/v1"
	"project/backend/middleware"
	"project/backend/services/auth"
	"project/backend/services/i18n"
	"project/backend/services/jwt"
)

func InitRouter(r *gin.Engine, authService auth.Service, jwtService jwt.Service, i18nService i18n.Service) {
	// 添加中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.XSSProtection())

	// 添加国际化中间件
	r.Use(middleware.I18n(i18nService))

	// 为每个请求存储i18n服务
	r.Use(func(c *gin.Context) {
		c.Set("i18n", i18nService)
		c.Next()
	})

	// 添加根健康检查端点，方便容器的健康检查
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// 添加调试端点
	r.GET("/api/debug", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "Debug API",
			"data": gin.H{
				"devices":  []gin.H{},
				"total":    0,
				"page":     1,
				"pageSize": 20,
			},
		})
	})

	// 添加JWT认证中间件，但不全局应用
	authMiddleware := middleware.Auth(jwtService)

	// 初始化API路由
	api := r.Group("/api")
	v1.RegisterRoutes(api, authService, jwtService, authMiddleware)
}
