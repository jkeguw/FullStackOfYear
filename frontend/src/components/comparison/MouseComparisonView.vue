<template>
  <div class="mouse-comparison-container">
    <div class="comparison-header">
      <h2 class="text-xl font-bold">Mouse Comparison</h2>
      <div class="comparison-controls">
        <el-radio-group v-model="comparisonMode" size="small">
          <el-radio-button label="overlay">Overlay Comparison</el-radio-button>
          <el-radio-button label="sideBySide">Side-by-side Comparison</el-radio-button>
        </el-radio-group>

        <el-select
          v-if="comparisonMode === 'overlay'"
          v-model="viewType"
          size="small"
          placeholder="Select View"
        >
          <el-option label="Top View" value="topView" />
          <el-option label="Side View" value="sideView" />
        </el-select>

        <el-slider
          v-if="comparisonMode === 'overlay'"
          v-model="overlayOpacity"
          :min="0.2"
          :max="0.8"
          :step="0.05"
          :format-tooltip="(value) => `Transparency: ${Math.round(value * 100)}%`"
          class="w-32"
        />
      </div>
    </div>

    <div v-if="!selectedMice.length" class="empty-state">
      <el-empty description="Please select at least one mouse to compare" />
      <el-button type="primary" @click="openMouseSelector">Select Mouse</el-button>
    </div>

    <div v-else class="comparison-content">
      <div class="svg-comparison-area">
        <!-- 比较视图容器 -->
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div class="svg-container" v-html="comparisonSvg"></div>

        <!-- 尺子工具按钮 -->
        <div class="ruler-tools">
          <el-tooltip content="添加拖动尺子" placement="top">
            <el-button
              type="text"
              :class="{ 'text-primary': showDraggableRuler }"
              @click="toggleDraggableRuler"
            >
              <i class="el-icon-ruler"></i>
            </el-button>
          </el-tooltip>
          <el-tooltip content="显示/隐藏刻度尺" placement="top">
            <el-button
              type="text"
              :class="{ 'text-primary': showScaleRuler }"
              @click="showScaleRuler = !showScaleRuler"
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
          <ScaleRuler :width="scaleRulerWidth" :markers="mouseMarkers" :show-markers="true" />
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
            type="primary"
            plain
            @click="openMouseSelector"
          >
            Add Mouse
          </el-button>
        </div>

        <div v-if="comparisonData" class="specs-comparison mt-6">
          <h3 class="text-lg font-medium mb-3">Parameter Comparison</h3>
          <div class="similarity-score mb-4">
            <div class="font-medium text-sm">Similarity Score</div>
            <el-progress
              :percentage="comparisonData.similarityScore"
              :color="getSimilarityColor(comparisonData.similarityScore)"
              :format="(percent) => `${percent}%`"
              :stroke-width="10"
            />
          </div>

          <el-table :data="comparisonTableData" border stripe>
            <el-table-column prop="property" label="Parameter" width="150" />
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
    <el-dialog v-model="mouseDialogVisible" title="Select Mouse" width="60%">
      <MouseSelector
        :initial-selected-mice="selectedMice"
        :max-selection="3"
        @select="handleMouseSelection"
        @cancel="mouseDialogVisible = false"
      />

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="mouseDialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="handleDialogConfirm">Confirm</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { DeviceListResponse } from '@/api/device';
import { ref, computed, onMounted, watch } from 'vue';
import { useComparisonStore } from '@/stores';
import type {
  MouseDevice,
  MouseComparisonResult,
  ComparisonMode,
  ViewType
} from '@/models/MouseModel';
import svgService from '@/services/svgService';
import comparisonService from '@/services/comparisonService';
import { getDevices, getMouseSVG } from '@/api/device';
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

