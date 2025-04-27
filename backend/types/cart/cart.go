package cart

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CartItemRequest 用于添加/更新购物车项的请求
type CartItemRequest struct {
	ProductID   string  `json:"product_id" binding:"required"`
	ProductType string  `json:"product_type" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1"`
	ImageURL    string  `json:"image_url,omitempty"`
}

// UpdateQuantityRequest 用于更新购物车项数量的请求
type UpdateQuantityRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=0"`
}

// CartItemResponse 表示购物车项的响应
type CartItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductType string  `json:"product_type"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	ImageURL    string  `json:"image_url,omitempty"`
}

// CartResponse 表示购物车的响应
type CartResponse struct {
	ID        string            `json:"id,omitempty"`
	Items     []CartItemResponse `json:"items"`
	Total     float64           `json:"total"`
	ItemCount int               `json:"item_count"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// CartDetail 包含购物车完整详情
type CartDetail struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id"`
	Items     []CartItemResponse `json:"items"`
	Total     float64           `json:"total"`
	ItemCount int               `json:"item_count"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}