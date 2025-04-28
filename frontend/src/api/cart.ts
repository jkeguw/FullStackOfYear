import { request } from '@/utils/request';
import type { CartItem, CartResponse } from '@/types/cart';

/**
 * 获取购物车
 */
export function getCart() {
  return request<CartResponse>({
    url: '/api/cart',
    method: 'get'
  });
}

/**
 * 添加商品到购物车
 */
export function addToCart(item: CartItem) {
  return request<{ message: string }>({
    url: '/api/cart',
    method: 'post',
    data: {
      product_id: item.id,
      product_type: item.type || 'mouse',
      name: item.name,
      price: item.price,
      quantity: item.quantity,
      image_url: item.image
    }
  });
}

/**
 * 更新购物车商品数量
 */
export function updateQuantity(productId: string, quantity: number) {
  return request<{ message: string }>({
    url: '/api/cart/quantity',
    method: 'patch',
    data: {
      product_id: productId,
      quantity
    }
  });
}

/**
 * 从购物车移除商品
 */
export function removeFromCart(productId: string) {
  return request<{ message: string }>({
    url: `/api/cart/${productId}`,
    method: 'delete'
  });
}

/**
 * 清空购物车
 */
export function clearCart() {
  return request<{ message: string }>({
    url: '/api/cart',
    method: 'delete'
  });
}
