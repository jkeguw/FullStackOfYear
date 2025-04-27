<template>
  <div class="mouse-detail-page">
    <div class="container mx-auto py-6 px-4">
      <div v-if="loading" class="loading-container flex justify-center items-center py-20">
        <el-skeleton style="width: 100%" animated>
          <template #template>
            <div class="flex flex-col md:flex-row gap-8">
              <div class="md:w-1/2">
                <el-skeleton-item variant="image" style="width: 100%; height: 400px" />
              </div>
              <div class="md:w-1/2">
                <el-skeleton-item variant="h1" style="width: 50%" />
                <div class="my-4">
                  <el-skeleton-item variant="text" style="margin-right: 16px" />
                  <el-skeleton-item variant="text" style="width: 30%" />
                </div>
                <el-skeleton :rows="6" />
              </div>
            </div>
          </template>
        </el-skeleton>
      </div>

      <div v-else-if="mouse" class="mouse-detail">
        <el-breadcrumb separator="/" class="mb-6">
          <el-breadcrumb-item :to="{ name: 'Home' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ name: 'DeviceList' }">鼠标数据库</el-breadcrumb-item>
          <el-breadcrumb-item>{{ mouse.brand }} {{ mouse.name }}</el-breadcrumb-item>
        </el-breadcrumb>
        
        <div class="flex flex-col md:flex-row gap-8">
          <!-- 左侧：图片/SVG区域 -->
          <div class="md:w-1/2">
            <div class="mouse-image-container bg-white rounded-lg p-6 mb-4 shadow-sm">
              <div class="relative aspect-video mb-2">
                <div v-if="activeImage.type === 'svg'" class="w-full h-full" v-html="activeImage.content"></div>
                <img v-else :src="activeImage.content" class="w-full h-full object-contain" />
              </div>
              
              <div class="flex gap-2 mt-4">
                <el-button-group>
                  <el-button 
                    :type="viewType === 'topView' ? 'primary' : 'default'" 
                    @click="viewType = 'topView'"
                    size="small"
                  >
                    顶视图
                  </el-button>
                  <el-button 
                    :type="viewType === 'sideView' ? 'primary' : 'default'" 
                    @click="viewType = 'sideView'"
                    size="small"
                  >
                    侧视图
                  </el-button>
                </el-button-group>
                
                <div class="flex-grow"></div>
                
                <el-button 
                  type="success" 
                  size="small"
                  @click="addToComparison"
                >
                  <template #icon><Plus /></template>
                  添加到比较
                </el-button>
                
                <add-to-cart-button
                  v-if="mouse && mouse.price"
                  :product="{ 
                    id: mouse.id, 
                    name: mouse.brand + ' ' + mouse.name, 
                    price: mouse.price,
                    image: mouse.imageUrl
                  }"
                  size="small"
                  type="primary"
                />
              </div>
            </div>
            
            <div class="mouse-dimensions bg-white rounded-lg p-6 shadow-sm">
              <h3 class="text-lg font-medium mb-4">尺寸信息</h3>
              <div class="dimensions-grid">
                <div class="dimension-item">
                  <div class="dimension-label">长度</div>
                  <div class="dimension-value">{{ mouse.dimensions.length }} mm</div>
                </div>
                <div class="dimension-item">
                  <div class="dimension-label">宽度</div>
                  <div class="dimension-value">{{ mouse.dimensions.width }} mm</div>
                </div>
                <div class="dimension-item">
                  <div class="dimension-label">高度</div>
                  <div class="dimension-value">{{ mouse.dimensions.height }} mm</div>
                </div>
                <div class="dimension-item">
                  <div class="dimension-label">重量</div>
                  <div class="dimension-value">{{ mouse.dimensions.weight }} g</div>
                </div>
                <div class="dimension-item" v-if="mouse.dimensions.gripWidth">
                  <div class="dimension-label">握宽</div>
                  <div class="dimension-value">{{ mouse.dimensions.gripWidth }} mm</div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 右侧：鼠标信息 -->
          <div class="md:w-1/2">
            <div class="bg-white rounded-lg p-6 shadow-sm">
              <div class="flex items-center mb-6">
                <h1 class="text-2xl font-bold">{{ mouse.brand }} {{ mouse.name }}</h1>
                <div class="flex-grow"></div>
                <div v-if="isInComparison" class="badge bg-success-100 text-success-700 px-2 py-1 rounded">
                  已加入比较
                </div>
              </div>
              
              <div class="mb-6" v-if="mouse.description">
                <h3 class="text-lg font-medium mb-2">产品简介</h3>
                <p class="text-gray-700">{{ mouse.description }}</p>
              </div>
              
              <el-tabs>
                <el-tab-pane label="技术规格">
                  <div class="specs-table">
                    <div class="spec-group">
                      <h4 class="text-base font-medium mb-2">形状特征</h4>
                      <div class="grid grid-cols-2 gap-4">
                        <div class="spec-item">
                          <div class="spec-label">形状类型</div>
                          <div class="spec-value">{{ formatShapeType(mouse.shape.type) }}</div>
                        </div>
                        <div class="spec-item">
                          <div class="spec-label">凸起位置</div>
                          <div class="spec-value">{{ formatHumpPlacement(mouse.shape.humpPlacement) }}</div>
                        </div>
                        <div class="spec-item">
                          <div class="spec-label">前部扩张</div>
                          <div class="spec-value">{{ mouse.shape.frontFlare }}</div>
                        </div>
                        <div class="spec-item">
                          <div class="spec-label">侧面曲率</div>
                          <div class="spec-value">{{ mouse.shape.sideCurvature }}</div>
                        </div>
                        <div class="spec-item">
                          <div class="spec-label">手型适配</div>
                          <div class="spec-value">{{ formatHandCompatibility(mouse.shape.handCompatibility) }}</div>
                        </div>
                        <div class="spec-item" v-if="mouse.shape.thumbRest !== undefined">
                          <div class="spec-label">拇指区</div>
                          <div class="spec-value">{{ mouse.shape.thumbRest ? '有拇指托' : '无拇指托' }}</div>
                        </div>
                      </div>
                    </div>
                    
                    <div class="spec-group mt-6">
                      <h4 class="text-base font-medium mb-2">技术参数</h4>
                      <div class="grid grid-cols-2 gap-4">
                        <div class="spec-item">
                          <div class="spec-label">连接方式</div>
                          <div class="spec-value">{{ formatConnectivity(mouse.technical.connectivity) }}</div>
                        </div>
                        <div class="spec-item">
                          <div class="spec-label">传感器</div>
                          <div class="spec-value">{{ mouse.technical.sensor }}</div>
                        </div>
                        <div class="spec-item">
                          <div class="spec-label">最大DPI</div>
                          <div class="spec-value">{{ formatNumber(mouse.technical.maxDPI) }}</div>
                        </div>
                        <div class="spec-item">
                          <div class="spec-label">轮询率</div>
                          <div class="spec-value">{{ mouse.technical.pollingRate }} Hz</div>
                        </div>
                        <div class="spec-item">
                          <div class="spec-label">侧键数量</div>
                          <div class="spec-value">{{ mouse.technical.sideButtons }}</div>
                        </div>
                        <div class="spec-item" v-if="mouse.technical.battery">
                          <div class="spec-label">电池寿命</div>
                          <div class="spec-value">{{ mouse.technical.battery.life }} 小时</div>
                        </div>
                      </div>
                    </div>
                  </div>
                </el-tab-pane>
                
                <el-tab-pane label="推荐用途">
                  <div class="recommendations">
                    <div class="mb-4">
                      <h4 class="text-base font-medium mb-2">适合游戏类型</h4>
                      <div class="flex gap-2 flex-wrap">
                        <el-tag v-for="game in mouse.recommended.gameTypes" :key="game" size="small">
                          {{ game }}
                        </el-tag>
                      </div>
                    </div>
                    
                    <div class="mb-4">
                      <h4 class="text-base font-medium mb-2">适合握持方式</h4>
                      <div class="flex gap-2 flex-wrap">
                        <el-tag 
                          v-for="grip in mouse.recommended.gripStyles" 
                          :key="grip"
                          size="small"
                          type="success"
                        >
                          {{ formatGripStyle(grip) }}
                        </el-tag>
                      </div>
                    </div>
                    
                    <div class="mb-4">
                      <h4 class="text-base font-medium mb-2">适合手型尺寸</h4>
                      <div class="flex gap-2 flex-wrap">
                        <el-tag 
                          v-for="hand in mouse.recommended.handSizes" 
                          :key="hand"
                          size="small"
                          type="warning"
                        >
                          {{ formatHandSize(hand) }}
                        </el-tag>
                      </div>
                    </div>
                    
                    <div class="other-recommendations grid grid-cols-2 gap-4 mt-6">
                      <div class="recommendation-item">
                        <h4 class="text-base font-medium mb-2">日常使用</h4>
                        <div class="flex items-center">
                          <el-icon v-if="mouse.recommended.dailyUse" class="text-success" size="large">
                            <CircleCheckFilled />
                          </el-icon>
                          <el-icon v-else class="text-danger" size="large">
                            <CircleCloseFilled />
                          </el-icon>
                          <span class="ml-2">{{ mouse.recommended.dailyUse ? '适合' : '不适合' }}</span>
                        </div>
                      </div>
                      
                      <div class="recommendation-item">
                        <h4 class="text-base font-medium mb-2">专业竞技</h4>
                        <div class="flex items-center">
                          <el-icon v-if="mouse.recommended.professional" class="text-success" size="large">
                            <CircleCheckFilled />
                          </el-icon>
                          <el-icon v-else class="text-danger" size="large">
                            <CircleCloseFilled />
                          </el-icon>
                          <span class="ml-2">{{ mouse.recommended.professional ? '适合' : '不适合' }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </el-tab-pane>
                
                <el-tab-pane label="相似鼠标">
                  <div class="similar-mice">
                    <div v-if="similarMice.length === 0" class="text-center py-8">
                      <el-icon class="text-3xl mb-3"><Loading /></el-icon>
                      <p>正在加载相似鼠标推荐...</p>
                    </div>
                    
                    <div v-else>
                      <p class="mb-4 text-gray-600">
                        根据形状和技术特性，以下鼠标与 {{ mouse.brand }} {{ mouse.name }} 比较相似：
                      </p>
                      
                      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <el-card 
                          v-for="similarMouse in similarMice" 
                          :key="similarMouse.id"
                          class="similar-mouse-card"
                          shadow="hover"
                          @click="viewSimilarMouse(similarMouse.id)"
                        >
                          <div class="flex items-center">
                            <div class="mouse-image w-16 h-16 mr-3 flex-shrink-0">
                              <template v-if="similarMouse.imageUrl">
                                <img :src="similarMouse.imageUrl" class="w-full h-full object-contain" :alt="similarMouse.name" />
                              </template>
                              <template v-else-if="similarMouse.svgData?.topView">
                                <div class="w-full h-full" v-html="similarMouse.svgData.topView"></div>
                              </template>
                              <template v-else>
                                <div class="w-full h-full bg-gray-100 flex items-center justify-center">
                                  <el-icon class="text-2xl text-gray-300"><Mouse /></el-icon>
                                </div>
                              </template>
                            </div>
                            <div>
                              <h3 class="text-base font-medium">{{ similarMouse.brand }} {{ similarMouse.name }}</h3>
                              <p class="text-xs text-gray-500 mt-1">
                                {{ similarMouse.dimensions.length }}×{{ similarMouse.dimensions.width }}×{{ similarMouse.dimensions.height }}mm
                              </p>
                              <p class="text-xs text-gray-500">
                                {{ similarMouse.dimensions.weight }}g · 
                                {{ formatShapeType(similarMouse.shape.type) }}
                              </p>
                            </div>
                          </div>
                        </el-card>
                      </div>
                      
                      <div class="text-center mt-6">
                        <el-button 
                          type="primary" 
                          @click="compareWithSimilar"
                          :disabled="!canCompareWithSimilar"
                        >
                          与相似鼠标比较
                        </el-button>
                      </div>
                    </div>
                  </div>
                </el-tab-pane>
              </el-tabs>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="error-container text-center py-20">
        <el-result
          icon="warning"
          title="未找到鼠标数据"
          sub-title="抱歉，我们找不到您请求的鼠标信息"
        >
          <template #extra>
            <el-button type="primary" @click="$router.push({ name: 'DeviceList' })">返回鼠标列表</el-button>
          </template>
        </el-result>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useComparisonStore } from '@/stores';
