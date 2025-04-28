<template>
  <div class="test-mouse-comparison">
    <h1 class="text-xl font-bold mb-4">鼠标比较测试页面</h1>
    <p class="mb-4">使用硬编码SVG数据进行鼠标对比测试</p>
    
    <div class="mb-6">
      <el-alert
        type="info"
        :closable="false"
        show-icon
      >
        <div>选择以下鼠标进行对比：</div>
        <div class="mt-2">
          <el-checkbox-group v-model="selectedMiceIds">
            <el-checkbox v-for="mouse in mockMice" :key="mouse.id" :label="mouse.id">
              {{ mouse.brand }} {{ mouse.name }}
            </el-checkbox>
          </el-checkbox-group>
        </div>
      </el-alert>
    </div>
    
    <div v-if="selectedMiceIds.length > 0" class="mb-6">
      <el-button type="primary" @click="addSelectedMice">
        添加选中的鼠标到比较视图
      </el-button>
      <el-button @click="clearComparison">
        清空比较视图
      </el-button>
    </div>
    
    <!-- 鼠标比较视图组件 -->
    <MouseComparisonView />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useComparisonStore } from '@/stores';
import MouseComparisonView from './MouseComparisonView.vue';
import { mockMice } from '@/data/mockMice';

const comparisonStore = useComparisonStore();
const selectedMiceIds = ref<string[]>([]);

// 添加选择的鼠标到比较视图
function addSelectedMice() {
  // 先清空当前选择
  clearComparison();
  
  // 添加新选择的鼠标
  for (const id of selectedMiceIds.value) {
    const mouse = mockMice.find(m => m.id === id);
    if (mouse) {
      comparisonStore.addMouse(mouse);
    }
  }
}

// 清空比较视图
function clearComparison() {
  comparisonStore.clearMice();
}

// 初始默认选择第一个鼠标
onMounted(() => {
  if (mockMice.length > 0) {
    selectedMiceIds.value = [mockMice[0].id];
    addSelectedMice();
  }
});
</script>

<style scoped>
.test-mouse-comparison {
  padding: 1rem;
  max-width: 1200px;
  margin: 0 auto;
}
</style>