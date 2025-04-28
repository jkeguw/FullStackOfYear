package order

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"project/backend/internal/errors"
	"project/backend/models"
	"project/backend/types/order"
)

// 创建一个空的订单服务，用于数据库连接失败时的降级
func NewEmptyService() *Service {
	return &Service{
		db:            nil,
		cartService:   nil,
		deviceService: nil,
	}
}

// CreateOrder 创建订单
func (s *Service) CreateOrderMock(ctx context.Context, userID primitive.ObjectID, req order.CreateOrderRequest) (*models.Order, error) {
	if s.db == nil {
		return nil, errors.NewInternalServerError("数据库连接失败，订单服务暂不可用")
	}
	return nil, nil
}

// GetOrder 获取订单详情
func (s *Service) GetOrderMock(ctx context.Context, userID, orderID primitive.ObjectID) (*models.Order, error) {
	if s.db == nil {
		return nil, errors.NewInternalServerError("数据库连接失败，订单服务暂不可用")
	}
	return nil, nil
}

// ListUserOrders 获取用户订单列表
func (s *Service) ListUserOrdersMock(ctx context.Context, userID primitive.ObjectID, page, pageSize int) (*OrderListResponse, error) {
	if s.db == nil {
		return nil, errors.NewInternalServerError("数据库连接失败，订单服务暂不可用")
	}

	// 返回空列表
	return &OrderListResponse{
		Orders:      []OrderResponse{},
		TotalCount:  0,
		CurrentPage: page,
		PageSize:    pageSize,
	}, nil
}
