<template>
  <div class="mouse-database-container">
    <div class="database-header">
      <h1 class="text-2xl font-bold">鼠标数据库</h1>
      <div class="header-controls">
        <el-input
          v-model="searchQuery"
          placeholder="搜索鼠标名称、品牌或形状"
          prefix-icon="el-icon-search"
          clearable
          class="search-input"
        />
        <ViewToggle 
          :initial-view-mode="viewMode" 
          @view-change="handleViewModeChange" 
        />
      </div>
    </div>
    
    <div class="database-content">
      <div class="filter-panel-container">
        <FilterPanel 
          :filters="filters"
          :applied-filters="appliedFilters"
          @filter-change="handleFilterChange"
        />
      </div>
      
      <div class="results-container">
        <div class="results-header">
          <div class="results-count">
            找到 <span class="font-bold">{{ filteredMice.length }}</span> 个鼠标
          </div>
          <SortControls 
            :sort-options="sortOptions"
            :initial-sort-by="sortBy"
            :initial-sort-direction="sortDirection"
            @sort-change="handleSortChange"
          />
        </div>
        
        <div v-if="loading" class="loading-container">
          <el-spinner size="large" />
        </div>
        
        <div v-else-if="!filteredMice.length" class="empty-state">
          <el-empty description="没有找到符合条件的鼠标" />
        </div>
        
        <!-- 网格视图 -->
        <div v-else-if="viewMode === 'grid'" class="grid-view">
          <el-card
            v-for="mouse in displayedMice"
            :key="mouse.id"
            class="mouse-card"
            shadow="hover"
            @click="navigateToMouseDetail(mouse.id)"
          >
            <div class="mouse-card-content">
              <div class="mouse-image">
                <img v-if="mouse.imageUrl" :src="mouse.imageUrl" alt="鼠标图片" class="img-fluid" />
                <div v-else class="placeholder-image">
                  <i class="el-icon-mouse"></i>
                </div>
              </div>
              <div class="mouse-info">
                <h3 class="mouse-name">{{ mouse.brand }} {{ mouse.name }}</h3>
                <div class="mouse-specs">
                  <div class="spec-item">
                    <span class="spec-label">尺寸:</span>
                    <span class="spec-value">{{ mouse.dimensions.length }}x{{ mouse.dimensions.width }}x{{ mouse.dimensions.height }}mm</span>
                  </div>
                  <div class="spec-item">
                    <span class="spec-label">重量:</span>
                    <span class="spec-value">{{ mouse.weight }}g</span>
                  </div>
                  <div class="spec-item">
                    <span class="spec-label">形状:</span>
                    <span class="spec-value">{{ mouse.shape.type }}</span>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </div>
        
        <!-- 列表视图 -->
        <div v-else class="list-view">
          <el-table
            :data="displayedMice"
            stripe
            border
            @row-click="(row) => navigateToMouseDetail(row.id)"
            class="mouse-table"
          >
            <el-table-column label="品牌" prop="brand" min-width="80" />
            <el-table-column label="型号" prop="name" min-width="120" />
            <el-table-column label="尺寸 (mm)" min-width="150">
              <template #default="{ row }">
                {{ row.dimensions.length }}x{{ row.dimensions.width }}x{{ row.dimensions.height }}
              </template>
            </el-table-column>
            <el-table-column label="重量 (g)" prop="weight" min-width="80" />
            <el-table-column label="形状" min-width="100">
              <template #default="{ row }">
                {{ row.shape.type }}
              </template>
            </el-table-column>
            <el-table-column label="连接方式" prop="connectivity" min-width="100" />
            <el-table-column label="操作" fixed="right" width="120">
              <template #default="{ row }">
                <el-button 
                  type="text"
                  @click.stop="addToComparison(row)"
                  :disabled="isInComparison(row.id)"
                >
                  {{ isInComparison(row.id) ? '已添加' : '添加比较' }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <!-- 分页控件 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="filteredMice.length"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useComparisonStore } from '@/stores';