// 不再需要单独加载SVG数据函数，现在直接通过API获取比较结果
// 此函数保留为备用但不再使用
const loadSvgData = async () => {
  console.warn('loadSvgData is deprecated, using API comparison instead');
  if (!selectedMice.value.length) return [];
  
  loading.value = true;
  
  try {
    const svgPromises = selectedMice.value.map(async (mouse) => {
      // 从API获取SVG数据
      try {
        const result = await getMouseSVG(mouse.id, viewType.value as 'top' | 'side');
        if (result && result.data) {
          return result.data.svgData;
        }
      } catch (err) {
        console.error(`Error fetching SVG for mouse ${mouse.id}:`, err);
      }
      
      // 如果获取失败，返回占位SVG
      console.warn(`鼠标 ${mouse.brand} ${mouse.name} 缺少 ${viewType.value} 视图SVG数据`);
      return '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text x="50" y="50" text-anchor="middle" fill="#999">暂无SVG数据</text></svg>';
    });
    
    return await Promise.all(svgPromises);
  } catch (error) {
    console.error('Error loading SVG data:', error);
    return [];
  } finally {
    loading.value = false;
  }
};

// 比较SVG
const comparisonSvg = ref('');

// 更新比较SVG
const updateComparisonSvg = async () => {
  if (!selectedMice.value.length) {
    comparisonSvg.value = '';
    return;
  }
  
  try {
    // 用新的后端API进行SVG比较
    const deviceIds = selectedMice.value.map(mouse => mouse.id);
    const view = viewType.value === 'topView' ? 'top' : 'side';
    
    if (comparisonMode.value === 'overlay') {
      const opacities = selectedMice.value.map((_, index) =>
        index === 0 ? 1.0 : overlayOpacity.value
      );
      const colors = ['#000000', '#FF0000', '#0000FF'];
      
      try {
        // 使用API进行SVG叠加比较
        const result = await svgService.createOverlaySvg(
          deviceIds,
          view,
          opacities,
          colors.slice(0, selectedMice.value.length)
        );
        comparisonSvg.value = result;
      } catch (apiError) {
        console.warn('使用API获取SVG失败，尝试使用本地SVG文件:', apiError);
        // 回退到本地SVG文件
        const svgsPromises = selectedMice.value.map(mouse => {
          const fileName = mouse.id.replace(/-/g, ' ');
          return fetch(`/svg/${fileName} ${view}.svg`)
            .then(response => {
              if (!response.ok) {
                throw new Error(`无法加载SVG: ${response.status}`);
              }
              return response.text();
            })
            .catch(err => {
              console.error(`加载SVG失败: ${mouse.id}`, err);
              return '';
            });
        });
        
        const svgs = await Promise.all(svgsPromises);
        if (svgs.some(svg => svg)) {
          const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
          svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
          svgContainer.setAttribute('viewBox', '0 0 1250 400');
          
          svgs.forEach((svgContent, index) => {
            if (!svgContent) return;
            
            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = svgContent;
            const svgElement = tempDiv.querySelector('svg');
            const pathElement = svgElement?.querySelector('path');
            
            if (pathElement) {
              const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
              g.setAttribute('opacity', opacities[index].toString());
              g.setAttribute('fill', 'none');
              g.setAttribute('stroke', colors[index]);
              g.setAttribute('stroke-width', '2');
              
              const clonedPath = pathElement.cloneNode(true);
              g.appendChild(clonedPath);
              svgContainer.appendChild(g);
            }
          });
          
          comparisonSvg.value = svgContainer.outerHTML;
        } else {
          throw new Error('无法加载SVG文件');
        }
      }
    } else {
      try {
        // 使用API进行SVG并排比较
        const result = await svgService.createSideBySideSvg(
          deviceIds,
          view
        );
        comparisonSvg.value = result;
      } catch (apiError) {
        console.warn('使用API获取并排SVG失败，尝试使用本地SVG文件:', apiError);
        // 回退到本地SVG文件
        const svgsPromises = selectedMice.value.map(mouse => {
          const fileName = mouse.id.replace(/-/g, ' ');
          return fetch(`/svg/${fileName} ${view}.svg`)
            .then(response => {
              if (!response.ok) {
                throw new Error(`无法加载SVG: ${response.status}`);
              }
              return response.text();
            })
            .catch(err => {
              console.error(`加载SVG失败: ${mouse.id}`, err);
              return '';
            });
        });
        
        const svgs = await Promise.all(svgsPromises);
        if (svgs.some(svg => svg)) {
          const width = 1250;
          const height = 400;
          const padding = 20;
          const mouseWidth = (width - (padding * (selectedMice.value.length + 1))) / selectedMice.value.length;
          
          const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
          svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
          svgContainer.setAttribute('viewBox', `0 0 ${width} ${height}`);
          
          svgs.forEach((svgContent, index) => {
            if (!svgContent) return;
            
            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = svgContent;
            const svgElement = tempDiv.querySelector('svg');
            const pathElement = svgElement?.querySelector('path');
            
            if (pathElement) {
              const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
              g.setAttribute('transform', `translate(${padding + (mouseWidth + padding) * index}, 0) scale(${mouseWidth/width})`);
              g.setAttribute('fill', 'none');
              g.setAttribute('stroke', '#000');
              g.setAttribute('stroke-width', '2');
              
              const clonedPath = pathElement.cloneNode(true);
              g.appendChild(clonedPath);
              
              // 添加鼠标名称
              const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
              text.setAttribute('x', (mouseWidth / 2).toString());
              text.setAttribute('y', (height - 20).toString());
              text.setAttribute('text-anchor', 'middle');
              text.setAttribute('font-family', 'Arial, sans-serif');
              text.setAttribute('font-size', '14');
              text.textContent = `${selectedMice.value[index].brand} ${selectedMice.value[index].name}`;
              g.appendChild(text);
              
              svgContainer.appendChild(g);
            }
          });
          
          comparisonSvg.value = svgContainer.outerHTML;
        } else {
          throw new Error('无法加载SVG文件');
        }
      }
    }
  } catch (error) {
    console.error('Error creating comparison SVG:', error);
    comparisonSvg.value = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text x="50" y="50" text-anchor="middle" fill="#f56c6c">渲染SVG比较图失败</text></svg>';
  }
};

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

