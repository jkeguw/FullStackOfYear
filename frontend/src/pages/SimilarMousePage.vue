<template>
  <div class="similar-mouse-page">
    <div class="container mx-auto py-6 px-4">
      <h1 class="text-2xl font-bold mb-6">寻找相似鼠标</h1>
      
      <el-card class="mb-6">
        <template #header>
          <div class="flex justify-between items-center">
            <h2 class="text-xl">参考鼠标</h2>
            <el-button type="primary" @click="resetSelection" plain v-if="selectedMouse">
              重新选择
            </el-button>
          </div>
        </template>
        
        <div v-if="!selectedMouse" class="py-4">
          <p class="text-center text-gray-500 mb-4">请选择一个参考鼠标</p>
          <el-card shadow="hover" @click="openMouseSelector">
            <div class="flex flex-col items-center py-6">
              <el-icon class="text-3xl mb-3"><Plus /></el-icon>
              <h3 class="text-lg font-medium">选择鼠标</h3>
              <p class="text-sm text-gray-500 mt-2">从我们的数据库中选择一个鼠标作为参考</p>
            </div>
          </el-card>
        </div>
        
        <div v-else>
          <div class="selected-mouse">
            <el-card shadow="hover" class="selected-mouse-card">
              <div class="flex flex-col items-center">
                <div class="mouse-image w-32 h-32 mb-2">
                  <template v-if="selectedMouse.imageUrl">
                    <img :src="selectedMouse.imageUrl" class="w-full h-full object-contain" :alt="selectedMouse.name" />
                  </template>
                  <template v-else-if="selectedMouse.svgData?.topView">
                    <div class="w-full h-full" v-html="selectedMouse.svgData.topView"></div>
                  </template>
                  <template v-else>
                    <div class="w-full h-full bg-gray-100 flex items-center justify-center">
                      <el-icon class="text-4xl text-gray-300"><Mouse /></el-icon>
                    </div>
                  </template>
                </div>
                <h3 class="text-base font-medium">{{ selectedMouse.brand }} {{ selectedMouse.name }}</h3>
                <p class="text-xs text-gray-500 mt-1">{{ selectedMouse.dimensions?.length || '未知' }}×{{ selectedMouse.dimensions?.width || '未知' }}×{{ selectedMouse.dimensions?.height || '未知' }}mm</p>
                <p class="text-xs text-gray-500">{{ selectedMouse.weight || '未知' }}g · {{ selectedMouse.shape?.type || '未知' }}</p>
              </div>
            </el-card>
          </div>
        </div>
      </el-card>
      
      <!-- 相似鼠标列表 -->
      <el-card v-if="selectedMouse" class="mt-6">
        <template #header>
          <h2 class="text-xl">相似鼠标推荐</h2>
        </template>
        <div class="py-4">
          <p class="text-gray-500 mb-4">
            根据您选择的 {{ selectedMouse.brand }} {{ selectedMouse.name }} 鼠标，我们为您推荐了以下相似的鼠标：
          </p>
          
          <div v-if="loading" class="text-center py-8">
            <el-icon class="text-3xl mb-3"><Loading /></el-icon>
            <p>正在加载推荐...</p>
          </div>
          
          <div v-else-if="similarMice.length === 0" class="text-center py-8">
            <el-empty description="未找到相似鼠标" />
          </div>
          
          <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <el-card 
              v-for="mouse in similarMice" 
              :key="mouse.id"
              class="similar-mouse-card"
              shadow="hover"
              @click="goToCompare(mouse)"
            >
              <div class="flex items-center">
                <div class="mouse-image w-16 h-16 mr-3 flex-shrink-0">
                  <template v-if="mouse.imageUrl">
                    <img :src="mouse.imageUrl" class="w-full h-full object-contain" :alt="mouse.name" />
                  </template>
                  <template v-else>
                    <div class="w-full h-full bg-gray-100 flex items-center justify-center">
                      <el-icon class="text-2xl text-gray-300"><Mouse /></el-icon>
                    </div>
                  </template>
                </div>
                <div>
                  <h3 class="text-base font-medium">{{ mouse.brand }} {{ mouse.name }}</h3>
                  <p class="text-xs text-gray-500 mt-1">{{ mouse.dimensions?.length || '未知' }}×{{ mouse.dimensions?.width || '未知' }}×{{ mouse.dimensions?.height || '未知' }}mm</p>
                  <p class="text-xs text-gray-500">{{ mouse.weight || '未知' }}g · {{ mouse.shape?.type || '未知' }}</p>
                  
                  <div class="similarity-score flex items-center mt-2">
                    <div class="text-xs font-medium mr-2">相似度:</div>
                    <el-progress 
                      :percentage="mouse.similarityScore || 0" 
                      :color="getSimilarityColor(mouse.similarityScore || 0)"
                      :show-text="true"
                      :stroke-width="6"
                      style="width: 100px"
                    />
                  </div>
                </div>
              </div>
            </el-card>
          </div>
        </div>
      </el-card>
    </div>
    
    <!-- 鼠标选择器对话框 -->
    <el-dialog 
      v-model="showMouseSelector" 
      title="选择参考鼠标"
      width="70%"
      destroy-on-close
    >
      <div class="mouse-selector">
        <el-input 
          v-model="searchQuery" 
          placeholder="搜索鼠标品牌或名称" 
          prefix-icon="el-icon-search" 
          clearable 
        />
        
        <div class="mouse-cards-container mt-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <el-card 
            v-for="mouse in filteredMice" 
            :key="mouse.id" 
            class="mouse-card" 
            shadow="hover"
            @click="selectReferenceMouseAndFindSimilar(mouse)"
          >
            <div class="flex items-center">
              <div class="mouse-image w-16 h-16 mr-3 flex-shrink-0">
                <img v-if="mouse.imageUrl" :src="mouse.imageUrl" alt="鼠标图片" class="w-full h-full object-contain" />
                <div v-else class="w-full h-full bg-gray-100 flex items-center justify-center text-gray-400">
                  无图片
                </div>
              </div>
              <div>
                <div class="font-medium">{{ mouse.brand }} {{ mouse.name }}</div>
                <div class="text-sm text-gray-500">
                  {{ mouse.dimensions?.length || '未知' }}x{{ mouse.dimensions?.width || '未知' }}x{{ mouse.dimensions?.height || '未知' }}mm
                </div>
                <div class="text-sm text-gray-500">{{ mouse.weight || '未知' }}g · {{ mouse.shape?.type || '未知' }}</div>
              </div>
            </div>
          </el-card>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showMouseSelector = false">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useComparisonStore, type MouseDevice } from '@/stores';