import type { MouseDevice } from '@/models/MouseModel';
import { ElMessage } from 'element-plus';
import { getDevices } from '@/api/device';
import { testMice } from '@/data/testMice';
import comparisonService from '@/services/comparisonService';
import { Plus, Loading, CircleCheckFilled, CircleCloseFilled, Mouse } from '@element-plus/icons-vue';
import AddToCartButton from '@/components/cart/AddToCartButton.vue';

// 路由
const route = useRoute();
const router = useRouter();

// 状态管理
const comparisonStore = useComparisonStore();

// 状态
const loading = ref(true);
const mouse = ref<MouseDevice | null>(null);
const viewType = ref<'topView' | 'sideView'>('topView');
const similarMice = ref<MouseDevice[]>([]);
const allMice = ref<MouseDevice[]>([]);

// 计算属性
const mouseId = computed(() => route.params.id as string);

const activeImage = computed(() => {
  if (!mouse.value || !mouse.value.svgData) return { type: 'placeholder', content: '' };
  
  const svgContent = mouse.value.svgData[viewType.value];
  if (svgContent) {
    return { type: 'svg', content: svgContent };
  } else if (mouse.value.imageUrl) {
    return { type: 'image', content: mouse.value.imageUrl };
  }
  
  return { type: 'placeholder', content: '' };
});

