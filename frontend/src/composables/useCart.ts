import { ref, computed } from 'vue';
import { ElMessage } from 'element-plus';
import type { CartItem, CartResponse } from '@/types/cart';
import * as cartApi from '@/api/cart';
import { useAuth } from './useAuth';

export function useCart() {
  const { isAuthenticated } = useAuth();
  const cart = ref<CartResponse>({
    items: [],
    total: 0,
    item_count: 0,
    updated_at: '',
  });
  const loading = ref(false);

  // 购物车商品总数
  const itemCount = computed(() => cart.value.item_count);
  
  // 购物车总价
  const total = computed(() => cart.value.total);

  /**
   * 获取购物车
   */
  async function fetchCart() {
    if (!isAuthenticated.value) return;
    
    loading.value = true;
    try {
      const response = await cartApi.getCart();
      cart.value = response;
    } catch (error) {
      console.error('获取购物车失败', error);
    } finally {
      loading.value = false;
    }
  }

  /**
   * 添加商品到购物车
   */
  async function addToCart(item: CartItem) {
    if (!isAuthenticated.value) {
      ElMessage.warning('请先登录');
      return false;
    }

    loading.value = true;
    try {
      await cartApi.addToCart(item);
      ElMessage.success('已添加到购物车');
      await fetchCart();
      return true;
    } catch (error) {
      console.error('添加到购物车失败', error);
      ElMessage.error('添加到购物车失败');
      return false;
    } finally {
      loading.value = false;
    }
  }

  /**
   * 更新购物车商品数量
   */
  async function updateQuantity(productId: string, quantity: number) {
    if (!isAuthenticated.value) return false;

    loading.value = true;
    try {
      await cartApi.updateQuantity(productId, quantity);
      await fetchCart();
      return true;
    } catch (error) {
      console.error('更新数量失败', error);
      ElMessage.error('更新数量失败');
      return false;
    } finally {
      loading.value = false;
    }
  }

  /**
   * 从购物车移除商品
   */
  async function removeFromCart(productId: string) {
    if (!isAuthenticated.value) return false;

    loading.value = true;
    try {
      await cartApi.removeFromCart(productId);
      ElMessage.success('已从购物车移除');
      await fetchCart();
      return true;
    } catch (error) {
      console.error('移除购物车商品失败', error);
      ElMessage.error('移除购物车商品失败');
      return false;
    } finally {
      loading.value = false;
    }
  }

  /**
   * 清空购物车
   */
  async function clearCart() {
    if (!isAuthenticated.value) return false;

    loading.value = true;
    try {
      await cartApi.clearCart();
      ElMessage.success('购物车已清空');
      await fetchCart();
      return true;
    } catch (error) {
      console.error('清空购物车失败', error);
      ElMessage.error('清空购物车失败');
      return false;
    } finally {
      loading.value = false;
    }
  }

  // 初始化时获取购物车
  if (isAuthenticated.value) {
    fetchCart();
  }

  return {
    cart,
    loading,
    itemCount,
    total,
    fetchCart,
    addToCart,
    updateQuantity,
    removeFromCart,
    clearCart
  };
}