import comparisonService from '@/services/comparisonService';
import { getDevices, DeviceListResponse } from '@/api/device';
import { Plus, Mouse, Loading } from '@element-plus/icons-vue';

// 状态
const comparisonStore = useComparisonStore();
const router = useRouter();
const showMouseSelector = ref(false);
const searchQuery = ref('');
const allMice = ref<MouseDevice[]>([]);
const similarMice = ref<(MouseDevice & { similarityScore?: number })[]>([]);
const selectedMouse = ref<MouseDevice | null>(null);
const loading = ref(false);

// 计算属性
const filteredMice = computed(() => {
  if (!searchQuery.value) return allMice.value;
  
  const query = searchQuery.value.toLowerCase();
  return allMice.value.filter(mouse => 
    mouse.name.toLowerCase().includes(query) || 
    mouse.brand.toLowerCase().includes(query)
  );
});

// 方法
function openMouseSelector() {
  showMouseSelector.value = true;
}

function resetSelection() {
  selectedMouse.value = null;
  similarMice.value = [];
}

function selectReferenceMouseAndFindSimilar(mouse: MouseDevice) {
  selectedMouse.value = mouse;
  comparisonStore.addToRecentlyViewed(mouse);
  showMouseSelector.value = false;
  findSimilarMice();
}

function goToCompare(mouse: MouseDevice) {
  // 添加两个鼠标到比较
  comparisonStore.clearSelection();
  if (selectedMouse.value) {
    comparisonStore.addMouse(selectedMouse.value);
  }
  comparisonStore.addMouse(mouse);
  
  // 导航到比较页面
  router.push('/compare');
}

// 根据当前选择的鼠标寻找相似鼠标
async function findSimilarMice() {
  if (!selectedMouse.value || allMice.value.length === 0) {
    similarMice.value = [];
    return;
  }
  
  loading.value = true;
  
  try {
    // 使用本地相似度服务找出相似鼠标
    const results = comparisonService.findSimilarMice(
      selectedMouse.value, 
      allMice.value,
      5
    );
    
    // 为每个鼠标添加相似度得分
    similarMice.value = results.map(mouse => {
      const result = comparisonService.generateComparisonResult([selectedMouse.value as MouseDevice, mouse]);
      return {
        ...mouse,
        similarityScore: result.similarityScore
      };
    });
  } catch (error) {
    console.error('Error finding similar mice:', error);
  } finally {
    loading.value = false;
  }
}

// 获取所有鼠标数据
async function fetchAllMice() {
  loading.value = true;
  try {
    const response = await getDevices({ type: 'mouse' });
    const deviceListResponse = response as unknown as DeviceListResponse;
    allMice.value = deviceListResponse.devices as MouseDevice[];
  } catch (error) {
    console.error('Error fetching mice:', error);
  } finally {
    loading.value = false;
  }
}

// 获取相似度评分颜色
function getSimilarityColor(score: number): string {
  if (score >= 90) return '#67C23A'; // 绿色
  if (score >= 75) return '#409EFF'; // 蓝色
  if (score >= 50) return '#E6A23C'; // 橙色
  return '#F56C6C'; // 红色
}

// 初始化
onMounted(() => {
  fetchAllMice();
});
</script>

<style scoped>
.selected-mouse-card {
  max-width: 200px;
  margin: 0 auto;
}

.similar-mouse-card {
  cursor: pointer;
  transition: all 0.2s ease;
}

.similar-mouse-card:hover {
  transform: translateY(-2px);
  border-color: var(--el-color-primary-light-5);
}

.mouse-image :deep(svg) {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>