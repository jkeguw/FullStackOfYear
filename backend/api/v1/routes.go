// api/v1/routes.go

package v1

import (
	authHandler "project/backend/handlers/auth"
	// cartHandler "project/backend/handlers/cart"
	deviceHandler "project/backend/handlers/device"
	// orderHandler "project/backend/handlers/order"
	userHandler "project/backend/handlers/user"
	authService "project/backend/services/auth"
	// cartService "project/backend/services/cart"
	deviceService "project/backend/services/device"
	"project/backend/services/email"
	"project/backend/services/i18n"
	"project/backend/services/jwt"
	// orderService "project/backend/services/order"
	reviewService "project/backend/services/review"
	userService "project/backend/services/user"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authService     authService.Service
	emailService    *email.Service
	deviceService   deviceService.Service
	userService     userService.Service
	reviewService   reviewService.Service
	i18nService     i18n.Service
	jwtService      jwt.Service
	authMiddleware  gin.HandlerFunc
	// cartService     cartService.Service
	// orderService    *orderService.Service
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
	// cartSvc cartService.Service,
	// orderSvc *orderService.Service,
) *Router {
	return &Router{
		authService:     authService,
		emailService:    emailService,
		deviceService:   deviceService,
		userService:     userService,
		reviewService:   reviewSvc,
		i18nService:     i18nService,
		jwtService:      jwtService,
		authMiddleware:  authMiddleware,
		// cartService:     cartSvc,
		// orderService:    orderSvc,
	}
}

func RegisterRoutes(router *gin.RouterGroup, authService authService.Service, jwtService jwt.Service, authMiddleware gin.HandlerFunc) {
	// 实例化服务 - 简化版本，使用空实现
	emailService := &email.Service{}
	deviceService := &deviceService.DefaultService{} // 使用实际结构体而非nil
	userService := &userService.DefaultService{} // 使用实际结构体而非nil
	reviewSvc := &reviewService.DefaultService{} // 使用实际结构体而非nil
	i18nSvc := &i18n.DefaultService{} // 使用实际结构体而非nil

	r := NewRouter(
		authService,
		emailService,
		deviceService,
		userService,
		reviewSvc,
		i18nSvc,
		jwtService,
		authMiddleware,
	)

	r.RegisterRoutes(router)
}

func (r *Router) RegisterRoutes(router *gin.RouterGroup) {
	emailHandler := authHandler.NewEmailVerificationHandler(r.authService, r.emailService)
	dHandler := deviceHandler.NewHandler(r.deviceService)
	uHandler := userHandler.NewHandler(r.userService)
	
	// 注册评测路由
	r.RegisterReviewRoutes(router)
	
	// 注册设备路由
	r.RegisterDeviceRoutes(router)
	
	// 注册语言路由
	r.registerLanguageRoutes(router)

	// 注册购物车路由
	r.RegisterCartRoutes(router)

	// 认证相关路由
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
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

	// 需要认证的路由
	authenticated := router.Group("")
	authenticated.Use(r.authMiddleware)
	{
		// 用户相关路由
		userGroup := authenticated.Group("/users")
		{
			userGroup.GET("/me", uHandler.GetUserProfile)
			userGroup.PUT("/me", uHandler.UpdateUserProfile)
		}
	}
	
	// 设备相关路由已移至device_routes.go

	// 设备评测相关路由 - 部分公开
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