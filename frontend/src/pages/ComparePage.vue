<template>
  <div class="compare-page">
    <div class="container mx-auto py-6 px-4">
      <h1 class="text-2xl font-bold mb-6">鼠标比较工具</h1>

      <el-card class="mb-6">
        <template #header>
          <div class="flex justify-between items-center">
            <h2 class="text-xl">选择鼠标</h2>
            <el-button v-if="selectedMice.length > 0" type="primary" plain @click="resetComparison">
              重置比较
            </el-button>
          </div>
        </template>

        <div v-if="selectedMice.length === 0" class="py-4">
          <p class="text-center text-gray-500 mb-4">请选择要比较的鼠标</p>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-card class="h-full" shadow="hover" @click="openMouseSelectorMode('mice')">
                <div class="flex flex-col items-center py-6">
                  <el-icon class="text-3xl mb-3"><Plus /></el-icon>
                  <h3 class="text-lg font-medium">选择鼠标</h3>
                  <p class="text-sm text-gray-500 mt-2">从我们的数据库中选择鼠标进行比较</p>
                </div>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card class="h-full" shadow="hover" @click="openMouseSelectorMode('recent')">
                <div class="flex flex-col items-center py-6">
                  <el-icon class="text-3xl mb-3"><Clock /></el-icon>
                  <h3 class="text-lg font-medium">最近浏览的鼠标</h3>
                  <p class="text-sm text-gray-500 mt-2">从最近浏览的鼠标中选择</p>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <div v-else>
          <div class="selected-mice flex flex-wrap gap-4">
            <el-card
              v-for="mouse in selectedMice"
              :key="mouse.id"
              class="selected-mouse-card"
              shadow="hover"
            >
              <div class="flex flex-col items-center">
                <div class="mouse-image w-32 h-32 mb-2">
                  <template v-if="mouse.imageUrl">
                    <img
                      :src="mouse.imageUrl"
                      class="w-full h-full object-contain"
                      :alt="mouse.name"
                    />
                  </template>
                  <template v-else-if="mouse.svgData?.topView">
                    <div class="w-full h-full svg-container">{{ mouse.svgData.topView }}</div>
                  </template>
                  <template v-else>
                    <div class="w-full h-full bg-gray-100 flex items-center justify-center">
                      <el-icon class="text-4xl text-gray-300"><Mouse /></el-icon>
                    </div>
                  </template>
                </div>
                <h3 class="text-base font-medium">{{ mouse.brand }} {{ mouse.name }}</h3>
                <p class="text-xs text-gray-500 mt-1">
                  {{ mouse.dimensions?.length || '未知' }}×{{
                    mouse.dimensions?.width || '未知'
                  }}×{{ mouse.dimensions?.height || '未知' }}mm
                </p>
                <p class="text-xs text-gray-500">
                  {{ mouse.dimensions?.weight || '未知' }}g · {{ mouse.shape?.type || '未知' }}
                </p>
                <el-button
                  size="small"
                  type="danger"
                  text
                  class="mt-2"
                  @click="removeMouse(mouse.id)"
                >
                  移除
                </el-button>
              </div>
            </el-card>

            <!-- 添加鼠标卡片 -->
            <el-card
              v-if="selectedMice.length < 3"
              class="selected-mouse-card add-card"
              shadow="hover"
              @click="openMouseSelectorMode('mice')"
            >
              <div class="flex flex-col items-center justify-center h-full">
                <el-icon class="text-4xl mb-2"><Plus /></el-icon>
                <p class="text-sm">添加鼠标</p>
              </div>
            </el-card>
          </div>
        </div>
      </el-card>

      <!-- 比较内容 -->
      <MouseComparisonView v-if="selectedMice.length >= 2" />

      <!-- 移除相似鼠标推荐窗口，改为提示去相似页面 -->
      <el-card v-if="selectedMice.length === 1" class="mt-6">
        <div class="py-4 text-center">
          <p class="text-gray-500 mb-4">您当前只选择了一个鼠标，需要至少两个鼠标才能进行比较。</p>
          <div class="flex justify-center gap-4">
            <el-button type="primary" @click="openMouseSelectorMode('mice')">
              添加另一个鼠标
            </el-button>
            <router-link to="/similar">
              <el-button> 查找相似鼠标 </el-button>
            </router-link>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 选择鼠标对话框 -->
    <el-dialog v-model="showMouseDialog" :title="dialogTitle" width="70%" destroy-on-close>
      <MouseSelector
        v-if="selectMode === 'mice'"
        :initial-selected-mice="selectedMice"
        :max-selection="3"
        @select="handleMouseSelection"
        @cancel="showMouseDialog = false"
      />

      <div v-else-if="selectMode === 'recent'" class="recent-mice">
        <div v-if="recentlyViewedMice.length === 0" class="empty-state text-center py-12">
          <el-icon class="text-3xl mb-3"><InfoFilled /></el-icon>
          <p class="text-gray-500">您还没有浏览过任何鼠标</p>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <el-card
            v-for="mouse in recentlyViewedMice"
            :key="mouse.id"
            class="recent-mouse-card"
            shadow="hover"
            @click="addMouse(mouse)"
          >
            <div class="flex items-center">
              <div class="mouse-image w-16 h-16 mr-3 flex-shrink-0">
                <template v-if="mouse.imageUrl">
                  <img
                    :src="mouse.imageUrl"
                    class="w-full h-full object-contain"
                    :alt="mouse.name"
                  />
                </template>
                <template v-else>
                  <div class="w-full h-full bg-gray-100 flex items-center justify-center">
                    <el-icon class="text-2xl text-gray-300"><Mouse /></el-icon>
                  </div>
                </template>
              </div>
              <div>
                <h3 class="text-base font-medium">{{ mouse.brand }} {{ mouse.name }}</h3>
                <p class="text-xs text-gray-500 mt-1">
                  {{ mouse.dimensions?.length || '未知' }}×{{
                    mouse.dimensions?.width || '未知'
                  }}×{{ mouse.dimensions?.height || '未知' }}mm
                </p>
                <p class="text-xs text-gray-500">
                  {{ mouse.dimensions?.weight || '未知' }}g · {{ mouse.shape?.type || '未知' }}
                </p>
              </div>
            </div>
          </el-card>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showMouseDialog = false">取消</el-button>
          <el-button type="primary" @click="showMouseDialog = false">确认</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useComparisonStore, type MouseDevice } from '@/stores';