import type { MouseDevice } from '@/models/MouseModel';
import { getDevices } from '@/api/device';
import FilterPanel from '@/components/database/FilterPanel.vue';
import SortControls from '@/components/database/SortControls.vue';
import ViewToggle from '@/components/database/ViewToggle.vue';
import type { SortDirection } from '@/components/database/SortControls.vue';
import type { ViewMode } from '@/components/database/ViewToggle.vue';

// 路由
const router = useRouter();
const comparisonStore = useComparisonStore();

// 状态
const allMice = ref<MouseDevice[]>([]);
const loading = ref(true);
const searchQuery = ref('');
const viewMode = ref<ViewMode>('grid');
const sortBy = ref('weight');
const sortDirection = ref<SortDirection>('asc');
const currentPage = ref(1);
const pageSize = ref(20);

// 过滤器配置
const filters = [
  {
    id: 'brand',
    label: '品牌',
    type: 'checkbox',
    options: [
      { label: 'Logitech', value: 'Logitech' },
      { label: 'Razer', value: 'Razer' },
      { label: 'Glorious', value: 'Glorious' },
      { label: 'Zowie', value: 'Zowie' },
      { label: 'SteelSeries', value: 'SteelSeries' }
    ]
  },
  {
    id: 'shape.type',
    label: '形状',
    type: 'checkbox',
    options: [
      { label: '对称', value: 'Symmetrical' },
      { label: '右手人体工学', value: 'Ergonomic' },
      { label: '低矮', value: 'Low Profile' },
      { label: '大型', value: 'Large' },
      { label: '中型', value: 'Medium' },
      { label: '小型', value: 'Small' }
    ]
  },
  {
    id: 'connectivity',
    label: '连接方式',
    type: 'checkbox',
    options: [
      { label: '有线', value: 'Wired' },
      { label: '无线', value: 'Wireless' },
      { label: '双模', value: 'Dual-mode' }
    ]
  },
  {
    id: 'weight',
    label: '重量范围',
    type: 'range',
    min: 40,
    max: 150,
    unit: 'g'
  }
];

// 排序选项
const sortOptions = [
  { label: '重量', value: 'weight' },
  { label: '长度', value: 'dimensions.length' },
  { label: '宽度', value: 'dimensions.width' },
  { label: '高度', value: 'dimensions.height' },
  { label: '品牌', value: 'brand' },
  { label: '型号', value: 'name' }
];

// 应用的过滤器
const appliedFilters = ref<Record<string, any>>({});

// 计算属性: 过滤后的鼠标列表
const filteredMice = computed(() => {
  return allMice.value.filter(mouse => {
    // 搜索过滤
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const nameMatch = `${mouse.brand} ${mouse.name}`.toLowerCase().includes(query);
      const shapeMatch = mouse.shape.type.toLowerCase().includes(query);
      
      if (!nameMatch && !shapeMatch) {
        return false;
      }
    }
    
    // 应用过滤器
    for (const [filterId, filterValue] of Object.entries(appliedFilters.value)) {
      // 跳过空过滤器
      if (!filterValue || (Array.isArray(filterValue) && filterValue.length === 0)) {
        continue;
      }
      
      // 范围过滤器
      if (filterId === 'weight' && typeof filterValue === 'object') {
        const { min, max } = filterValue;
        if (mouse.weight < min || mouse.weight > max) {
          return false;
        }
        continue;
      }
      
      // 复选框过滤器
      if (Array.isArray(filterValue) && filterValue.length > 0) {
        const parts = filterId.split('.');
        let value;
        
        if (parts.length === 1) {
          value = mouse[filterId as keyof typeof mouse];
        } else if (parts.length === 2) {
          const [obj, prop] = parts;
          value = mouse[obj as keyof typeof mouse]?.[prop];
        }
        
        if (!filterValue.includes(value)) {
          return false;
        }
      }
    }
    
    return true;
  });
});

