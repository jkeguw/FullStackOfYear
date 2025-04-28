package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/backend/handlers/order"
	"project/backend/internal/errors"
)

// RegisterOrderRoutes 注册订单相关路由
func (r *Router) RegisterOrderRoutes(router *gin.RouterGroup, handler *order.Handler) {
	orderRoutes := router.Group("/orders")
	orderRoutes.Use(r.authMiddleware) // 添加认证中间件 - 订单路由需要认证
	{
		// 创建订单
		orderRoutes.POST("", handler.CreateOrder)

		// 获取订单详情
		orderRoutes.GET("/:id", handler.GetOrder)

		// 获取订单列表
		orderRoutes.GET("", handler.ListUserOrders)

		// 通过订单号获取订单
		orderRoutes.GET("/number/:number", handler.GetOrderByNumber)

		// 更新订单状态
		orderRoutes.PATCH("/:id/status", handler.UpdateOrderStatus)

		// 处理支付
		orderRoutes.POST("/:id/payment", handler.ProcessPayment)

		// 获取订单统计
		orderRoutes.GET("/stats", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, errors.NewAppError(errors.NotImplemented, "Order stats not implemented"))
		})
	}
}
