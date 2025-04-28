// api/v1/routes.go

package v1

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"project/backend/config"
	authHandler "project/backend/handlers/auth"
	deviceHandler "project/backend/handlers/device"
	orderHandler "project/backend/handlers/order"
	userHandler "project/backend/handlers/user"
	"project/backend/internal/database"
	authService "project/backend/services/auth"
	cartService "project/backend/services/cart"
	deviceService "project/backend/services/device"
	"project/backend/services/email"
	"project/backend/services/i18n"
	"project/backend/services/jwt"
	orderService "project/backend/services/order"
	reviewService "project/backend/services/review"
	userService "project/backend/services/user"
)

type Router struct {
	authService    authService.Service
	emailService   *email.Service
	deviceService  deviceService.Service
	userService    userService.Service
	reviewService  reviewService.Service
	i18nService    i18n.Service
	jwtService     jwt.Service
	authMiddleware gin.HandlerFunc
	cartService    cartService.Service
	orderService   *orderService.Service
}

func NewRouter(
	authService authService.Service,
	emailService *email.Service,
	deviceService deviceService.Service,
	userService userService.Service,
	reviewSvc reviewService.Service,
	i18nService i18n.Service,
	jwtService jwt.Service,
	authMiddleware gin.HandlerFunc,
	cartSvc cartService.Service,
	orderSvc *orderService.Service,
) *Router {
	return &Router{
		authService:    authService,
		emailService:   emailService,
		deviceService:  deviceService,
		userService:    userService,
		reviewService:  reviewSvc,
		i18nService:    i18nService,
		jwtService:     jwtService,
		authMiddleware: authMiddleware,
		cartService:    cartSvc,
		orderService:   orderSvc,
	}
}

func RegisterRoutes(router *gin.RouterGroup, authService authService.Service, jwtService jwt.Service, authMiddleware gin.HandlerFunc) {
	// 不在这里添加健康检查端点，避免路由重复

	// 获取MongoDB数据库连接
	var db *mongo.Database
	if database.MongoClient != nil {
		db = database.MongoClient.Database(config.GetConfig().MongoDB.Database)
	} else {
		// 如果MongoDB连接失败，使用Mock服务
		log.Println("警告: MongoDB连接不可用，服务将使用空实现")
	}

	emailService := &email.Service{}
	deviceService := deviceService.New(db) // 使用实际的MongoDB连接，即使数据库连接为nil也使用完整实现
	userService := &userService.DefaultService{}
	reviewSvc := &reviewService.DefaultService{}
	i18nSvc := i18n.NewService() // 使用工厂方法创建i18n服务

	// 安全地创建服务
	var cartSvc cartService.Service
	var orderSvc *orderService.Service

	if db != nil {
		cartSvc = cartService.NewService(db)
		orderSvc = orderService.NewService(db, nil, nil) // 购物车和设备服务暂时为nil
	} else {
		// 使用mock实现避免空指针
		cartSvc = &cartService.MockService{}
		orderSvc = orderService.NewEmptyService()
	}

	r := NewRouter(
		authService,
		emailService,
		deviceService,
		userService,
		reviewSvc,
		i18nSvc,
		jwtService,
		authMiddleware,
		cartSvc,
		orderSvc,
	)

	r.RegisterRoutes(router)
}

func (r *Router) RegisterRoutes(router *gin.RouterGroup) {
	emailHandler := authHandler.NewEmailVerificationHandler(r.authService, r.emailService)
	dHandler := deviceHandler.NewHandler(r.deviceService)
	uHandler := userHandler.NewHandler(r.userService)
	oHandler := orderHandler.NewHandler(r.orderService)

	// 评测
	r.RegisterReviewRoutes(router)

	// 设备
	r.RegisterDeviceRoutes(router)

	// 语言
	r.registerLanguageRoutes(router)

	// 购物车
	r.RegisterCartRoutes(router)

	// 订单
	r.RegisterOrderRoutes(router, oHandler)

	// 添加 authService 到上下文的中间件
	authServiceMiddleware := func(c *gin.Context) {
		c.Set("authService", r.authService)
		c.Next()
	}

	// 认证
	authGroup := router.Group("/auth")
	authGroup.Use(authServiceMiddleware) // 注入认证服务
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		// OAuth
		oAuthHandler := authHandler.NewOAuthHandler(r.authService)
		authGroup.GET("/oauth/google", oAuthHandler.HandleOAuthLogin)
		authGroup.GET("/oauth/google/callback", oAuthHandler.HandleOAuthCallback)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.POST("/logout", r.authMiddleware, authHandler.Logout)

		authGroup.GET("/verify-email", emailHandler.VerifyEmail)

		emailGroup := authGroup.Group("")
		emailGroup.Use(r.authMiddleware)
		{
			emailGroup.POST("/send-verification", emailHandler.SendVerification)
			emailGroup.POST("/update-email", emailHandler.UpdateEmail)
		}
	}

	// 需要认证
	authenticated := router.Group("")
	authenticated.Use(r.authMiddleware)
	{
		// 用户
		userGroup := authenticated.Group("/users")
		{
			userGroup.GET("/me", uHandler.GetUserProfile)
			userGroup.PUT("/me", uHandler.UpdateUserProfile)
		}
	}

	// 设备相关路由已移至device_routes.go

	// 设备评测
	reviewsGroup := router.Group("/device-reviews")
	{
		// 公开路由，无需认证
		reviewsGroup.GET("", dHandler.ListDeviceReviews)
		reviewsGroup.GET("/:id", dHandler.GetDeviceReview)

		// 创建评测路由
		authReviewsGroup := reviewsGroup.Group("")
		authReviewsGroup.Use(r.authMiddleware)
		{
			authReviewsGroup.POST("", dHandler.CreateDeviceReview)
			authReviewsGroup.PUT("/:id", dHandler.UpdateDeviceReview)
			authReviewsGroup.DELETE("/:id", dHandler.DeleteDeviceReview)
		}
	}
}