const isInComparison = computed(() => {
  return comparisonStore.selectedMice.some(m => m.id === mouseId.value);
});

const canCompareWithSimilar = computed(() => {
  return similarMice.value.length > 0;
});

// 方法
async function fetchMouseData() {
  loading.value = true;
  try {
    // 从测试数据中获取鼠标信息
    const selectedMouse = testMice.find(m => m.id === mouseId.value);
    if (selectedMouse) {
      mouse.value = selectedMouse;
      comparisonStore.addToRecentlyViewed(selectedMouse);
    } else {
      mouse.value = null;
    }
  } catch (error) {
    console.error('Error fetching mouse data:', error);
    mouse.value = null;
  } finally {
    loading.value = false;
  }
}

async function fetchAllMice() {
  try {
    allMice.value = testMice;
    findSimilarMice();
  } catch (error) {
    console.error('Error fetching all mice:', error);
  }
}

function findSimilarMice() {
  if (!mouse.value || allMice.value.length === 0) {
    similarMice.value = [];
    return;
  }
  
  // 使用比较服务找出相似鼠标
  similarMice.value = comparisonService.findSimilarMice(
    mouse.value,
    allMice.value,
    3
  );
}

function addToComparison() {
  if (!mouse.value) return;
  
  comparisonStore.addMouse(mouse.value);
  ElMessage.success(`已将 ${mouse.value.brand} ${mouse.value.name} 添加到比较`);
}

