<template>
  <div class="mouse-comparison-container">
    <div class="comparison-header">
      <h2 class="text-xl font-bold">鼠标比较</h2>
      <div class="comparison-controls">
        <el-radio-group v-model="comparisonMode" size="small">
          <el-radio-button label="overlay">重叠比较</el-radio-button>
          <el-radio-button label="sideBySide">并排比较</el-radio-button>
        </el-radio-group>
        
        <el-select 
          v-if="comparisonMode === 'overlay'" 
          v-model="viewType" 
          size="small" 
          placeholder="选择视图"
        >
          <el-option label="顶视图" value="topView" />
          <el-option label="侧视图" value="sideView" />
        </el-select>
        
        <el-slider 
          v-if="comparisonMode === 'overlay'" 
          v-model="overlayOpacity" 
          :min="0.2" 
          :max="0.8" 
          :step="0.05" 
          :format-tooltip="value => `透明度: ${Math.round(value * 100)}%`" 
          class="w-32"
        />
      </div>
    </div>
    
    <div v-if="!selectedMice.length" class="empty-state">
      <el-empty description="请选择至少一个鼠标进行比较" />
      <el-button type="primary" @click="openMouseSelector">选择鼠标</el-button>
    </div>
    
    <div v-else class="comparison-content">
      <div class="svg-comparison-area">
        <!-- 比较视图容器 -->
        <div class="svg-container" v-html="comparisonSvg"></div>
        
        <!-- 尺子工具按钮 -->
        <div class="ruler-tools">
          <el-tooltip content="添加拖动尺子" placement="top">
            <el-button 
              type="text" 
              @click="toggleDraggableRuler" 
              :class="{'text-primary': showDraggableRuler}"
            >
              <i class="el-icon-ruler"></i>
            </el-button>
          </el-tooltip>
          <el-tooltip content="显示/隐藏刻度尺" placement="top">
            <el-button 
              type="text" 
              @click="showScaleRuler = !showScaleRuler"
              :class="{'text-primary': showScaleRuler}"
            >
              <i class="el-icon-scale"></i>
            </el-button>
          </el-tooltip>
        </div>
        
        <!-- 可拖动尺子 -->
        <DraggableRuler 
          v-if="showDraggableRuler" 
          :key="draggableRulerKey"
          :initial-position="initialRulerPosition"
          :scale="rulerScale"
        />
        
        <!-- 刻度尺 -->
        <div v-if="showScaleRuler" class="scale-ruler-container">
          <ScaleRuler 
            :width="scaleRulerWidth"
            :markers="mouseMarkers"
            :show-markers="true"
          />
        </div>
        
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-overlay">
          <el-icon><Loading /></el-icon>
        </div>
      </div>
      
      <div class="comparison-details">
        <div>
          <el-button 
            v-if="selectedMice.length < 3" 
            size="small" 
            @click="openMouseSelector" 
            type="primary" 
            plain
          >
            添加鼠标
          </el-button>
        </div>
        
        <div v-if="comparisonData" class="specs-comparison mt-6">
          <h3 class="text-lg font-medium mb-3">参数对比</h3>
          <div class="similarity-score mb-4">
            <div class="font-medium text-sm">相似度评分</div>
            <el-progress 
              :percentage="comparisonData.similarityScore" 
              :color="getSimilarityColor(comparisonData.similarityScore)"
              :format="percent => `${percent}%`"
              :stroke-width="10"
            />
          </div>
          
          <el-table :data="comparisonTableData" border stripe>
            <el-table-column prop="property" label="参数" width="150" />
            <el-table-column 
              v-for="(mouse, index) in selectedMice" 
              :key="mouse.id" 
              :label="`${mouse.brand} ${mouse.name}`"
              :prop="`values[${index}]`"
              :min-width="120" 
            />
          </el-table>
        </div>
      </div>
    </div>
    
    <!-- 使用MouseSelector替换原有对话框 -->
    <el-dialog v-model="mouseDialogVisible" title="选择鼠标" width="60%">
      <MouseSelector 
        :initial-selected-mice="selectedMice"
        :max-selection="3"
        @select="handleMouseSelection"
        @cancel="mouseDialogVisible = false"
      />
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="mouseDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="mouseDialogVisible = false">确认</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useComparisonStore } from '@/stores';
import type { MouseDevice, MouseComparisonResult, ComparisonMode, ViewType } from '@/models/MouseModel';
import svgService from '@/services/svgService';
import comparisonService from '@/services/comparisonService';
import { getDevices } from '@/api/device';
import DraggableRuler from '@/components/tools/DraggableRuler.vue';
import ScaleRuler from '@/components/tools/ScaleRuler.vue';
import MouseSelector from '@/components/comparison/MouseSelector.vue';

// 状态
const comparisonStore = useComparisonStore();
const comparisonMode = ref<ComparisonMode>(comparisonStore.comparisonMode);
const viewType = ref<ViewType>(comparisonStore.viewType);
const overlayOpacity = ref(comparisonStore.overlayOpacity);
const loading = ref(false);
const mouseDialogVisible = ref(false);
const availableMice = ref<MouseDevice[]>([]);
const comparisonData = ref<MouseComparisonResult | null>(null);

// 尺子工具相关状态
const showDraggableRuler = ref(false);
const showScaleRuler = ref(false);
const draggableRulerKey = ref(0);
const initialRulerPosition = ref({ x: 200, y: 100 });
const rulerScale = ref(5); // 像素/mm 比例
const scaleRulerWidth = ref(400);

