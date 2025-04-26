<template>
  <div class="filter-panel">
    <div class="filter-header">
      <h3 class="text-lg font-medium">筛选条件</h3>
      <el-button 
        v-if="hasActiveFilters" 
        type="primary"
        text
        @click="clearAllFilters"
        class="text-claude-primary-purple hover:text-claude-primary-purple-light"
      >
        清除全部
      </el-button>
    </div>
    
    <el-divider />
    
    <div class="filter-groups">
      <div 
        v-for="filter in filters" 
        :key="filter.id" 
        class="filter-group"
      >
        <div class="filter-title">
          <h4 class="text-base font-medium">{{ filter.label }}</h4>
        </div>
        
        <!-- 复选框过滤器 -->
        <template v-if="filter.type === 'checkbox'">
          <el-checkbox-group 
            v-model="activeFilters[filter.id]" 
            @change="(val) => emitFilterChange(filter.id, val)"
          >
            <div class="checkbox-list">
              <el-checkbox 
                v-for="option in filter.options" 
                :key="option.value" 
                :label="option.value"
              >
                {{ option.label }}
              </el-checkbox>
            </div>
          </el-checkbox-group>
        </template>
        
        <!-- 范围过滤器 -->
        <template v-else-if="filter.type === 'range'">
          <div class="range-filter">
            <div class="range-values">
              <span>{{ activeRanges[filter.id]?.min || filter.min }}</span>
              <span>{{ activeRanges[filter.id]?.max || filter.max }}{{ filter.unit }}</span>
            </div>
            
            <el-slider
              v-model="activeRanges[filter.id]"
              range
              :min="filter.min"
              :max="filter.max"
              @change="(val) => emitFilterChange(filter.id, val)"
              :button-color="'var(--claude-primary-purple)'"
              :active-color="'var(--claude-primary-purple)'"
            />
          </div>
        </template>
        
        <!-- 单选按钮过滤器 -->
        <template v-else-if="filter.type === 'radio'">
          <el-radio-group 
            v-model="activeRadios[filter.id]" 
            @change="(val) => emitFilterChange(filter.id, val)"
          >
            <div class="radio-list">
              <el-radio 
                v-for="option in filter.options" 
                :key="option.value" 
                :label="option.value"
              >
                {{ option.label }}
              </el-radio>
            </div>
          </el-radio-group>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, defineProps, defineEmits } from 'vue';
import { ElMessage } from 'element-plus';

interface FilterOption {
  label: string;
  value: string | number;
}

interface CheckboxFilter {
  id: string;
  label: string;
  type: 'checkbox';
  options: FilterOption[];
}

interface RadioFilter {
  id: string;
  label: string;
  type: 'radio';
  options: FilterOption[];
}

interface RangeFilter {
  id: string;
  label: string;
  type: 'range';
  min: number;
  max: number;
  unit?: string;
}

type Filter = CheckboxFilter | RadioFilter | RangeFilter;

const props = defineProps<{
  filters: Filter[];
  appliedFilters?: Record<string, any>;
}>();

const emit = defineEmits<{
  (e: 'filter-change', filterId: string, value: any): void;
}>();

// 本地状态
const activeFilters = ref<Record<string, string[]>>({});
const activeRanges = ref<Record<string, { min: number; max: number }>>({});
const activeRadios = ref<Record<string, string | number>>({});

// 初始化过滤器状态
function initializeFilters() {
  if (props.appliedFilters) {
    for (const [key, value] of Object.entries(props.appliedFilters)) {
      const filter = props.filters.find(f => f.id === key);
      
      if (!filter) continue;
      
      if (filter.type === 'checkbox') {
        activeFilters.value[key] = Array.isArray(value) ? value : [];
      } else if (filter.type === 'range') {
        activeRanges.value[key] = typeof value === 'object' ? value : { min: filter.min, max: filter.max };
      } else if (filter.type === 'radio') {
        activeRadios.value[key] = value;
      }
    }
  }
  
  // 为空的过滤器初始化默认值
  props.filters.forEach(filter => {
    if (filter.type === 'checkbox' && !activeFilters.value[filter.id]) {
      activeFilters.value[filter.id] = [];
    } else if (filter.type === 'range' && !activeRanges.value[filter.id]) {
      activeRanges.value[filter.id] = { min: filter.min, max: filter.max };
    }
  });
}

// 监听过滤器变化和重置
watch(() => props.filters, () => {
  initializeFilters();
}, { immediate: true });

watch(() => props.appliedFilters, (newFilters) => {
  if (newFilters) {
    initializeFilters();
  }
}, { deep: true });

// 清除所有过滤器
function clearAllFilters() {
  activeFilters.value = {};
  activeRanges.value = {};
  activeRadios.value = {};
  
  props.filters.forEach(filter => {
    if (filter.type === 'checkbox') {
      activeFilters.value[filter.id] = [];
      emit('filter-change', filter.id, []);
    } else if (filter.type === 'range') {
      activeRanges.value[filter.id] = { min: filter.min, max: filter.max };
      emit('filter-change', filter.id, { min: filter.min, max: filter.max });
    } else if (filter.type === 'radio') {
      activeRadios.value[filter.id] = '';
      emit('filter-change', filter.id, '');
    }
  });
  
  // 添加清除按钮动效
  ElMessage({
    message: '已清除所有筛选条件',
    type: 'success',
    duration: 2000
  });
}

// 发射过滤器变化事件
function emitFilterChange(filterId: string, value: any) {
  emit('filter-change', filterId, value);
}

// 计算是否有激活的过滤器
const hasActiveFilters = computed(() => {
  // 检查复选框过滤器
  for (const [_, value] of Object.entries(activeFilters.value)) {
    if (Array.isArray(value) && value.length > 0) {
      return true;
    }
  }
  
  // 检查范围过滤器
  for (const [filterId, value] of Object.entries(activeRanges.value)) {
    const filter = props.filters.find(f => f.id === filterId && f.type === 'range') as RangeFilter | undefined;
    if (filter && (value.min > filter.min || value.max < filter.max)) {
      return true;
    }
  }
  
  // 检查单选按钮过滤器
  for (const [_, value] of Object.entries(activeRadios.value)) {
    if (value) {
      return true;
    }
  }
  
  return false;
});

// 初始化
initializeFilters();
</script>

<style scoped>
.filter-panel {
  width: 100%;
  background-color: var(--claude-bg-medium);
  border-radius: var(--el-border-radius-base);
  border: 1px solid var(--claude-border-dark);
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: var(--claude-text-white);
}

.filter-groups {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.filter-group {
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--claude-border-dark);
}

.filter-group:last-child {
  border-bottom: none;
}

.filter-title {
  margin-bottom: 0.75rem;
  color: var(--claude-primary-purple-light);
}

.checkbox-list, .radio-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.range-filter {
  padding: 0 0.5rem;
}

.range-values {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: var(--claude-text-light);
}
</style>