function viewSimilarMouse(id: string) {
  router.push({ name: 'MouseDetail', params: { id } });
}

function compareWithSimilar() {
  if (!mouse.value) return;
  
  // 清空当前比较
  comparisonStore.clearSelection();
  
  // 添加当前鼠标
  comparisonStore.addMouse(mouse.value);
  
  // 添加第一个相似鼠标
  if (similarMice.value.length > 0) {
    comparisonStore.addMouse(similarMice.value[0]);
  }
  
  // 导航到比较页面
  router.push({ name: 'Compare' });
}

// 格式化工具函数
function formatShapeType(type: string): string {
  const shapeMap: Record<string, string> = {
    'ergo': '人体工学',
    'symmetrical': '对称',
    'asymmetrical': '非对称',
    'fingertip': '指尖式',
    'ambi': '双手通用'
  };
  
  return shapeMap[type] || type;
}

function formatHandCompatibility(compatibility: string): string {
  const compatMap: Record<string, string> = {
    'right': '右手',
    'left': '左手',
    'ambidextrous': '双手通用'
  };
  
  return compatMap[compatibility] || compatibility;
}

function formatHumpPlacement(placement: string): string {
  const placementMap: Record<string, string> = {
    'front': '前部',
    'center': '中部',
    'back': '后部',
    'none': '无明显凸起'
  };
  
  return placementMap[placement] || placement;
}