import MouseComparisonView from '@/components/comparison/MouseComparisonView.vue';
import MouseSelector from '@/components/comparison/MouseSelector.vue';
import comparisonService from '@/services/comparisonService';
import { getDevices, DeviceListResponse } from '@/api/device';
import { Plus, Clock, Mouse, Loading, InfoFilled } from '@element-plus/icons-vue';

// 状态
const comparisonStore = useComparisonStore();
const showMouseDialog = ref(false);
const selectMode = ref<'mice' | 'recent'>('mice');
const similarMice = ref<MouseDevice[]>([]);
const allMice = ref<MouseDevice[]>([]);
const loading = ref(false);

// 计算属性
const selectedMice = computed(() => comparisonStore.selectedMice);
const recentlyViewedMice = computed(() => comparisonStore.recentlyViewedMice);

const dialogTitle = computed(() => {
  if (selectMode.value === 'mice') return '选择鼠标';
  return '最近浏览的鼠标';
});

// 方法
function removeMouse(mouseId: string) {
  comparisonStore.removeMouse(mouseId);
}

function addMouse(mouse: MouseDevice) {
  comparisonStore.addMouse(mouse);
  comparisonStore.addToRecentlyViewed(mouse);
  showMouseDialog.value = false;
  // 解决最近浏览只出现一次的问题
  selectMode.value = 'mice'; // 重置选择模式，这样下次点击"最近浏览"时会重新加载
}

function resetComparison() {
  comparisonStore.clearSelection();
}

function handleMouseSelection(mice: MouseDevice[]) {
  // MouseSelector组件已处理选择逻辑
  // 这里只需要关闭对话框
  showMouseDialog.value = false;
}

// 根据当前选择的鼠标寻找相似鼠标
async function findSimilarMice() {
  if (selectedMice.value.length !== 1 || allMice.value.length === 0) {
    similarMice.value = [];
    return;
  }

  // 使用本地相似度服务找出相似鼠标
  similarMice.value = comparisonService.findSimilarMice(selectedMice.value[0], allMice.value, 5);
}

// 获取所有鼠标数据
async function fetchAllMice() {
  loading.value = true;
  try {
    const response = await getDevices({ type: 'mouse' });
    const deviceListResponse = response as unknown as DeviceListResponse;
    allMice.value = deviceListResponse.devices as MouseDevice[];
    findSimilarMice();
  } catch (error) {
    console.error('Error fetching mice:', error);
  } finally {
    loading.value = false;
  }
}

// 打开鼠标选择器
function openMouseSelectorMode(mode: 'mice' | 'recent') {
  selectMode.value = mode;
  showMouseDialog.value = true;
}

// 监听选择的鼠标变化，更新相似鼠标
watch(selectedMice, () => {
  findSimilarMice();
});

// 初始化
onMounted(() => {
  fetchAllMice();
});
</script>

<style scoped>
.selected-mouse-card {
  width: 200px;
}

.add-card {
  width: 200px;
  height: 275px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  border: 2px dashed var(--el-border-color);
}

.add-card:hover {
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
}

.mouse-image :deep(svg) {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.similar-mouse-card,
.recent-mouse-card {
  cursor: pointer;
  transition: all 0.2s ease;
}

.similar-mouse-card:hover,
.recent-mouse-card:hover {
  transform: translateY(-2px);
  border-color: var(--el-color-primary-light-5);
}
</style>
