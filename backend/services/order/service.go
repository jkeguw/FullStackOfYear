package order

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"project/backend/internal/errors"
	"project/backend/models"
	"project/backend/services/cart"
	"project/backend/services/device"
	"project/backend/types/order"
)

// 类型别名，避免使用带包名的类型引用
type (
	OrderResponse        = order.OrderResponse
	OrderItemResponse    = order.OrderItemResponse
	ShippingInfoResponse = order.ShippingInfoResponse
	PaymentInfoResponse  = order.PaymentInfoResponse
	OrderListResponse    = order.OrderListResponse
)

// Service 订单服务
type Service struct {
	db            *mongo.Database
	cartService   *cart.Service
	deviceService *device.Service
}

// NewService 创建订单服务
func NewService(db *mongo.Database, cartService *cart.Service, deviceService *device.Service) *Service {
	if db == nil {
		// 如果数据库为nil，返回一个基本服务实例
		// 各方法会检查db是否为nil
		return &Service{
			db:            nil,
			cartService:   nil,
			deviceService: nil,
		}
	}
	return &Service{
		db:            db,
		cartService:   cartService,
		deviceService: deviceService,
	}
}

// CreateOrder 创建订单
func (s *Service) CreateOrder(ctx context.Context, userID primitive.ObjectID, req order.CreateOrderRequest) (*models.Order, error) {
	// 检查数据库连接
	if s.db == nil {
		return nil, errors.NewInternalServerError("数据库连接失败，订单服务暂不可用")
	}

	// 验证请求
	if len(req.Items) == 0 {
		return nil, errors.NewBadRequestError("订单必须包含至少一件商品")
	}

	// 创建订单项
	orderItems := make([]models.OrderItem, 0, len(req.Items))
	subtotal := 0.0

	for _, item := range req.Items {
		productID, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return nil, errors.NewBadRequestError(fmt.Sprintf("无效的商品ID: %s", item.ProductID))
		}

		// 这里应该从产品服务获取商品信息
		// 简化起见，我们使用固定值
		price := 99.99                              // 这应该从产品服务获取
		name := "鼠标设备"                              // 这应该从产品服务获取
		imageURL := "https://example.com/image.jpg" // 这应该从产品服务获取

		// 计算小计
		itemSubtotal := price * float64(item.Quantity)
		subtotal += itemSubtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:   productID,
			ProductType: item.ProductType,
			Name:        name,
			Price:       price,
			Quantity:    item.Quantity,
			Subtotal:    itemSubtotal,
			ImageURL:    imageURL,
		})
	}

	// 计算费用
	shippingFee := 10.0    // 固定运费，实际应该根据配送方式和地址计算
	tax := subtotal * 0.05 // 5%的税率，实际应该根据地区计算
	discount := 0.0        // 无折扣
	total := subtotal + shippingFee + tax - discount

	// 创建订单
	now := time.Now()
	order := &models.Order{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		OrderNumber: generateOrderNumber(),
		Status:      models.OrderStatusPending,
		Items:       orderItems,
		ShippingInfo: models.ShippingInfo{
			Name:           req.ShippingInfo.Name,
			Phone:          req.ShippingInfo.Phone,
			Email:          req.ShippingInfo.Email,
			Address:        req.ShippingInfo.Address,
			City:           req.ShippingInfo.City,
			State:          req.ShippingInfo.State,
			ZipCode:        req.ShippingInfo.ZipCode,
			Country:        req.ShippingInfo.Country,
			ShippingMethod: req.ShippingInfo.ShippingMethod,
		},
		PaymentInfo: models.PaymentInfo{
			Method:        models.PaymentMethodEnum(req.PaymentMethod),
			PaymentStatus: "pending",
		},
		Subtotal:    subtotal,
		ShippingFee: shippingFee,
		Tax:         tax,
		Discount:    discount,
		Total:       total,
		Notes:       req.Notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 保存到数据库
	_, err := s.db.Collection(models.OrdersCollection).InsertOne(ctx, order)
	if err != nil {
		return nil, errors.NewInternalServerError("创建订单失败")
	}

	// 成功后应该清空购物车，但这里简化处理

	return order, nil
}

