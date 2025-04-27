package order

import (
	"time"
)

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Items         []OrderItemRequest `json:"items"`                   // 订单商品
	ShippingInfo  ShippingInfoRequest `json:"shippingInfo"`           // 配送信息
	PaymentMethod string             `json:"paymentMethod"`           // 支付方式
	Notes         string             `json:"notes,omitempty"`         // 备注
}

// OrderItemRequest 订单商品请求
type OrderItemRequest struct {
	ProductID   string  `json:"productId"`     // 商品ID
	ProductType string  `json:"productType"`   // 商品类型
	Quantity    int     `json:"quantity"`      // 数量
}

// ShippingInfoRequest 配送信息请求
type ShippingInfoRequest struct {
	Name           string `json:"name"`             // 收件人姓名
	Phone          string `json:"phone"`            // 电话
	Email          string `json:"email,omitempty"`  // 邮箱
	Address        string `json:"address"`          // 地址
	City           string `json:"city"`             // 城市
	State          string `json:"state"`            // 省/州
	ZipCode        string `json:"zipCode"`          // 邮编
	Country        string `json:"country"`          // 国家
	ShippingMethod string `json:"shippingMethod"`   // 配送方式
}

// UpdateOrderStatusRequest 更新订单状态请求
type UpdateOrderStatusRequest struct {
	Status       string `json:"status"`                      // 订单状态
	CancelReason string `json:"cancelReason,omitempty"`      // 取消原因(仅在取消订单时需要)
}

// OrderResponse 订单响应
type OrderResponse struct {
	ID            string             `json:"id"`
	UserID        string             `json:"userId"`
	OrderNumber   string             `json:"orderNumber"`         // 订单编号
	Status        string             `json:"status"`              // 订单状态
	Items         []OrderItemResponse `json:"items"`              // 订单商品
	ShippingInfo  ShippingInfoResponse `json:"shippingInfo"`      // 配送信息
	PaymentInfo   PaymentInfoResponse `json:"paymentInfo"`        // 支付信息
	Subtotal      float64            `json:"subtotal"`            // 商品小计
	ShippingFee   float64            `json:"shippingFee"`         // 配送费
	Tax           float64            `json:"tax"`                 // 税费
	Discount      float64            `json:"discount"`            // 折扣
	Total         float64            `json:"total"`               // 总价
	Notes         string             `json:"notes,omitempty"`     // 备注
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	PaidAt        *time.Time         `json:"paidAt,omitempty"`
	ShippedAt     *time.Time         `json:"shippedAt,omitempty"`
	DeliveredAt   *time.Time         `json:"deliveredAt,omitempty"`
	CancelledAt   *time.Time         `json:"cancelledAt,omitempty"`
	CancelReason  string             `json:"cancelReason,omitempty"`
}

// OrderItemResponse 订单商品响应
type OrderItemResponse struct {
	ProductID   string  `json:"productId"`     // 商品ID
	ProductType string  `json:"productType"`   // 商品类型
	Name        string  `json:"name"`          // 商品名称
	Price       float64 `json:"price"`         // 单价
	Quantity    int     `json:"quantity"`      // 数量
	Subtotal    float64 `json:"subtotal"`      // 小计
	ImageURL    string  `json:"imageUrl,omitempty"` // 图片URL
}

// ShippingInfoResponse 配送信息响应
type ShippingInfoResponse struct {
	Name           string `json:"name"`            // 收件人姓名
	Phone          string `json:"phone"`           // 电话
	Email          string `json:"email,omitempty"` // 邮箱
	Address        string `json:"address"`         // 地址
	City           string `json:"city"`            // 城市
	State          string `json:"state"`           // 省/州
	ZipCode        string `json:"zipCode"`         // 邮编
	Country        string `json:"country"`         // 国家
	ShippingMethod string `json:"shippingMethod"`  // 配送方式
}

// PaymentInfoResponse 支付信息响应
type PaymentInfoResponse struct {
	Method          string `json:"method"`                     // 支付方式
	TransactionID   string `json:"transactionId,omitempty"`    // 交易ID
	LastFourDigits  string `json:"lastFourDigits,omitempty"`   // 卡号后四位
	PaymentStatus   string `json:"paymentStatus"`              // 支付状态
	PaymentProvider string `json:"paymentProvider,omitempty"`  // 支付提供商
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Orders      []OrderResponse `json:"orders"`
	TotalCount  int64           `json:"totalCount"`
	CurrentPage int             `json:"currentPage"`
	PageSize    int             `json:"pageSize"`
}

// PaymentCompleteRequest 支付完成请求
type PaymentCompleteRequest struct {
	OrderID       string `json:"orderId"`        // 订单ID
	TransactionID string `json:"transactionId"`  // 交易ID
	PaymentStatus string `json:"paymentStatus"`  // 支付状态
}