function formatGripStyle(style: string): string {
  const styleMap: Record<string, string> = {
    'palm': '掌握',
    'claw': '爪握',
    'fingertip': '指尖',
    'hybrid': '混合'
  };
  
  return styleMap[style] || style;
}

function formatHandSize(size: string): string {
  const sizeMap: Record<string, string> = {
    'small': '小型手',
    'medium': '中型手',
    'large': '大型手',
    'extra-large': '超大型手'
  };
  
  return sizeMap[size] || size;
}

function formatConnectivity(connectivity: string[]): string {
  if (!connectivity || connectivity.length === 0) return '未知';
  
  const connectivityMap: Record<string, string> = {
    'wired': '有线',
    'wireless': '无线',
    'bluetooth': '蓝牙',
    'dual': '双模'
  };
  
  return connectivity.map(c => connectivityMap[c] || c).join('，');
}

function formatNumber(num: number): string {
  return num.toLocaleString();
}

// 生命周期
onMounted(() => {
  fetchMouseData();
  fetchAllMice();
});

// 监听路由变化，以便在导航到不同的鼠标时更新数据
watch(() => route.params.id, () => {
  fetchMouseData();
  findSimilarMice();
});
</script>

<style scoped>
.mouse-detail-page {
  min-height: 100vh;
  background-color: var(--claude-gray-100);
}

.dimensions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.dimension-item {
  padding: 0.5rem;
  border-radius: 0.25rem;
  background-color: var(--claude-gray-50);
  border: 1px solid var(--claude-gray-200);
}

.dimension-label {
  font-size: 0.875rem;
  color: var(--claude-gray-600);
  margin-bottom: 0.25rem;
}

.dimension-value {
  font-size: 1rem;
  font-weight: 500;
  color: var(--claude-gray-900);
}

.specs-table {
  margin-top: 1rem;
}

.spec-item {
  padding: 0.5rem;
  background-color: var(--claude-gray-50);
  border-radius: 0.25rem;
  margin-bottom: 0.5rem;
}

.spec-label {
  font-size: 0.75rem;
  color: var(--claude-gray-600);
  margin-bottom: 0.25rem;
}

.spec-value {
  font-size: 0.875rem;
  color: var(--claude-gray-900);
}

.similar-mouse-card {
  cursor: pointer;
  transition: all 0.2s ease;
}

.similar-mouse-card:hover {
  transform: translateY(-2px);
  border-color: var(--claude-primary-light-5);
}

.badge {
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: 0.25rem;
}

.text-success {
  color: var(--success-color, #10b981);
}

.text-danger {
  color: var(--error-color, #ef4444);
}
</style>