// GetOrder 获取订单详情
func (s *Service) GetOrder(ctx context.Context, userID, orderID primitive.ObjectID) (*models.Order, error) {
	// 检查数据库连接
	if s.db == nil {
		return nil, errors.NewInternalServerError("数据库连接失败，订单服务暂不可用")
	}

	var order models.Order

	err := s.db.Collection(models.OrdersCollection).FindOne(ctx, bson.M{
		"_id":    orderID,
		"userId": userID,
	}).Decode(&order)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("订单不存在")
		}
		return nil, errors.NewInternalServerError("获取订单失败")
	}

	return &order, nil
}

// GetOrderByNumber 通过订单号获取订单
func (s *Service) GetOrderByNumber(ctx context.Context, userID primitive.ObjectID, orderNumber string) (*models.Order, error) {
	// 检查数据库连接
	if s.db == nil {
		return nil, errors.NewInternalServerError("数据库连接失败，订单服务暂不可用")
	}

	var order models.Order

	err := s.db.Collection(models.OrdersCollection).FindOne(ctx, bson.M{
		"orderNumber": orderNumber,
		"userId":      userID,
	}).Decode(&order)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("订单不存在")
		}
		return nil, errors.NewInternalServerError("获取订单失败")
	}

	return &order, nil
}

// ListUserOrders 获取用户订单列表
func (s *Service) ListUserOrders(ctx context.Context, userID primitive.ObjectID, page, pageSize int) (*OrderListResponse, error) {
	// 检查数据库连接
	if s.db == nil {
		return nil, errors.NewInternalServerError("数据库连接失败，订单服务暂不可用")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize

	// 创建查询选项
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"createdAt": -1}) // 最新的订单优先

	// 执行查询
	cursor, err := s.db.Collection(models.OrdersCollection).Find(ctx, bson.M{"userId": userID}, findOptions)
	if err != nil {
		return nil, errors.NewInternalServerError("获取订单列表失败")
	}
	defer cursor.Close(ctx)

	// 解析结果
	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, errors.NewInternalServerError("解析订单列表失败")
	}

	// 获取总数
	total, err := s.db.Collection(models.OrdersCollection).CountDocuments(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, errors.NewInternalServerError("获取订单数量失败")
	}

	// 转换为响应格式
	orderResponses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = convertToOrderResponse(order)
	}

	return &OrderListResponse{
		Orders:      orderResponses,
		TotalCount:  total,
		CurrentPage: page,
		PageSize:    pageSize,
	}, nil
}

// UpdateOrderStatus 更新订单状态
func (s *Service) UpdateOrderStatus(ctx context.Context, userID, orderID primitive.ObjectID, status models.OrderStatusEnum, cancelReason string) (*models.Order, error) {
	// 检查数据库连接
	if s.db == nil {
		return nil, errors.NewInternalServerError("数据库连接失败，订单服务暂不可用")
	}

	// 先获取订单
	order, err := s.GetOrder(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}

	// 验证状态变更是否合法
	if !isValidStatusTransition(order.Status, status) {
		return nil, errors.NewBadRequestError(fmt.Sprintf("无法从 %s 状态变更为 %s 状态", order.Status, status))
	}

	// 准备更新数据
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}

	// 根据状态添加额外字段
	switch status {
	case models.OrderStatusPaid:
		now := time.Now()
		update["$set"].(bson.M)["paidAt"] = now
	case models.OrderStatusShipped:
		now := time.Now()
		update["$set"].(bson.M)["shippedAt"] = now
	case models.OrderStatusDelivered:
		now := time.Now()
		update["$set"].(bson.M)["deliveredAt"] = now
	case models.OrderStatusCancelled:
		now := time.Now()
		update["$set"].(bson.M)["cancelledAt"] = now
		update["$set"].(bson.M)["cancelReason"] = cancelReason
	}

	// 执行更新
	result := s.db.Collection(models.OrdersCollection).FindOneAndUpdate(
		ctx,
		bson.M{"_id": orderID, "userId": userID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		return nil, errors.NewInternalServerError("更新订单状态失败")
	}

	// 解析更新后的订单
	var updatedOrder models.Order
	if err := result.Decode(&updatedOrder); err != nil {
		return nil, errors.NewInternalServerError("解析更新后的订单失败")
	}

	return &updatedOrder, nil
}

