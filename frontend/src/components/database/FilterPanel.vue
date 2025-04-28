<template>
  <div class="filter-panel bg-[#1A1A1A] border border-[#333333] rounded-lg p-4">
    <h3 class="text-lg font-medium text-white mb-4">Filter Criteria</h3>

    <div class="space-y-6">
      <!-- Search box -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">Search Name</label>
        <el-input
          v-model="searchQuery"
          placeholder="Enter mouse name..."
          :prefix-icon="Search"
          class="search-input"
        />
      </div>

      <!-- Brand filter (dropdown) -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">Brand</label>
        <el-select 
          v-model="selectedBrands" 
          multiple 
          clearable 
          collapse-tags
          class="w-full filter-select" 
          placeholder="Select brands">
          <el-option
            v-for="brand in brands"
            :key="brand"
            :label="brand"
            :value="brand"
          />
        </el-select>
      </div>

      <!-- Shape filter (dropdown) -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">Shape</label>
        <el-select 
          v-model="selectedShapes" 
          multiple 
          clearable 
          collapse-tags
          class="w-full filter-select" 
          placeholder="Select shapes">
          <el-option
            v-for="shape in shapes"
            :key="shape"
            :label="shapeLabels[shape] || shape"
            :value="shape"
          />
        </el-select>
      </div>

      <!-- Weight range -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">
          Weight Range: {{ weightRange[0] }}g - {{ weightRange[1] }}g
        </label>
        <el-slider v-model="weightRange" range :min="40" :max="150" class="filter-slider" />
      </div>

      <!-- Connection type (dropdown) -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">Connection Type</label>
        <el-select 
          v-model="selectedConnections" 
          multiple 
          clearable 
          collapse-tags
          class="w-full filter-select" 
          placeholder="Select connection types">
          <el-option
            v-for="conn in connections"
            :key="conn"
            :label="connectionLabels[conn] || conn"
            :value="conn"
          />
        </el-select>
      </div>

      <!-- Button count range -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">
          Button Count: {{ buttonRange[0] }} - {{ buttonRange[1] }}
        </label>
        <el-slider v-model="buttonRange" range :min="2" :max="20" :step="1" class="filter-slider" />
      </div>

      <!-- Action buttons -->
      <div class="flex justify-between mt-6">
        <el-button @click="resetFilters" :icon="RefreshRight" class="dark-button">
          Reset Filters
        </el-button>
        <el-button type="primary" @click="applyFilters" :icon="Filter"> Apply Filters </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import { Search, RefreshRight, Filter } from '@element-plus/icons-vue';

// 品牌列表
const brands = [
  'Logitech',
  'Razer',
  'SteelSeries',
  'Zowie',
  'Glorious',
  'Vaxee',
  'Pulsar',
  'Cooler Master',
  'HyperX',
  'Endgame Gear'
];

// 鼠标形状
const shapes = ['ergo', 'ambi', 'symmetrical'];
const shapeLabels = {
  ergo: 'Ergonomic',
  ambi: 'Ambidextrous',
  symmetrical: 'Symmetrical'
};

// Connection types
const connections = ['wired', 'wireless', 'hybrid'];
const connectionLabels = {
  wired: 'Wired',
  wireless: 'Wireless',
  hybrid: 'Hybrid'
};

// 筛选状态
const searchQuery = ref('');
const selectedBrands = ref([]);
const selectedShapes = ref([]);
const weightRange = ref([40, 150]);
const defaultWeightRange = [40, 150];
const selectedConnections = ref([]);
const buttonRange = ref([2, 20]);
const defaultButtonRange = [2, 20];

// 重置筛选
const resetFilters = () => {
  searchQuery.value = '';
  selectedBrands.value = [];
  selectedShapes.value = [];
  weightRange.value = [...defaultWeightRange];
  selectedConnections.value = [];
  buttonRange.value = [...defaultButtonRange];

  // 发出事件通知父组件筛选已重置
  emits('reset-filters');
};

// 应用筛选
const applyFilters = () => {
  const filters = {
    query: searchQuery.value,
    brands: selectedBrands.value,
    shapes: selectedShapes.value,
    weight: weightRange.value,
    connections: selectedConnections.value,
    buttons: buttonRange.value
  };

  // 发出筛选事件
  emits('apply-filters', filters);
};

// 定义组件事件
const emits = defineEmits(['apply-filters', 'reset-filters']);

// 监听筛选条件变化，实时更新
watch(
  [searchQuery, selectedBrands, selectedShapes, weightRange, selectedConnections, buttonRange],
  () => {
    // 自动应用筛选，可以根据需要取消注释
    // applyFilters()
  }
);
</script>

<style scoped>
.filter-panel {
  background-color: var(--claude-bg-medium);
}

.search-input :deep(.el-input__inner) {
  background-color: var(--claude-bg-light);
  border-color: var(--claude-border-light);
  color: var(--claude-text-light);
}

.filter-select {
  width: 100%;
}

.filter-select :deep(.el-input__inner) {
  background-color: var(--claude-bg-light);
  border-color: var(--claude-border-light);
  color: var(--claude-text-light);
}

.filter-select :deep(.el-select__tags) {
  background-color: var(--claude-bg-light);
}

.filter-select :deep(.el-tag) {
  background-color: var(--claude-primary-purple);
  border-color: var(--claude-primary-purple);
  color: white;
}

.filter-slider :deep(.el-slider__runway) {
  background-color: var(--claude-border-light);
}

.filter-slider :deep(.el-slider__bar) {
  background-color: var(--claude-primary-purple);
}

.filter-slider :deep(.el-slider__button) {
  border-color: var(--claude-primary-purple);
  background-color: var(--claude-primary-purple);
}

.dark-button {
  color: var(--claude-text-light);
  background-color: var(--claude-bg-light);
  border-color: var(--claude-border-light);
}

.dark-button:hover {
  background-color: var(--claude-bg-medium);
  border-color: var(--claude-border-dark);
}
</style>
