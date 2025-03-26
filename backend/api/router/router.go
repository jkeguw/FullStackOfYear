package router

import (
	"FullStackOfYear/backend/api/v1"
	"FullStackOfYear/backend/middleware"
	"FullStackOfYear/backend/services/auth"
	"FullStackOfYear/backend/services/email"
	"FullStackOfYear/backend/services/measurement" // 新增
	"github.com/gin-gonic/gin"
)

func InitRouter(authService auth.Service, emailService *email.Service, measurementService measurement.Service) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 创建 router 实例并传入三个服务
	v1Router := v1.NewRouter(authService, emailService, measurementService)

	apiV1 := r.Group("/api/v1")
	v1Router.RegisterRoutes(apiV1)

	return r
}
