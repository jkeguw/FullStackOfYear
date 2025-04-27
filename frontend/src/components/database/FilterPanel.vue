<template>
  <div class="filter-panel bg-[#1A1A1A] border border-[#333333] rounded-lg p-4">
    <h3 class="text-lg font-medium text-white mb-4">筛选条件</h3>
    
    <div class="space-y-6">
      <!-- 搜索框 -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">搜索名称</label>
        <el-input 
          v-model="searchQuery" 
          placeholder="输入鼠标名称..." 
          :prefix-icon="Search"
          class="search-input"
        />
      </div>
      
      <!-- 品牌筛选 -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">品牌</label>
        <el-checkbox-group v-model="selectedBrands" class="filter-checkbox-group">
          <el-checkbox 
            v-for="brand in brands" 
            :key="brand" 
            :label="brand"
            class="filter-checkbox text-gray-300"
          >
            {{ brand }}
          </el-checkbox>
        </el-checkbox-group>
      </div>
      
      <!-- 形状筛选 -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">形状</label>
        <el-checkbox-group v-model="selectedShapes" class="filter-checkbox-group">
          <el-checkbox 
            v-for="shape in shapes" 
            :key="shape" 
            :label="shape"
            class="filter-checkbox text-gray-300"
          >
            {{ shapeLabels[shape] || shape }}
          </el-checkbox>
        </el-checkbox-group>
      </div>
      
      <!-- 重量范围 -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">
          重量范围: {{ weightRange[0] }}g - {{ weightRange[1] }}g
        </label>
        <el-slider
          v-model="weightRange"
          range
          :min="40"
          :max="150"
          class="filter-slider"
        />
      </div>
      
      <!-- 连接方式 -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">连接方式</label>
        <el-checkbox-group v-model="selectedConnections" class="filter-checkbox-group">
          <el-checkbox 
            v-for="conn in connections" 
            :key="conn" 
            :label="conn"
            class="filter-checkbox text-gray-300"
          >
            {{ connectionLabels[conn] || conn }}
          </el-checkbox>
        </el-checkbox-group>
      </div>
      
      <!-- 按钮数量范围 -->
      <div>
        <label class="block text-gray-400 text-sm font-medium mb-2">
          按钮数量: {{ buttonRange[0] }} - {{ buttonRange[1] }}
        </label>
        <el-slider
          v-model="buttonRange"
          range
          :min="2"
          :max="20"
          :step="1"
          class="filter-slider"
        />
      </div>
      
      <!-- 操作按钮 -->
      <div class="flex justify-between mt-6">
        <el-button 
          @click="resetFilters" 
          :icon="RefreshRight"
          class="dark-button"
        >
          重置筛选
        </el-button>
        <el-button 
          type="primary" 
          @click="applyFilters"
          :icon="Filter"
        >
          应用筛选
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { Search, RefreshRight, Filter } from '@element-plus/icons-vue'

// 品牌列表
const brands = [
  'Logitech', 'Razer', 'SteelSeries', 'Zowie', 'Glorious', 
  'Vaxee', 'Pulsar', 'Cooler Master', 'HyperX', 'Endgame Gear'
]

// 鼠标形状
const shapes = ['ergo', 'ambi', 'symmetrical']
const shapeLabels = {
  'ergo': '人体工学',
  'ambi': '双手通用',
  'symmetrical': '对称'
}

// 连接方式
const connections = ['wired', 'wireless', 'hybrid']
const connectionLabels = {
  'wired': '有线',
  'wireless': '无线',
  'hybrid': '双模'
}

// 筛选状态
const searchQuery = ref('')
const selectedBrands = ref([])
const selectedShapes = ref([])
const weightRange = ref([40, 150])
const defaultWeightRange = [40, 150]
const selectedConnections = ref([])
const buttonRange = ref([2, 20])
const defaultButtonRange = [2, 20]

// 重置筛选
const resetFilters = () => {
  searchQuery.value = ''
  selectedBrands.value = []
  selectedShapes.value = []
  weightRange.value = [...defaultWeightRange]
  selectedConnections.value = []
  buttonRange.value = [...defaultButtonRange]
  
  // 发出事件通知父组件筛选已重置
  emits('reset-filters')
}

// 应用筛选
const applyFilters = () => {
  const filters = {
    query: searchQuery.value,
    brands: selectedBrands.value,
    shapes: selectedShapes.value,
    weight: weightRange.value,
    connections: selectedConnections.value,
    buttons: buttonRange.value
  }
  
  // 发出筛选事件
  emits('apply-filters', filters)
}

// 定义组件事件
const emits = defineEmits(['apply-filters', 'reset-filters'])

// 监听筛选条件变化，实时更新
watch([searchQuery, selectedBrands, selectedShapes, weightRange, selectedConnections, buttonRange], () => {
  // 自动应用筛选，可以根据需要取消注释
  // applyFilters()
})
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

.filter-checkbox-group {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.filter-checkbox :deep(.el-checkbox__label) {
  color: var(--claude-text-light);
}

.filter-checkbox :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: var(--claude-primary-purple);
  border-color: var(--claude-primary-purple);
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