package cart

import (
	"context"
	"project/backend/internal/errors"
	"project/backend/models"
)

// MockService 实现购物车服务接口的Mock版本
type MockService struct{}

// GetCart 获取用户的购物车
func (s *MockService) GetCart(ctx context.Context, userID string) (*models.Cart, error) {
	return &models.Cart{
		Items: []models.CartItem{},
	}, nil
}

// AddToCart 向购物车添加商品
func (s *MockService) AddToCart(ctx context.Context, userID string, item models.CartItem) error {
	return errors.NewInternalServerError("数据库连接失败，购物车服务暂不可用")
}

// UpdateQuantity 更新购物车商品数量
func (s *MockService) UpdateQuantity(ctx context.Context, userID string, productID string, quantity int) error {
	return errors.NewInternalServerError("数据库连接失败，购物车服务暂不可用")
}

// RemoveFromCart 从购物车移除商品
func (s *MockService) RemoveFromCart(ctx context.Context, userID string, productID string) error {
	return errors.NewInternalServerError("数据库连接失败，购物车服务暂不可用")
}

// ClearCart 清空购物车
func (s *MockService) ClearCart(ctx context.Context, userID string) error {
	return errors.NewInternalServerError("数据库连接失败，购物车服务暂不可用")
}
