/**
 * 创建订单请求
 */
export interface CreateOrderRequest {
  items: OrderItemRequest[];
  shippingInfo: ShippingInfoRequest;
  paymentMethod: string;
  notes?: string;
}

/**
 * 订单商品请求
 */
export interface OrderItemRequest {
  productId: string;
  productType: string;
  quantity: number;
}

/**
 * 配送信息请求
 */
export interface ShippingInfoRequest {
  name: string;
  phone: string;
  email?: string;
  address: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
  shippingMethod: string;
}

/**
 * 更新订单状态请求
 */
export interface UpdateOrderStatusRequest {
  status: string;
  cancelReason?: string;
}

/**
 * 支付完成请求
 */
export interface PaymentCompleteRequest {
  orderId: string;
  transactionId: string;
  paymentStatus: string;
}

/**
 * 订单响应
 */
export interface OrderResponse {
  id: string;
  userId: string;
  orderNumber: string;
  status: string;
  items: OrderItemResponse[];
  shippingInfo: ShippingInfoResponse;
  paymentInfo: PaymentInfoResponse;
  subtotal: number;
  shippingFee: number;
  tax: number;
  discount: number;
  total: number;
  notes?: string;
  createdAt: string;
  updatedAt: string;
  paidAt?: string;
  shippedAt?: string;
  deliveredAt?: string;
  cancelledAt?: string;
  cancelReason?: string;
}

/**
 * 订单商品响应
 */
export interface OrderItemResponse {
  productId: string;
  productType: string;
  name: string;
  price: number;
  quantity: number;
  subtotal: number;
  imageUrl?: string;
}

/**
 * 配送信息响应
 */
export interface ShippingInfoResponse {
  name: string;
  phone: string;
  email?: string;
  address: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
  shippingMethod: string;
}

/**
 * 支付信息响应
 */
export interface PaymentInfoResponse {
  method: string;
  transactionId?: string;
  lastFourDigits?: string;
  paymentStatus: string;
  paymentProvider?: string;
}

/**
 * 订单列表响应
 */
export interface OrderListResponse {
  orders: OrderResponse[];
  totalCount: number;
  currentPage: number;
  pageSize: number;
}

/**
 * 订单状态枚举
 */
export enum OrderStatus {
  PENDING = 'pending',
  PAID = 'paid',
  SHIPPED = 'shipped',
  DELIVERED = 'delivered',
  CANCELLED = 'cancelled',
  REFUNDED = 'refunded'
}

/**
 * 支付方式枚举
 */
export enum PaymentMethod {
  CREDIT_CARD = 'credit_card',
  DEBIT_CARD = 'debit_card',
  PAYPAL = 'paypal',
  ALIPAY = 'alipay',
  WECHAT = 'wechat'
}

/**
 * 订单状态文本映射
 */
export const OrderStatusText = {
  [OrderStatus.PENDING]: '待支付',
  [OrderStatus.PAID]: '已支付',
  [OrderStatus.SHIPPED]: '已发货',
  [OrderStatus.DELIVERED]: '已送达',
  [OrderStatus.CANCELLED]: '已取消',
  [OrderStatus.REFUNDED]: '已退款'
};

/**
 * 支付方式文本映射
 */
export const PaymentMethodText = {
  [PaymentMethod.CREDIT_CARD]: '信用卡',
  [PaymentMethod.DEBIT_CARD]: '借记卡',
  [PaymentMethod.PAYPAL]: 'PayPal',
  [PaymentMethod.ALIPAY]: '支付宝',
  [PaymentMethod.WECHAT]: '微信支付'
};