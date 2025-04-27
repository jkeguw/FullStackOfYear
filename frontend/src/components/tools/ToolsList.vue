<template>
  <div class="tools-list">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div 
        v-for="tool in tools" 
        :key="tool.id" 
        class="tool-card bg-gray-800 rounded-lg overflow-hidden hover:shadow-xl transition duration-300"
      >
        <div class="p-6">
          <div class="tool-icon mb-4 flex justify-center">
            <el-icon :size="40" class="text-indigo-400">
              <component :is="tool.icon" />
            </el-icon>
          </div>
          <h3 class="text-xl font-bold mb-3 text-center">{{ tool.name }}</h3>
          <p class="text-gray-400 mb-4 text-center">{{ tool.description }}</p>
          <div class="text-center">
            <router-link :to="tool.route" class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition">
              使用工具
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { 
  ScaleToOriginal, 
  Monitor
} from '@element-plus/icons-vue';

// 工具列表
const tools = [
  {
    id: 'dpi',
    name: 'DPI 计算器',
    description: '在不同DPI设置之间转换灵敏度，确保一致的手感',
    icon: Monitor,
    route: '/tools/dpi'
  },
  {
    id: 'ruler',
    name: '厘米尺',
    description: '使用可调节长度的厘米尺进行简单测量',
    icon: ScaleToOriginal,
    route: '/tools/ruler'
  }
];
</script>

<style scoped>
.tool-card {
  transition: transform 0.3s ease;
  border: 1px solid rgba(75, 85, 99, 0.4);
}

.tool-card:hover {
  transform: translateY(-5px);
  border-color: rgba(124, 58, 237, 0.5);
}
</style>