// ProcessPayment 处理支付
func (s *Service) ProcessPayment(ctx context.Context, orderID primitive.ObjectID, transactionID, paymentStatus string) (*models.Order, error) {
	// 检查数据库连接
	if s.db == nil {
		return nil, errors.NewInternalServerError("数据库连接失败，订单服务暂不可用")
	}

	// 查找订单
	var order models.Order
	err := s.db.Collection(models.OrdersCollection).FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("订单不存在")
		}
		return nil, errors.NewInternalServerError("获取订单失败")
	}

	// 验证订单状态
	if order.Status != models.OrderStatusPending {
		return nil, errors.NewBadRequestError("只有待支付的订单可以处理支付")
	}

	// 更新支付信息
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"paymentInfo.transactionId": transactionID,
			"paymentInfo.paymentStatus": paymentStatus,
			"updatedAt":                 now,
		},
	}

	// 如果支付成功，更新订单状态为已支付
	if paymentStatus == "success" {
		update["$set"].(bson.M)["status"] = models.OrderStatusPaid
		update["$set"].(bson.M)["paidAt"] = now
	}

	// 执行更新
	result := s.db.Collection(models.OrdersCollection).FindOneAndUpdate(
		ctx,
		bson.M{"_id": orderID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		return nil, errors.NewInternalServerError("更新支付信息失败")
	}

	// 解析更新后的订单
	var updatedOrder models.Order
	if err := result.Decode(&updatedOrder); err != nil {
		return nil, errors.NewInternalServerError("解析更新后的订单失败")
	}

	return &updatedOrder, nil
}

// 生成订单号
func generateOrderNumber() string {
	now := time.Now()
	// 在Go 1.20+中不再需要显式设置随机数种子

	// 格式：年月日-随机数字
	return fmt.Sprintf("%s-%04d",
		now.Format("20060102"),
		rand.Intn(10000),
	)
}

// 验证订单状态变更是否合法
func isValidStatusTransition(from, to models.OrderStatusEnum) bool {
	// 定义允许的状态转换
	allowedTransitions := map[models.OrderStatusEnum][]models.OrderStatusEnum{
		models.OrderStatusPending: {
			models.OrderStatusPaid,
			models.OrderStatusCancelled,
		},
		models.OrderStatusPaid: {
			models.OrderStatusShipped,
			models.OrderStatusRefunded,
		},
		models.OrderStatusShipped: {
			models.OrderStatusDelivered,
		},
		models.OrderStatusDelivered: {
			models.OrderStatusRefunded,
		},
		// 已取消和已退款是终态，不允许再变更
	}

	allowed, exists := allowedTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}

	return false
}

// 将订单模型转换为响应格式
func convertToOrderResponse(order models.Order) OrderResponse {
	// 转换订单项
	items := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemResponse{
			ProductID:   item.ProductID.Hex(),
			ProductType: item.ProductType,
			Name:        item.Name,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
			ImageURL:    item.ImageURL,
		}
	}

	// 转换配送信息
	shippingInfo := ShippingInfoResponse{
		Name:           order.ShippingInfo.Name,
		Phone:          order.ShippingInfo.Phone,
		Email:          order.ShippingInfo.Email,
		Address:        order.ShippingInfo.Address,
		City:           order.ShippingInfo.City,
		State:          order.ShippingInfo.State,
		ZipCode:        order.ShippingInfo.ZipCode,
		Country:        order.ShippingInfo.Country,
		ShippingMethod: order.ShippingInfo.ShippingMethod,
	}

	// 转换支付信息
	paymentInfo := PaymentInfoResponse{
		Method:          string(order.PaymentInfo.Method),
		TransactionID:   order.PaymentInfo.TransactionID,
		LastFourDigits:  order.PaymentInfo.LastFourDigits,
		PaymentStatus:   order.PaymentInfo.PaymentStatus,
		PaymentProvider: order.PaymentInfo.PaymentProvider,
	}

	return OrderResponse{
		ID:           order.ID.Hex(),
		UserID:       order.UserID.Hex(),
		OrderNumber:  order.OrderNumber,
		Status:       string(order.Status),
		Items:        items,
		ShippingInfo: shippingInfo,
		PaymentInfo:  paymentInfo,
		Subtotal:     order.Subtotal,
		ShippingFee:  order.ShippingFee,
		Tax:          order.Tax,
		Discount:     order.Discount,
		Total:        order.Total,
		Notes:        order.Notes,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
		PaidAt:       order.PaidAt,
		ShippedAt:    order.ShippedAt,
		DeliveredAt:  order.DeliveredAt,
		CancelledAt:  order.CancelledAt,
		CancelReason: order.CancelReason,
	}
}
