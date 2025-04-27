package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// OrderStatusEnum 订单状态枚举
type OrderStatusEnum string

const (
	// 订单状态
	OrderStatusPending   OrderStatusEnum = "pending"   // 待支付
	OrderStatusPaid      OrderStatusEnum = "paid"      // 已支付
	OrderStatusShipped   OrderStatusEnum = "shipped"   // 已发货
	OrderStatusDelivered OrderStatusEnum = "delivered" // 已送达
	OrderStatusCancelled OrderStatusEnum = "cancelled" // 已取消
	OrderStatusRefunded  OrderStatusEnum = "refunded"  // 已退款
)

// PaymentMethodEnum 支付方式枚举
type PaymentMethodEnum string

const (
	// 支付方式
	PaymentMethodCreditCard PaymentMethodEnum = "credit_card" // 信用卡
	PaymentMethodDebitCard  PaymentMethodEnum = "debit_card"  // 借记卡
	PaymentMethodPaypal     PaymentMethodEnum = "paypal"      // PayPal
	PaymentMethodAlipay     PaymentMethodEnum = "alipay"      // 支付宝
	PaymentMethodWechat     PaymentMethodEnum = "wechat"      // 微信支付
)

// Order 订单
type Order struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
	OrderNumber   string             `bson:"orderNumber" json:"orderNumber"` // 订单编号
	Status        OrderStatusEnum    `bson:"status" json:"status"`           // 订单状态
	Items         []OrderItem        `bson:"items" json:"items"`             // 订单商品
	ShippingInfo  ShippingInfo       `bson:"shippingInfo" json:"shippingInfo"` // 配送信息
	PaymentInfo   PaymentInfo        `bson:"paymentInfo" json:"paymentInfo"`   // 支付信息
	Subtotal      float64            `bson:"subtotal" json:"subtotal"`         // 商品小计
	ShippingFee   float64            `bson:"shippingFee" json:"shippingFee"`   // 配送费
	Tax           float64            `bson:"tax" json:"tax"`                   // 税费
	Discount      float64            `bson:"discount" json:"discount"`         // 折扣
	Total         float64            `bson:"total" json:"total"`               // 总价
	Notes         string             `bson:"notes,omitempty" json:"notes,omitempty"` // 备注
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
	PaidAt        *time.Time         `bson:"paidAt,omitempty" json:"paidAt,omitempty"`
	ShippedAt     *time.Time         `bson:"shippedAt,omitempty" json:"shippedAt,omitempty"`
	DeliveredAt   *time.Time         `bson:"deliveredAt,omitempty" json:"deliveredAt,omitempty"`
	CancelledAt   *time.Time         `bson:"cancelledAt,omitempty" json:"cancelledAt,omitempty"`
	CancelReason  string             `bson:"cancelReason,omitempty" json:"cancelReason,omitempty"`
}

// OrderItem 订单商品
type OrderItem struct {
	ProductID   primitive.ObjectID `bson:"productId" json:"productId"`     // 商品ID
	ProductType string             `bson:"productType" json:"productType"` // 商品类型
	Name        string             `bson:"name" json:"name"`               // 商品名称
	Price       float64            `bson:"price" json:"price"`             // 单价
	Quantity    int                `bson:"quantity" json:"quantity"`       // 数量
	Subtotal    float64            `bson:"subtotal" json:"subtotal"`       // 小计
	ImageURL    string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"` // 图片URL
}

// ShippingInfo 配送信息
type ShippingInfo struct {
	Name         string  `bson:"name" json:"name"`                   // 收件人姓名
	Phone        string  `bson:"phone" json:"phone"`                 // 电话
	Email        string  `bson:"email,omitempty" json:"email,omitempty"` // 邮箱
	Address      string  `bson:"address" json:"address"`             // 地址
	City         string  `bson:"city" json:"city"`                   // 城市
	State        string  `bson:"state" json:"state"`                 // 省/州
	ZipCode      string  `bson:"zipCode" json:"zipCode"`             // 邮编
	Country      string  `bson:"country" json:"country"`             // 国家
	ShippingMethod string `bson:"shippingMethod" json:"shippingMethod"` // 配送方式
}

// PaymentInfo 支付信息
type PaymentInfo struct {
	Method          PaymentMethodEnum `bson:"method" json:"method"`                   // 支付方式
	TransactionID   string            `bson:"transactionId,omitempty" json:"transactionId,omitempty"` // 交易ID
	LastFourDigits  string            `bson:"lastFourDigits,omitempty" json:"lastFourDigits,omitempty"` // 卡号后四位
	PaymentStatus   string            `bson:"paymentStatus" json:"paymentStatus"`     // 支付状态
	PaymentProvider string            `bson:"paymentProvider,omitempty" json:"paymentProvider,omitempty"` // 支付提供商
}

// 集合名常量
const (
	OrdersCollection = "orders"
)