<template>
  <div class="view-toggle bg-[#1A1A1A] border border-[#333333] rounded-lg p-3 flex items-center justify-between">
    <span class="text-sm text-white font-medium">显示方式:</span>
    
    <div class="flex gap-2">
      <el-button
        :type="viewMode === 'grid' ? 'primary' : 'default'"
        :icon="Grid"
        size="small"
        @click="changeViewMode('grid')"
        class="dark-button"
      >
        网格视图
      </el-button>
      
      <el-button
        :type="viewMode === 'list' ? 'primary' : 'default'"
        :icon="List"
        size="small"
        @click="changeViewMode('list')"
        class="dark-button"
      >
        列表视图
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Grid, List } from '@element-plus/icons-vue'

// 视图模式
const viewMode = ref('grid')

// 定义事件
const emit = defineEmits(['change-view'])

// 改变视图模式
const changeViewMode = (mode: string) => {
  viewMode.value = mode
  emit('change-view', mode)
}
</script>

<style scoped>
.dark-button {
  background-color: var(--claude-bg-light);
  border-color: var(--claude-border-light);
  color: var(--claude-text-light);
}

.dark-button:hover {
  background-color: var(--claude-bg-medium);
  border-color: var(--claude-border-dark);
}

.el-button.is-active {
  background-color: var(--claude-primary-purple);
  border-color: var(--claude-primary-purple);
  color: white;
}
</style>