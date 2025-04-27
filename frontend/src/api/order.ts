import { AxiosResponse } from 'axios';
import { request } from '../utils/request';
import { 
  CreateOrderRequest,
  OrderResponse,
  OrderListResponse,
  UpdateOrderStatusRequest,
  PaymentCompleteRequest
} from '../types/order';

/**
 * 创建订单
 * @param data 订单信息
 * @returns 订单响应
 */
export const createOrder = (data: CreateOrderRequest): Promise<AxiosResponse<OrderResponse>> => {
  return request({
    url: '/v1/orders',
    method: 'post',
    data
  });
};

/**
 * 获取订单详情
 * @param orderId 订单ID
 * @returns 订单详情
 */
export const getOrder = (orderId: string): Promise<AxiosResponse<OrderResponse>> => {
  return request({
    url: `/v1/orders/${orderId}`,
    method: 'get'
  });
};

/**
 * 获取订单列表
 * @param page 页码
 * @param pageSize 每页数量
 * @returns 订单列表
 */
export const getOrderList = (page = 1, pageSize = 10): Promise<AxiosResponse<OrderListResponse>> => {
  return request({
    url: '/v1/orders',
    method: 'get',
    params: {
      page,
      page_size: pageSize
    }
  });
};

/**
 * 更新订单状态
 * @param orderId 订单ID
 * @param data 状态更新数据
 * @returns 更新后的订单
 */
export const updateOrderStatus = (
  orderId: string,
  data: UpdateOrderStatusRequest
): Promise<AxiosResponse<OrderResponse>> => {
  return request({
    url: `/v1/orders/${orderId}/status`,
    method: 'patch',
    data
  });
};

/**
 * 处理订单支付
 * @param orderId 订单ID
 * @param data 支付数据
 * @returns 支付结果
 */
export const processPayment = (
  orderId: string,
  data: PaymentCompleteRequest
): Promise<AxiosResponse<OrderResponse>> => {
  return request({
    url: `/v1/orders/${orderId}/payment`,
    method: 'post',
    data
  });
};

/**
 * 获取订单统计信息
 * @returns 订单统计
 */
export const getOrderStats = (): Promise<AxiosResponse<any>> => {
  return request({
    url: '/v1/orders/stats',
    method: 'get'
  });
};