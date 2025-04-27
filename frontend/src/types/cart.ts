/**
 * 购物车商品
 */
export interface CartItem {
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
  items: CartItem[];
  total: number;
  item_count: number;
  updated_at: string;
}