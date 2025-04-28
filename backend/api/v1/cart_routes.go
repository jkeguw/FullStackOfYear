package v1

import (
	"github.com/gin-gonic/gin"
	cartHandler "project/backend/handlers/cart"
)

// RegisterCartRoutes 注册购物车相关路由
func (r *Router) RegisterCartRoutes(router *gin.RouterGroup) {
	// 创建购物车处理器
	handler := cartHandler.NewHandler(r.cartService)

	// 购物车路由组 - 需要认证
	cartGroup := router.Group("/cart")
	cartGroup.Use(r.authMiddleware) // 添加认证中间件
	{
		// 获取购物车
		cartGroup.GET("", handler.GetCart)

		// 添加商品到购物车
		cartGroup.POST("", handler.AddToCart)

		// 更新购物车商品数量
		cartGroup.PATCH("/quantity", handler.UpdateQuantity)

		// 从购物车中移除商品
		cartGroup.DELETE("/:productID", handler.RemoveFromCart)

		// 清空购物车
		cartGroup.DELETE("", handler.ClearCart)
	}
}
