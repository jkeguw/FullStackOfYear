<template>
  <el-button
    :type="buttonType"
    :size="buttonSize"
    :disabled="disabled || loading"
    :loading="loading"
    :icon="buttonIcon"
    @click="handleAddToCart"
  >
    <slot>{{ buttonText }}</slot>
  </el-button>
</template>

<script setup lang="ts">
import { ref, defineProps, computed } from 'vue';
import { ElMessage } from 'element-plus';
import { ShoppingCart, Plus } from '@element-plus/icons-vue';
import { useCart } from '@/composables/useCart';
import { useAuth } from '@/composables/useAuth';
import type { CartItem } from '@/types/cart';

const props = defineProps({
  // 商品信息
  product: {
    type: Object,
    required: true,
    validator: (value: any) => {
      return value && value.id && value.name && value.price !== undefined;
    }
  },
  // 商品类型
  productType: {
    type: String,
    default: 'mouse'
  },
  // 按钮类型
  type: {
    type: String,
    default: 'primary'
  },
  // 按钮大小
  size: {
    type: String,
    default: 'default'
  },
  // 是否只显示图标
  iconOnly: {
    type: Boolean,
    default: false
  },
  // 自定义按钮文本
  text: {
    type: String,
    default: '加入购物车'
  },
  // 商品数量
  quantity: {
    type: Number,
    default: 1
  },
  // 是否禁用
  disabled: {
    type: Boolean,
    default: false
  }
});

const loading = ref(false);
const { addToCart } = useCart();
const { isAuthenticated } = useAuth();

// 计算属性
const buttonType = computed(() => props.type);
const buttonSize = computed(() => props.size);
const buttonText = computed(() => props.iconOnly ? '' : props.text);
const buttonIcon = computed(() => props.iconOnly ? ShoppingCart : Plus);

// 添加到购物车
async function handleAddToCart() {
  if (!isAuthenticated.value) {
    ElMessage.warning('请先登录');
    return;
  }

  loading.value = true;
  
  try {
    const cartItem: CartItem = {
      product_id: props.product.id,
      product_type: props.productType,
      name: props.product.name,
      price: props.product.price,
      quantity: props.quantity,
      image_url: props.product.image || props.product.image_url
    };
    
    await addToCart(cartItem);
  } catch (error) {
    console.error('添加到购物车失败', error);
  } finally {
    loading.value = false;
  }
}
</script>