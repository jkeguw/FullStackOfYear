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

	// 添加JWT认证中间件，但不全局应用
	authMiddleware := middleware.Auth(jwtService)

	// 初始化API路由
	api := r.Group("/api")
	v1.RegisterRoutes(api, authService, jwtService, authMiddleware)
}