<template>
  <div class="cart-icon-container">
    <el-badge :value="itemCount" :hidden="itemCount <= 0">
      <el-button circle @click="openCart">
        <el-icon><ShoppingCart /></el-icon>
      </el-button>
    </el-badge>

    <CartDrawer v-model:visible="cartVisible" @close="onCartClose" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ShoppingCart } from '@element-plus/icons-vue';
import { useCart } from '@/composables/useCart';
import CartDrawer from './CartDrawer.vue';

const cartVisible = ref(false);
const { itemCount, fetchCart } = useCart();

function openCart() {
  // 打开购物车前刷新购物车数据
  fetchCart();
  cartVisible.value = true;
}

function onCartClose() {
  cartVisible.value = false;
}
</script>

<style scoped>
.cart-icon-container {
  display: inline-block;
  cursor: pointer;
}
</style>
