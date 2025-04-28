<template>
  <div>
    <!-- 抽屉开关按钮 -->
    <div class="drawer-toggle-button" @click="drawerOpen = true">
      <div class="handle-icon">
        <el-icon><Menu /></el-icon>
      </div>
    </div>

    <!-- 抽屉菜单 -->
    <el-drawer
      v-model="drawerOpen"
      direction="ltr"
      size="280px"
      :with-header="false"
      :destroy-on-close="false"
      custom-class="nav-drawer"
    >
      <div class="drawer-content">
        <div class="drawer-body px-4 pb-4 pt-6">
          <div class="mb-6">
            <router-link to="/" class="logo flex items-center mb-4" @click="drawerOpen = false">
              <el-icon class="mr-2 text-xl"><Mouse /></el-icon>
              <span class="font-bold text-xl text-white">{{ $t('common.app_name') }}</span>
            </router-link>

            <LanguageSwitcherFlags class="mt-4" />
          </div>

          <div class="mt-4">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">Main Navigation</h3>
            <div class="flex flex-col space-y-4">
              <router-link
                to="/"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><HomeFilled /></el-icon>
                <span>{{ $t('common.home') }}</span>
              </router-link>
              <router-link
                to="/database"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><DataAnalysis /></el-icon>
                <span>{{ $t('mouse.database') }}</span>
              </router-link>
              <router-link
                to="/reviews"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><ChatDotRound /></el-icon>
                <span>{{ $t('review.reviews') }}</span>
              </router-link>
              <router-link
                to="/compare"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><Switch /></el-icon>
                <span>{{ $t('comparison.title') }}</span>
              </router-link>
            </div>
          </div>

          <div class="mt-6">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">Tools</h3>
            <div class="flex flex-col space-y-4">
              <router-link
                to="/tools/sensitivity"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><Operation /></el-icon>
                <span>Sensitivity Tool</span>
              </router-link>
              <router-link
                to="/tools/ruler"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><ScaleToOriginal /></el-icon>
                <span>Measurement Ruler</span>
              </router-link>
            </div>
          </div>

          <div class="mt-6">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">Shopping</h3>
            <div class="flex flex-col space-y-4">
              <router-link
                to="/cart"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><ShoppingCart /></el-icon>
                <span>Cart</span>
              </router-link>
              <router-link
                to="/checkout"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><ShoppingBag /></el-icon>
                <span>Checkout</span>
              </router-link>
              <router-link
                v-if="isLoggedIn"
                to="/orders"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><Finished /></el-icon>
                <span>{{ $t('order.my_orders') }}</span>
              </router-link>
            </div>
          </div>

          <div class="mt-6" v-if="isLoggedIn">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">User Center</h3>
            <div class="flex flex-col space-y-4">
              <router-link
                to="/profile"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><User /></el-icon>
                <span>Profile</span>
              </router-link>
              <router-link
                to="/orders"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><Document /></el-icon>
                <span>Order History</span>
              </router-link>
              <div
                class="drawer-link cursor-pointer"
                @click="
                  logout();
                  drawerOpen = false;
                "
              >
                <el-icon><SwitchButton /></el-icon>
                <span>Logout</span>
              </div>
            </div>
          </div>

          <div class="mt-6" v-else>
            <div class="flex flex-col space-y-4">
              <router-link
                to="/login"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><Key /></el-icon>
                <span>{{ $t('common.login') }}</span>
              </router-link>
              <router-link
                to="/register"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><UserFilled /></el-icon>
                <span>{{ $t('common.register') }}</span>
              </router-link>
            </div>
          </div>

          <div class="mt-6">
            <h3 class="text-sm font-bold text-gray-400 mb-2 uppercase">About</h3>
            <div class="flex flex-col space-y-4">
              <router-link
                to="/about"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><InfoFilled /></el-icon>
                <span>About Us</span>
              </router-link>
              <router-link
                to="/contact"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><Message /></el-icon>
                <span>Contact Us</span>
              </router-link>
              <router-link
                to="/privacy"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><Lock /></el-icon>
                <span>Privacy Policy</span>
              </router-link>
              <router-link
                to="/terms"
                class="drawer-link"
                active-class="active"
                @click="drawerOpen = false"
              >
                <el-icon><Document /></el-icon>
                <span>Terms of Use</span>
              </router-link>
            </div>
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
  Mouse,
  Menu,
  Close,
  Plus,
  HomeFilled,
  DataAnalysis,
  ChatDotRound,
  Switch,
  Reading,
  ScaleToOriginal,
  Operation,
  Pointer,
  User,
  Monitor,
  Document,
  SwitchButton,
  ShoppingCart,
  ShoppingBag,
  Finished,
  InfoFilled,
  Message,
  Lock,
  Key,
  UserFilled
} from '@element-plus/icons-vue';
import LanguageSwitcherFlags from '@/components/common/LanguageSwitcherFlags.vue';
import { useI18n } from 'vue-i18n';

// Initialize i18n
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
.drawer-toggle-button {
  position: fixed;
  top: 50%;
  left: 0;
  transform: translateY(-50%);
  background: var(--claude-primary-purple);
  color: white;
  padding: 12px 8px;
  border-radius: 0 8px 8px 0;
  cursor: pointer;
  z-index: 40;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.2);
  transition: transform 0.3s ease;
}

.drawer-toggle-button:hover {
  transform: translateY(-50%) translateX(5px);
  background: var(--claude-primary-purple-darker);
}

.drawer-toggle-button .handle-icon {
  display: flex;
  align-items: center;
  justify-content: center;
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

.logo {
  color: white;
  text-decoration: none;
}
</style>
