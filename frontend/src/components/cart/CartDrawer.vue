<template>
  <el-drawer
    v-model="drawerVisible"
    title="购物车"
    direction="rtl"
    size="350px"
    :before-close="onClose"
  >
    <div class="cart-container">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>
      <template v-else>
        <div v-if="cart.items.length === 0" class="empty-cart">
          <el-empty description="购物车为空" />
          <el-button type="primary" @click="goToShop">浏览商品</el-button>
        </div>
        <template v-else>
          <div class="cart-list">
            <div v-for="item in cart.items" :key="item.product_id" class="cart-item">
              <div class="cart-item-header">
                <div class="cart-item-image">
                  <img 
                    :src="item.image_url || '/placeholder.png'" 
                    :alt="item.name"
                  />
                </div>
                <div class="cart-item-info">
                  <div class="cart-item-name">{{ item.name }}</div>
                  <div class="cart-item-price">¥{{ item.price.toFixed(2) }}</div>
                </div>
              </div>
              <div class="cart-item-actions">
                <el-input-number 
                  v-model="item.quantity" 
                  :min="1" 
                  :max="99"
                  size="small"
                  @change="(val) => handleQuantityChange(item.product_id, val)"
                />
                <el-button 
                  type="danger" 
                  size="small" 
                  icon="Delete" 
                  circle
                  @click="removeItem(item.product_id)"
                />
              </div>
            </div>
          </div>
          <div class="cart-footer">
            <div class="cart-total">
              <span>合计:</span>
              <span class="price">¥{{ total.toFixed(2) }}</span>
            </div>
            <div class="cart-actions">
              <el-button type="danger" @click="handleClearCart">清空购物车</el-button>
              <el-button type="primary" @click="goToCheckout">结算</el-button>
            </div>
          </div>
        </template>
      </template>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, defineEmits, defineProps, watch } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessageBox } from 'element-plus';
import { useCart } from '@/composables/useCart';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:visible', 'close']);
const router = useRouter();
const drawerVisible = ref(props.visible);

const { cart, loading, total, updateQuantity, removeFromCart, clearCart } = useCart();

// 同步drawer显示状态
watch(() => props.visible, (val) => {
  drawerVisible.value = val;
});

watch(() => drawerVisible.value, (val) => {
  emit('update:visible', val);
  if (!val) {
    emit('close');
  }
});

// 处理数量变更
async function handleQuantityChange(productId: string, quantity: number) {
  await updateQuantity(productId, quantity);
}

// 移除商品
async function removeItem(productId: string) {
  await removeFromCart(productId);
}

// 清空购物车
async function handleClearCart() {
  ElMessageBox.confirm('确定要清空购物车吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    await clearCart();
  }).catch(() => {});
}

// 前往商店
function goToShop() {
  router.push('/mouse-database');
  drawerVisible.value = false;
}

// 前往结算页面
function goToCheckout() {
  router.push('/checkout');
  drawerVisible.value = false;
}

// 关闭抽屉
function onClose() {
  drawerVisible.value = false;
}
</script>

<style scoped>
.cart-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.loading-container {
  padding: 20px;
}

.empty-cart {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 20px;
}

.cart-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.cart-item {
  border-bottom: 1px solid #f0f0f0;
  padding: 15px 0;
}

.cart-item-header {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.cart-item-image {
  width: 70px;
  height: 70px;
  overflow: hidden;
  border-radius: 4px;
  border: 1px solid #f0f0f0;
}

.cart-item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cart-item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.cart-item-name {
  font-weight: 500;
  font-size: 14px;
}

.cart-item-price {
  color: #ff6700;
  font-weight: 600;
  font-size: 16px;
}

.cart-item-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cart-footer {
  margin-top: auto;
  padding: 15px;
  border-top: 1px solid #f0f0f0;
}

.cart-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-size: 16px;
}

.price {
  color: #ff6700;
  font-weight: 600;
  font-size: 20px;
}

.cart-actions {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}
</style>