function handleMouseSelection(_selectedMice) {
  // 通过MouseSelector组件处理选择
  updateComparisonData();
  updateComparisonSvg();
}

function handleDialogConfirm() {
  updateComparisonData();
  updateComparisonSvg();
  mouseDialogVisible.value = false;
}

function _removeMouse(mouseId: string) {
  comparisonStore.removeMouse(mouseId);
  if (selectedMice.value.length < 2) {
    comparisonData.value = null;
  } else {
    updateComparisonData();
  }
  updateComparisonSvg();
}

async function fetchAvailableMice() {
  loading.value = true;
  try {
    const response = await getDevices({ type: 'mouse' });
    if (!response || !response.data) {
      throw new Error('获取鼠标数据失败: 没有返回数据');
    }
    availableMice.value = response.data.devices as MouseDevice[];
  } catch (error) {
    console.error('Error fetching mice:', error);
    availableMice.value = []; // 设置为空数组，避免引用空值
  } finally {
    loading.value = false;
  }
}

function updateComparisonData() {
  if (selectedMice.value.length < 2) {
    comparisonData.value = null;
    return;
  }

  // @ts-expect-error - Type inconsistency between MouseDevice definitions
  comparisonData.value = comparisonService.generateComparisonResult(selectedMice.value);
}

// 格式化差异百分比
function _formatDifference(value: number) {
  return value === 0 ? '相同' : `${value.toFixed(1)}%`;
}

// 获取差异显示样式
function _getDifferenceClass(value: number) {
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
  updateComparisonSvg();
});

watch(overlayOpacity, (newOpacity) => {
  comparisonStore.setOverlayOpacity(newOpacity);
  updateComparisonSvg();
});

watch(viewType, () => {
  updateComparisonSvg();
});

watch(selectedMice, () => {
  updateComparisonSvg();
}, { deep: true });

// 生命周期钩子
onMounted(() => {
  fetchAvailableMice();
  updateComparisonData();
  updateComparisonSvg();
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
  width: 100%;
  background-color: #ffffff;
  border-radius: 8px;
  padding: 10px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
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
  border-color: #409eff;
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
  color: #409eff;
}

.scale-ruler-container {
  position: absolute;
  bottom: 1rem;
  left: 1rem;
  right: 1rem;
  z-index: 20;
}
</style>
