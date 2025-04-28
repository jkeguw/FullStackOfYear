<template>
  <div class="tools-page bg-gray-900 text-white min-h-screen">
    <div class="container mx-auto px-4 py-8">
      <h1 class="text-3xl font-bold mb-8 text-center">游戏外设工具集</h1>

      <div class="tool-content bg-gray-800 rounded-lg p-4 sm:p-6">
        <div v-if="hasToolRoute" class="mb-4">
          <router-link to="/tools" class="text-indigo-400 hover:text-indigo-300 flex items-center">
            <el-icon class="mr-1"><Back /></el-icon>
            返回工具集
          </router-link>
        </div>

        <!-- 工具主页 -->
        <div v-if="!hasToolRoute">
          <div class="tools-grid grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <!-- 尺子工具已废弃 -->
            <div
              class="tool-card p-4 bg-gray-700 rounded-lg shadow-lg border border-gray-600 opacity-50"
            >
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-xl font-semibold">尺子工具</h3>
                <el-tag type="danger" size="small">已废弃</el-tag>
              </div>
              <p class="text-gray-300 mb-4">此功能已被移除，请使用其他工具</p>
              <div class="flex justify-between text-sm text-gray-400">
                <span>功能已弃用</span>
                <span>不再可用</span>
              </div>
            </div>

            <!-- DPI计算器 -->
            <div
              class="tool-card p-4 bg-gray-700 rounded-lg shadow-lg border border-gray-600 hover:border-gray-500 transition-all"
            >
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-xl font-semibold">DPI计算器</h3>
                <router-link to="/tools/dpi" class="el-button el-button--primary el-button--small"
                  >启动</router-link
                >
              </div>
              <p class="text-gray-300 mb-4">
                计算鼠标DPI与实际移动距离的关系，帮助找到最适合的灵敏度
              </p>
              <div class="flex justify-between text-sm text-gray-400">
                <span>游戏必备</span>
                <span>精确调整</span>
              </div>
            </div>

            <!-- 灵敏度转换 -->
            <div
              class="tool-card p-4 bg-gray-700 rounded-lg shadow-lg border border-gray-600 hover:border-gray-500 transition-all"
            >
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-xl font-semibold">灵敏度转换</h3>
                <router-link
                  to="/tools/sensitivity"
                  class="el-button el-button--primary el-button--small"
                  >启动</router-link
                >
              </div>
              <p class="text-gray-300 mb-4">
                在不同游戏和应用程序之间转换鼠标灵敏度，保持相同的手感
              </p>
              <div class="flex justify-between text-sm text-gray-400">
                <span>多游戏支持</span>
                <span>手感一致</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 工具内容（子路由) -->
        <router-view v-if="hasToolRoute" />
      </div>
    </div>

    <!-- 移除浮动尺子组件 -->
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';
import { Back } from '@element-plus/icons-vue';

const route = useRoute();
// 尺子工具已废弃

// 检查是否在工具子路由
const hasToolRoute = computed(() => {
  return route.path !== '/tools';
});
</script>

<style scoped>
.tools-page {
  padding-bottom: 2rem;
}

.tool-card {
  transition: all 0.3s ease;
}

.tool-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3);
}
</style>