// 计算属性
const selectedMice = computed(() => comparisonStore.selectedMice);

// 鼠标标记点
const mouseMarkers = computed(() => {
  if (!selectedMice.value.length) return [];
  
  const colors = ['#000000', '#FF0000', '#0000FF'];
  return selectedMice.value.map((mouse, index) => {
    // 标记鼠标的长度位置
    const position = ((mouse.dimensions?.length || 0) / 2) * rulerScale.value;
    return {
      position,
      label: `${mouse.brand} ${mouse.name}`,
      color: colors[index % colors.length]
    };
  });
});

// 比较SVG
const comparisonSvg = computed(() => {
  if (!selectedMice.value.length) return '';
  
  try {
    const svgs = selectedMice.value.map(mouse => {
      const svgData = mouse.svgData?.[viewType.value];
      if (!svgData) return '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text x="10" y="50">No SVG data</text></svg>';
      return svgData;
    });
    
    if (comparisonMode.value === 'overlay') {
      const opacities = selectedMice.value.map((_, index) => 
        index === 0 ? 1.0 : overlayOpacity.value
      );
      const colors = ['#000000', '#FF0000', '#0000FF'];
      return svgService.createOverlaySvg(svgs, opacities, colors.slice(0, selectedMice.value.length));
    } else {
      const labels = selectedMice.value.map(mouse => `${mouse.brand} ${mouse.name}`);
      return svgService.createSideBySideSvg(svgs, labels);
    }
  } catch (error) {
    console.error('Error creating comparison SVG:', error);
    return '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text x="10" y="50">Error rendering comparison</text></svg>';
  }
});

// 对比表格数据
const comparisonTableData = computed(() => {
  if (!comparisonData.value) return [];
  
  return Object.values(comparisonData.value.differences).sort((a, b) => {
    // 按差异百分比降序排序
    return b.differencePercent - a.differencePercent;
  });
});


// 方法
function openMouseSelector() {
  mouseDialogVisible.value = true;
}

function handleMouseSelection(selectedMice) {
  // 通过MouseSelector组件处理选择
  updateComparisonData();
  mouseDialogVisible.value = false;
}

function removeMouse(mouseId: string) {
  comparisonStore.removeMouse(mouseId);
  if (selectedMice.value.length < 2) {
    comparisonData.value = null;
  } else {
    updateComparisonData();
  }
}

async function fetchAvailableMice() {
  loading.value = true;
  try {
    const response = await getDevices({ type: 'mouse' });
    const deviceListResponse = response as unknown as DeviceListResponse;
    availableMice.value = deviceListResponse.devices as MouseDevice[];
  } catch (error) {
    console.error('Error fetching mice:', error);
  } finally {
    loading.value = false;
  }
}

function updateComparisonData() {
  if (selectedMice.value.length < 2) {
    comparisonData.value = null;
    return;
  }
  
  comparisonData.value = comparisonService.generateComparisonResult(selectedMice.value);
}

// 格式化差异百分比
function formatDifference(value: number) {
  return value === 0 ? '相同' : `${value.toFixed(1)}%`;
}

// 获取差异显示样式
function getDifferenceClass(value: number) {
  if (value === 0) return 'text-green-500';
  if (value < 10) return 'text-blue-500';
  if (value < 25) return 'text-amber-500';
  return 'text-red-500';
}

// 获取相似度评分颜色
function getSimilarityColor(score: number) {
  if (score >= 90) return '#67C23A'; // 绿色
  if (score >= 75) return '#409EFF'; // 蓝色
  if (score >= 50) return '#E6A23C'; // 橙色
  return '#F56C6C'; // 红色
}

// 切换可拖动尺子
function toggleDraggableRuler() {
  if (showDraggableRuler.value) {
    showDraggableRuler.value = false;
  } else {
    // 刷新尺子，避免状态残留问题
    draggableRulerKey.value++;
    showDraggableRuler.value = true;
  }
}

// 监听比较模式和透明度变化
watch(comparisonMode, (newMode) => {
  comparisonStore.setComparisonMode(newMode);
});

watch(overlayOpacity, (newOpacity) => {
  comparisonStore.setOverlayOpacity(newOpacity);
});

// 生命周期钩子
onMounted(() => {
  fetchAvailableMice();
  updateComparisonData();
});
</script>

<style scoped>
.mouse-comparison-container {
  padding: 1rem;
}

.comparison-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.comparison-controls {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.empty-state {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 300px;
  gap: 1rem;
}

.comparison-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
}

@media (min-width: 1024px) {
  .comparison-content {
    grid-template-columns: 1fr 1fr;
  }
}

.svg-comparison-area {
  position: relative;
  background-color: #f9f9f9;
  border-radius: 0.5rem;
  padding: 1rem;
  min-height: 400px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.svg-container {
  max-width: 100%;
  height: auto;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
}

.mouse-card.selected {
  border-color: #409EFF;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.ruler-tools {
  position: absolute;
  top: 1rem;
  right: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  background-color: rgba(255, 255, 255, 0.9);
  padding: 0.5rem;
  border-radius: 0.25rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.text-primary {
  color: #409EFF;
}

.scale-ruler-container {
  position: absolute;
  bottom: 1rem;
  left: 1rem;
  right: 1rem;
  z-index: 20;
}
</style>