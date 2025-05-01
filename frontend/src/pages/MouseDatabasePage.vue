<template>
  <div class="mouse-database-container">
    <div class="database-header">
      <h1 class="text-2xl font-bold">Mouse Database</h1>
      <div class="header-controls">
        <el-input
          v-model="searchQuery"
          placeholder="Search mouse name, brand, or shape"
          prefix-icon="el-icon-search"
          clearable
          class="search-input"
        />
        <ViewToggle :initial-view-mode="viewMode" @view-change="handleViewModeChange" />
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
            Found <span class="font-bold">{{ filteredMice.length }}</span> mice
          </div>
          <SortControls
            :sort-options="sortOptions"
            :initial-sort-by="sortBy"
            :initial-sort-direction="sortDirection"
            @sort-change="handleSortChange"
          />
        </div>

        <div v-if="loading" v-loading="true" class="loading-container"></div>

        <div v-else-if="!filteredMice.length" class="empty-state">
          <el-empty description="No mice found matching your criteria" />
        </div>

        <!-- 网格视图 -->
        <div v-else-if="viewMode === 'grid'" class="grid-view">
          <el-card v-for="mouse in displayedMice" :key="mouse.id" class="mouse-card">
            <div class="mouse-card-content">
              <div class="mouse-image" @click="navigateToMouseDetail(mouse.id)">
                <img v-if="mouse.imageUrl" :src="mouse.imageUrl" alt="Mouse image" class="img-fluid" />
                <div v-else class="placeholder-image">
                  <i class="el-icon-mouse"></i>
                </div>
              </div>
              <div class="mouse-info">
                <h3 class="mouse-name" @click="navigateToMouseDetail(mouse.id)">
                  {{ mouse.brand }} {{ mouse.name }}
                </h3>
                <div class="mouse-specs">
                  <div class="spec-item">
                    <span class="spec-label">Size:</span>
                    <span class="spec-value" v-if="mouse.dimensions">
                      {{ mouse.dimensions.length }}x{{ mouse.dimensions.width }}x{{ mouse.dimensions.height }}mm
                    </span>
                    <span class="spec-value" v-else>N/A</span>
                  </div>
                  <div class="spec-item">
                    <span class="spec-label">Weight:</span>
                    <span class="spec-value" v-if="mouse.weight">{{ mouse.weight }}g</span>
                    <span class="spec-value" v-else>N/A</span>
                  </div>
                  <div class="spec-item">
                    <span class="spec-label">Shape:</span>
                    <span class="spec-value" v-if="mouse.shape">{{ mouse.shape.type }}</span>
                    <span class="spec-value" v-else>N/A</span>
                  </div>
                </div>
                <div class="mouse-card-actions">
                  <div class="price-tag">¥{{ mouse.price || '699' }}</div>
                  <div class="action-buttons">
                    <el-button
                      type="primary"
                      size="small"
                      :disabled="isInComparison(mouse.id)"
                      icon="el-icon-sort"
                      circle
                      title="Add to comparison"
                      @click.stop="addToComparison(mouse)"
                    ></el-button>
                    <AddToCartButton
                      :product="{
                        id: mouse.id,
                        name: `${mouse.brand} ${mouse.name}`,
                        price: mouse.price || 699,
                        imageUrl: mouse.imageUrl
                      }"
                      product-type="mouse"
                      type="primary"
                      size="small"
                      icon-only
                    />
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </div>

        <!-- 列表视图 -->
        <div v-else class="list-view">
          <el-table :data="displayedMice" stripe border class="mouse-table">
            <el-table-column label="Brand" prop="brand" min-width="80" />
            <el-table-column label="Model" prop="name" min-width="120">
              <template #default="{ row }">
                <a class="mouse-link" @click="navigateToMouseDetail(row.id)">{{ row.name }}</a>
              </template>
            </el-table-column>
            <el-table-column label="Size (mm)" min-width="150">
              <template #default="{ row }">
                <template v-if="row.dimensions">
                  {{ row.dimensions.length }}x{{ row.dimensions.width }}x{{ row.dimensions.height }}
                </template>
                <template v-else>N/A</template>
              </template>
            </el-table-column>
            <el-table-column label="Weight (g)" min-width="80">
              <template #default="{ row }">
                {{ row.weight || 'N/A' }}
              </template>
            </el-table-column>
            <el-table-column label="Shape" min-width="100">
              <template #default="{ row }">
                <template v-if="row.shape">{{ row.shape.type }}</template>
                <template v-else>N/A</template>
              </template>
            </el-table-column>
            <el-table-column label="Connectivity" prop="connectivity" min-width="100" />
            <el-table-column label="Price" prop="price" min-width="80">
              <template #default="{ row }"> ¥{{ row.price || '699' }} </template>
            </el-table-column>
            <el-table-column label="Actions" fixed="right" width="160">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-button
                    type="primary"
                    size="small"
                    icon="el-icon-sort"
                    circle
                    :disabled="isInComparison(row.id)"
                    title="Add to comparison"
                    @click.stop="addToComparison(row)"
                  ></el-button>
                  <AddToCartButton
                    :product="{
                      id: row.id,
                      name: `${row.brand} ${row.name}`,
                      price: row.price || 699,
                      imageUrl: row.imageUrl
                    }"
                    product-type="mouse"
                    type="primary"
                    size="small"
                    icon-only
                  />
                </div>
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
import { DeviceListResponse } from '@/api/device';
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useComparisonStore } from '@/stores';
import type { MouseDevice } from '@/models/MouseModel';
import { getDevices } from '@/api/device';
import { ElMessage } from 'element-plus';
import FilterPanel from '@/components/database/FilterPanel.vue';
import SortControls from '@/components/database/SortControls.vue';
import ViewToggle from '@/components/database/ViewToggle.vue';
import AddToCartButton from '@/components/cart/AddToCartButton.vue';
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
    label: 'Brand',
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
    label: 'Shape',
    type: 'checkbox',
    options: [
      { label: 'Symmetrical', value: 'Symmetrical' },
      { label: 'Ergonomic', value: 'Ergonomic' },
      { label: 'Low Profile', value: 'Low Profile' },
      { label: 'Large', value: 'Large' },
      { label: 'Medium', value: 'Medium' },
      { label: 'Small', value: 'Small' }
    ]
  },
  {
    id: 'connectivity',
    label: 'Connectivity',
    type: 'checkbox',
    options: [
      { label: 'Wired', value: 'Wired' },
      { label: 'Wireless', value: 'Wireless' },
      { label: 'Dual-mode', value: 'Dual-mode' }
    ]
  },
  {
    id: 'weight',
    label: 'Weight Range',
    type: 'range',
    min: 40,
    max: 150,
    unit: 'g'
  }
];

