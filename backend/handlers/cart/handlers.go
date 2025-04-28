package cart

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"project/backend/models"
	"project/backend/services/cart"
	cartTypes "project/backend/types/cart"
)

// Handler 购物车处理器
type Handler struct {
	cartService cart.Service
}

// NewHandler 创建购物车处理器实例
func NewHandler(cartService cart.Service) *Handler {
	return &Handler{
		cartService: cartService,
	}
}

// GetCart 获取购物车
func (h *Handler) GetCart(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权访问",
			"data":    nil,
		})
		return
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 计算商品总数和总价
	total := 0.0
	itemCount := 0
	items := make([]cartTypes.CartItemResponse, 0, len(cart.Items))

	for _, item := range cart.Items {
		total += item.Price * float64(item.Quantity)
		itemCount += item.Quantity
		items = append(items, cartTypes.CartItemResponse{
			ProductID:   item.ProductID.Hex(),
			ProductType: item.ProductType,
			Name:        item.Name,
			Price:       item.Price,
			Quantity:    item.Quantity,
			ImageURL:    item.ImageURL,
		})
	}

	response := cartTypes.CartResponse{
		ID:        cart.ID.Hex(),
		Items:     items,
		Total:     total,
		ItemCount: itemCount,
		UpdatedAt: cart.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    response,
	})
}

// AddToCart 添加商品到购物车
func (h *Handler) AddToCart(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权访问",
			"data":    nil,
		})
		return
	}

	var req cartTypes.CartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的请求数据",
			"data":    nil,
		})
		return
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的商品ID",
			"data":    nil,
		})
		return
	}

	cartItem := models.CartItem{
		ProductID:   productID,
		ProductType: req.ProductType,
		Name:        req.Name,
		Price:       req.Price,
		Quantity:    req.Quantity,
		ImageURL:    req.ImageURL,
	}

	err = h.cartService.AddToCart(c.Request.Context(), userID.(string), cartItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "商品已添加到购物车",
		"data":    nil,
	})
}

// UpdateQuantity 更新购物车商品数量
func (h *Handler) UpdateQuantity(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权访问",
			"data":    nil,
		})
		return
	}

	var req cartTypes.UpdateQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的请求数据",
			"data":    nil,
		})
		return
	}

	err := h.cartService.UpdateQuantity(c.Request.Context(), userID.(string), req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "商品数量已更新",
		"data":    nil,
	})
}

// RemoveFromCart 从购物车移除商品
func (h *Handler) RemoveFromCart(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权访问",
			"data":    nil,
		})
		return
	}

	productID := c.Param("productID")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "缺少商品ID",
			"data":    nil,
		})
		return
	}

	err := h.cartService.RemoveFromCart(c.Request.Context(), userID.(string), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "商品已从购物车移除",
		"data":    nil,
	})
}

// ClearCart 清空购物车
func (h *Handler) ClearCart(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权访问",
			"data":    nil,
		})
		return
	}

	err := h.cartService.ClearCart(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "购物车已清空",
		"data":    nil,
	})
}
