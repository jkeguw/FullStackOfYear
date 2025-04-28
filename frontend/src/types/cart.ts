/**
 * 前端购物车商品类型 - 用于添加到购物车
 */
export interface CartItem {
  id: string;
  type?: string;
  name: string;
  price: number;
  quantity: number;
  image?: string;
}

/**
 * 后端返回的购物车商品类型
 */
export interface CartItemResponse {
  product_id: string;
  product_type: string;
  name: string;
  price: number;
  quantity: number;
  image_url?: string;
}

/**
 * 购物车响应
 */
export interface CartResponse {
  id?: string;
  items: CartItemResponse[];
  total: number;
  item_count: number;
  updated_at: string;
}