// 计算属性: 排序后的鼠标列表
const sortedMice = computed(() => {
  return [...filteredMice.value].sort((a, b) => {
    const parts = sortBy.value.split('.');
    let valueA, valueB;
    
    if (parts.length === 1) {
      valueA = a[sortBy.value as keyof typeof a];
      valueB = b[sortBy.value as keyof typeof b];
    } else if (parts.length === 2) {
      const [obj, prop] = parts;
      valueA = a[obj as keyof typeof a]?.[prop];
      valueB = b[obj as keyof typeof b]?.[prop];
    }
    
    if (typeof valueA === 'string' && typeof valueB === 'string') {
      return sortDirection.value === 'asc' 
        ? valueA.localeCompare(valueB) 
        : valueB.localeCompare(valueA);
    }
    
    // 数字比较
    return sortDirection.value === 'asc' 
      ? Number(valueA) - Number(valueB) 
      : Number(valueB) - Number(valueA);
  });
});

// 计算属性: 显示的鼠标列表（分页）
const displayedMice = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return sortedMice.value.slice(start, end);
});

// 方法
function handleFilterChange(filterId: string, value: any) {
  appliedFilters.value = { ...appliedFilters.value, [filterId]: value };
  // 过滤器变化时重置到第一页
  currentPage.value = 1;
}

function handleSortChange(newSortBy: string, newSortDirection: SortDirection) {
  sortBy.value = newSortBy;
  sortDirection.value = newSortDirection;
}

function handleViewModeChange(mode: ViewMode) {
  viewMode.value = mode;
}

function handleSizeChange(newSize: number) {
  pageSize.value = newSize;
}

function handleCurrentChange(newPage: number) {
  currentPage.value = newPage;
}

function navigateToMouseDetail(mouseId: string) {
  router.push(`/mice/${mouseId}`);
}

function addToComparison(mouse: MouseDevice) {
  if (comparisonStore.selectedMice.length >= 3) {
    // 已达到最大比较数量
    window.ElMessage.warning('最多只能选择3个鼠标进行比较');
    return;
  }
  
  comparisonStore.addMouse(mouse);
  window.ElMessage.success(`已添加 ${mouse.brand} ${mouse.name} 到比较列表`);
}

function isInComparison(mouseId: string) {
  return comparisonStore.selectedMice.some(m => m.id === mouseId);
}

// 获取鼠标数据
async function fetchMice() {
  loading.value = true;
  try {
    const response = await getDevices({ type: 'mouse' });
    const deviceListResponse = response as unknown as DeviceListResponse;
    allMice.value = deviceListResponse.devices as MouseDevice[];
  } catch (error) {
    console.error('Error fetching mice:', error);
    ElMessage.error('获取鼠标数据失败');
  } finally {
    loading.value = false;
  }
}

// 生命周期钩子
onMounted(() => {
  fetchMice();
});
</script>

<style scoped>
.mouse-database-container {
  padding: 1rem;
}

.database-header {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.header-controls {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  align-items: center;
}

.search-input {
  width: 300px;
}

.database-content {
  display: grid;
  grid-template-columns: 250px 1fr;
  gap: 1.5rem;
}

.filter-panel-container {
  position: sticky;
  top: 1rem;
  height: fit-content;
  background-color: var(--claude-bg-medium);
  border-radius: var(--el-border-radius-base);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  padding: 1rem;
  border: 1px solid var(--claude-border-dark);
}

.results-container {
  width: 100%;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}

.mouse-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.mouse-card:hover {
  transform: translateY(-2px);
}

.mouse-card-content {
  display: flex;
  flex-direction: column;
}

.mouse-image {
  height: 150px;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 1rem;
  background-color: #f9f9f9;
  border-radius: 0.25rem;
}

.placeholder-image {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #b0b0b0;
  font-size: 3rem;
}

.mouse-info {
  margin-top: auto;
}

.mouse-name {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.mouse-specs {
  font-size: 0.875rem;
}

.spec-item {
  display: flex;
  margin-bottom: 0.25rem;
}

.spec-label {
  color: #909399;
  margin-right: 0.5rem;
  width: 40px;
}

.mouse-table {
  margin-bottom: 1rem;
}

.mouse-table tbody tr {
  cursor: pointer;
}

.empty-state, .loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 1.5rem;
}

@media (max-width: 768px) {
  .database-content {
    grid-template-columns: 1fr;
  }
  
  .filter-panel-container {
    position: static;
    margin-bottom: 1rem;
  }
  
  .search-input {
    width: 100%;
  }
}
</style>