<template>
  <div
    class="sort-controls bg-[#1A1A1A] border border-[#333333] rounded-lg p-3 flex items-center justify-between"
  >
    <span class="text-sm text-white font-medium">排序方式:</span>

    <div class="flex items-center gap-2">
      <el-select
        v-model="sortBy"
        placeholder="选择排序字段"
        size="small"
        class="dark-select w-32"
        @change="changeSorting"
      >
        <el-option
          v-for="option in sortOptions"
          :key="option.value"
          :label="option.label"
          :value="option.value"
        />
      </el-select>

      <el-button
        :icon="sortOrder === 'asc' ? SortUp : SortDown"
        circle
        size="small"
        @click="toggleSortOrder"
        class="dark-button"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { SortUp, SortDown } from '@element-plus/icons-vue';

// 排序选项 (移除了最新上架选项)
const sortOptions = [
  { label: '名称', value: 'name' },
  { label: '品牌', value: 'brand' },
  { label: '重量', value: 'weight' },
  { label: '价格', value: 'price' }
];

// 排序状态
const sortBy = ref('name');
const sortOrder = ref('asc');

// 定义事件
const emit = defineEmits(['sort-change']);

// 改变排序字段
const changeSorting = () => {
  emit('sort-change', { field: sortBy.value, order: sortOrder.value });
};

// 切换排序顺序
const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
  changeSorting();
};
</script>

<style scoped>
.dark-select :deep(.el-input__inner) {
  background-color: var(--claude-bg-light);
  border-color: var(--claude-border-light);
  color: var(--claude-text-light);
}

.dark-button {
  background-color: var(--claude-bg-light);
  border-color: var(--claude-border-light);
  color: var(--claude-text-light);
}

.dark-button:hover {
  background-color: var(--claude-bg-medium);
  border-color: var(--claude-border-dark);
}
</style>
