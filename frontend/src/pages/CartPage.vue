<template>
  <div class="cart-page">
    <h1 class="page-title">我的购物车</h1>

    <div class="cart-container">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="8" animated />
      </div>
      <template v-else>
        <div v-if="cart.items.length === 0" class="empty-cart">
          <el-empty description="购物车为空" />
          <el-button type="primary" @click="goToShop">浏览商品</el-button>
        </div>
        <template v-else>
          <div class="cart-table">
            <el-table :data="cart.items" style="width: 100%">
              <el-table-column label="商品信息" min-width="300">
                <template #default="{ row }">
                  <div class="product-info">
                    <div class="product-image">
                      <img :src="row.imageUrl || '/placeholder.png'" :alt="row.name" />
                    </div>
                    <div class="product-name">{{ row.name }}</div>
                  </div>
                </template>
              </el-table-column>

              <el-table-column label="单价" width="120">
                <template #default="{ row }">
                  <div class="product-price">¥{{ row.price.toFixed(2) }}</div>
                </template>
              </el-table-column>

              <el-table-column label="数量" width="160">
                <template #default="{ row }">
                  <el-input-number
                    v-model="row.quantity"
                    :min="1"
                    :max="99"
                    size="small"
                    @change="(val) => handleQuantityChange(row.productId, val)"
                  />
                </template>
              </el-table-column>

              <el-table-column label="小计" width="120">
                <template #default="{ row }">
                  <div class="product-subtotal">¥{{ (row.price * row.quantity).toFixed(2) }}</div>
                </template>
              </el-table-column>

              <el-table-column label="操作" width="100">
                <template #default="{ row }">
                  <el-button type="danger" size="small" @click="removeItem(row.productId)">
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <div class="cart-summary">
            <div class="cart-actions">
              <el-button @click="goToShop">继续购物</el-button>
              <el-button type="danger" @click="handleClearCart">清空购物车</el-button>
            </div>

            <div class="cart-checkout">
              <div class="cart-total">
                <span>商品总计 ({{ cart.itemCount }} 件商品):</span>
                <span class="price">¥{{ total.toFixed(2) }}</span>
              </div>
              <el-button type="primary" size="large" @click="goToCheckout"> 去结算 </el-button>
            </div>
          </div>
        </template>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessageBox } from 'element-plus';
import { useCart } from '@/composables/useCart';

const router = useRouter();
const { cart, loading, total, fetchCart, updateQuantity, removeFromCart, clearCart } = useCart();

// 页面加载时获取购物车信息
onMounted(() => {
  fetchCart();
});

// 处理数量变更
async function handleQuantityChange(productId: string, quantity: number) {
  await updateQuantity(productId, quantity);
}

// 移除商品
async function removeItem(productId: string) {
  ElMessageBox.confirm('确定要从购物车中移除此商品吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      await removeFromCart(productId);
    })
    .catch(() => {});
}

// 清空购物车
async function handleClearCart() {
  ElMessageBox.confirm('确定要清空购物车吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      await clearCart();
    })
    .catch(() => {});
}

// 前往商店
function goToShop() {
  router.push('/devices');
}

// 前往结算页面
function goToCheckout() {
  router.push('/checkout');
}
</script>

<style scoped>
.cart-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: #333;
}

.cart-container {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.loading-container {
  padding: 40px 20px;
}

.empty-cart {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  gap: 20px;
}

.product-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.product-image {
  width: 80px;
  height: 80px;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid #f0f0f0;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-name {
  font-weight: 500;
}

.product-price,
.product-subtotal {
  font-weight: 600;
  color: #ff6700;
}

.cart-summary {
  margin-top: 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.cart-actions {
  display: flex;
  gap: 10px;
}

.cart-checkout {
  display: flex;
  align-items: center;
  gap: 20px;
}

.cart-total {
  font-size: 16px;
}

.price {
  color: #ff6700;
  font-weight: 600;
  font-size: 24px;
  margin-left: 10px;
}
</style>
