<template>
  <div>
    <nav class="navbar">
      <div class="container mx-auto px-4">
        <div class="flex justify-between items-center py-4">
          <router-link to="/" class="logo flex items-center">
            <el-icon class="mr-2 text-xl"><Mouse /></el-icon>
            <span class="font-bold text-xl text-white">{{ $t('common.app_name') }}</span>
          </router-link>
          
          <div class="flex items-center space-x-4">
            <el-button v-if="isInCompare" type="primary" size="small" @click="goToComparison">
              <el-badge :value="comparisonCount" :hidden="comparisonCount === 0">
                比较中
              </el-badge>
            </el-button>
            
            <router-link v-if="!isLoggedIn" to="/login">
              <el-button type="primary" size="small">{{ $t('common.login') }}</el-button>
            </router-link>
            
            <el-dropdown v-else>
              <el-avatar 
                :size="32" 
                :src="userAvatar" 
                class="cursor-pointer"
              >
                {{ userInitials }}
              </el-avatar>
              
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="$router.push('/profile')">个人资料</el-dropdown-item>
                  <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            
            <language-switcher class="hidden md:block" />
            
            <el-button @click="drawerOpen = true" type="primary" circle>
              <el-icon><Menu /></el-icon>
            </el-button>
          </div>
        </div>
      </div>
    </nav>
    
    <!-- 抽屉菜单 -->
    <el-drawer
      v-model="drawerOpen"
      title="菜单导航"
      direction="rtl"
      size="280px"
      :with-header="false"
      :destroy-on-close="false"
      custom-class="nav-drawer"
    >
      <div class="drawer-content">
        <div class="drawer-header">
          <div class="flex items-center justify-between px-4 py-4 border-b border-gray-700">
            <span class="text-lg font-bold text-white">导航菜单</span>
            <el-button @click="drawerOpen = false" type="text" circle>
              <el-icon><Close /></el-icon>
            </el-button>
          </div>
        </div>
        
        <div class="drawer-body px-4 pb-4">
          <div class="mt-4">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">主导航</h3>
            <div class="flex flex-col space-y-4">
              <router-link to="/" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><HomeFilled /></el-icon>
                <span>{{ $t('common.home') }}</span>
              </router-link>
              <router-link to="/database" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><DataAnalysis /></el-icon>
                <span>{{ $t('mouse.database') }}</span>
              </router-link>
              <router-link to="/reviews" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><ChatDotRound /></el-icon>
                <span>{{ $t('review.reviews') }}</span>
              </router-link>
              <router-link to="/compare" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><Switch /></el-icon>
                <span>{{ $t('comparison.title') }}</span>
              </router-link>
            </div>
          </div>
          
          <div class="mt-6">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">工具集</h3>
            <div class="flex flex-col space-y-4">
              <router-link to="/tools/dpi" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><Reading /></el-icon>
                <span>DPI 计算器</span>
              </router-link>
              <router-link to="/tools/ruler" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><ScaleToOriginal /></el-icon>
                <span>测量尺子</span>
              </router-link>
              <router-link to="/tools/sensitivity" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><Operation /></el-icon>
                <span>灵敏度转换</span>
              </router-link>
            </div>
          </div>
          
          <div class="mt-6">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">购物</h3>
            <div class="flex flex-col space-y-4">
              <router-link to="/cart" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><ShoppingCart /></el-icon>
                <span>{{ $t('cart.shopping_cart') }}</span>
              </router-link>
              <router-link to="/checkout" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><ShoppingBag /></el-icon>
                <span>{{ $t('order.checkout') }}</span>
              </router-link>
              <router-link v-if="isLoggedIn" to="/orders" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><Finished /></el-icon>
                <span>{{ $t('order.my_orders') }}</span>
              </router-link>
            </div>
          </div>
          
          <div class="mt-6" v-if="isLoggedIn">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">用户中心</h3>
            <div class="flex flex-col space-y-4">
              <router-link to="/profile" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><User /></el-icon>
                <span>个人资料</span>
              </router-link>
              <!-- 已移除个人设备管理功能 -->
              <router-link to="/orders" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><Document /></el-icon>
                <span>订单记录</span>
              </router-link>
              <div class="drawer-link cursor-pointer" @click="logout(); drawerOpen = false">
                <el-icon><SwitchButton /></el-icon>
                <span>退出登录</span>
              </div>
            </div>
          </div>
          
          <div class="mt-6">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">{{ $t('common.about') }}</h3>
            <div class="flex flex-col space-y-4">
              <router-link to="/about" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><InfoFilled /></el-icon>
                <span>{{ $t('common.about_us') }}</span>
              </router-link>
              <router-link to="/contact" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><Message /></el-icon>
                <span>{{ $t('common.contact') }}</span>
              </router-link>
              <router-link to="/privacy" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><Lock /></el-icon>
                <span>{{ $t('common.privacy') }}</span>
              </router-link>
              <router-link to="/terms" class="drawer-link" active-class="active" @click="drawerOpen = false">
                <el-icon><Document /></el-icon>
                <span>{{ $t('common.terms') }}</span>
              </router-link>
            </div>
          </div>
          
          <div class="mt-6">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">{{ $t('common.language') }}</h3>
            <language-switcher />
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore, useComparisonStore } from '@/stores';
import { 
  Mouse, Menu, Close, Plus, HomeFilled, DataAnalysis, 
  ChatDotRound, Switch, Reading, ScaleToOriginal, Operation, 
  Pointer, User, Monitor, Document, SwitchButton, ShoppingCart, 
  ShoppingBag, Finished, InfoFilled, Message, Lock
} from '@element-plus/icons-vue';
import LanguageSwitcher from '@/components/common/LanguageSwitcher.vue';
import { useI18n } from 'vue-i18n';

// 初始化i18n
const { t } = useI18n();

// 路由
const router = useRouter();

// Store
const userStore = useUserStore();
const comparisonStore = useComparisonStore();

// 状态
const drawerOpen = ref(false);

// 计算属性
const isLoggedIn = computed(() => !!userStore.token);

const userAvatar = computed(() => {
  return userStore.user?.avatar || '';
});

const userInitials = computed(() => {
  if (!userStore.user?.name) return '';
  return userStore.user.name.substring(0, 1).toUpperCase();
});

const comparisonCount = computed(() => {
  return comparisonStore.selectedMice.length;
});

const isInCompare = computed(() => {
  return comparisonCount.value > 0;
});

// 方法
function logout() {
  userStore.clearUser();
  router.push('/');
}

function goToComparison() {
  router.push('/compare');
}
</script>

<style scoped>
.navbar {
  background-color: var(--claude-bg-dark);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  position: sticky;
  top: 0;
  z-index: 50;
  border-bottom: 1px solid var(--claude-border-dark);
}

.nav-drawer :deep(.el-drawer__body) {
  padding: 0;
  background-color: var(--claude-bg-dark);
  color: white;
}

.drawer-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.drawer-body {
  flex: 1;
  overflow-y: auto;
}

.drawer-link {
  display: flex;
  align-items: center;
  color: var(--claude-text-light);
  font-weight: 500;
  padding: 0.75rem 0.5rem;
  border-radius: 0.375rem;
  transition: all 0.2s ease;
}

.drawer-link:hover {
  color: var(--claude-text-white);
  background-color: rgba(255, 255, 255, 0.05);
}

.drawer-link.active {
  color: var(--claude-primary-purple);
  font-weight: 600;
  background-color: rgba(125, 90, 243, 0.1);
}

.drawer-link .el-icon {
  margin-right: 12px;
  font-size: 18px;
}
</style>