// 排序选项
const sortOptions = [
  { label: 'Weight', value: 'weight' },
  { label: 'Length', value: 'dimensions.length' },
  { label: 'Width', value: 'dimensions.width' },
  { label: 'Height', value: 'dimensions.height' },
  { label: 'Brand', value: 'brand' },
  { label: 'Model', value: 'name' },
  { label: 'Price', value: 'price' }
  // 已移除'最新上架'选项
];

// 应用的过滤器
const appliedFilters = ref<Record<string, any>>({});

// 计算属性: 过滤后的鼠标列表
const filteredMice = computed(() => {
  return allMice.value.filter((mouse) => {
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
    ElMessage.warning('You can only select up to 3 mice for comparison');
    return;
  }

  comparisonStore.addMouse(mouse);
  ElMessage.success(`Added ${mouse.brand} ${mouse.name} to comparison list`);
}

function isInComparison(mouseId: string) {
  return comparisonStore.selectedMice.some((m) => m.id === mouseId);
}

// 获取鼠标数据
async function fetchMice() {
  loading.value = true;
  try {
    const response = await getDevices({ type: 'mouse' });
    console.log('API Response:', response); // 添加调试日志
    
    let deviceListResponse: DeviceListResponse;
    
    // 处理各种可能的响应格式
    if (response && 'devices' in response) {
      deviceListResponse = response as DeviceListResponse;
    } else if (response && response.data && 'devices' in response.data) {
      deviceListResponse = response.data as DeviceListResponse;
    } else {
      console.error('Unexpected API response format:', response);
      ElMessage.error('Invalid API response format');
      allMice.value = [];
      loading.value = false;
      return;
    }
    
    if (Array.isArray(deviceListResponse.devices)) {
      allMice.value = deviceListResponse.devices as MouseDevice[];
      console.log('Loaded mice:', allMice.value.length);
    } else {
      console.error('No devices array in response:', deviceListResponse);
      allMice.value = [];
    }
  } catch (error) {
    console.error('Error fetching mice:', error);
    ElMessage.error('Failed to fetch mouse data');
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
  transition: transform 0.2s;
}

.mouse-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
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
  cursor: pointer;
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
  cursor: pointer;
}

.mouse-name:hover {
  color: var(--claude-primary-purple);
}

.mouse-specs {
  font-size: 0.875rem;
  margin-bottom: 1rem;
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

.mouse-card-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px solid var(--claude-border-light);
}

.price-tag {
  font-weight: 600;
  font-size: 1.125rem;
  color: #e6603c;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
}

.mouse-table {
  margin-bottom: 1rem;
}

.table-actions {
  display: flex;
  gap: 0.5rem;
}

.mouse-link {
  color: var(--claude-primary-purple);
  text-decoration: none;
  cursor: pointer;
}

.mouse-link:hover {
  text-decoration: underline;
}

.empty-state,
.loading-container {
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
