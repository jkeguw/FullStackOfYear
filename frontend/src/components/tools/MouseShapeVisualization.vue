<template>
  <div class="mouse-shape-visualization">
    <el-card>
      <template #header>
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center space-y-3 sm:space-y-0">
          <h3 class="text-lg font-medium">鼠标形状可视化工具</h3>
          <div class="flex flex-wrap gap-2 w-full sm:w-auto">
            <el-select v-model="selectedMouse" size="small" class="w-full xs:w-auto">
              <el-option 
                v-for="mouse in availableMice" 
                :key="mouse.id" 
                :label="mouse.name" 
                :value="mouse.id"
              ></el-option>
            </el-select>
            <el-select v-model="unit" size="small" class="w-full xs:w-auto">
              <el-option label="厘米 (cm)" value="cm"></el-option>
              <el-option label="毫米 (mm)" value="mm"></el-option>
              <el-option label="英寸 (inch)" value="inch"></el-option>
            </el-select>
          </div>
        </div>
      </template>
      
      <!-- 校准对话框 -->
      <el-dialog
        v-model="calibrateDialogVisible"
        :title="calibrationStep === 1 ? '校准尺子 - 步骤1：标记位置' : '校准尺子 - 步骤2：输入实际长度'"
        width="400px"
      >
        <template v-if="calibrationStep === 1">
          <p>请在尺子上点击一个已知长度的位置，然后点击"下一步"。</p>
        </template>
        <template v-else>
          <p>请输入您标记的位置对应的实际长度：</p>
          <el-input 
            v-model="calibrationValue" 
            placeholder="请输入长度值" 
            type="number" 
            :min="0"
            :step="0.1"
          >
            <template #append>{{ unit }}</template>
          </el-input>
        </template>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="calibrateDialogVisible = false">取消</el-button>
            <el-button v-if="calibrationStep === 1" type="primary" @click="handleCalibrateConfirm">
              下一步
            </el-button>
            <el-button v-else type="primary" @click="handleCalibrateConfirm">
              确认
            </el-button>
          </span>
        </template>
      </el-dialog>

      <div class="flex flex-col space-y-6">
        <!-- Mouse visualization tabs -->
        <el-tabs v-model="activeView" stretch>
          <el-tab-pane label="顶视图" name="top">
            <div class="svg-container relative" ref="svgContainer">
              <!-- Zoom controls -->
              <div class="zoom-controls">
                <el-button-group>
                  <el-button size="small" icon="el-icon-zoom-in" @click="zoomIn">+</el-button>
                  <el-button size="small" icon="el-icon-zoom-out" @click="zoomOut">-</el-button>
                  <el-button size="small" icon="el-icon-refresh" @click="resetZoom">重置</el-button>
                </el-button-group>
              </div>
              
              <!-- SVG container with zoom and pan support -->
              <div class="svg-wrapper" 
                :style="{ transform: `scale(${zoomLevel})`, 
                        transformOrigin: 'center center' }"
                @mousedown="startPan"
                @touchstart="handleTouchStart"
                @mousemove="pan"
                @touchmove="handleTouchMove"
                @mouseup="endPan"
                @touchend="handleTouchEnd"
                @mouseleave="endPan">
                
                <!-- Dynamic SVG import with Suspense for lazy loading -->
                <Suspense>
                  <template #default>
                    <img 
                      v-if="svgLoaded" 
                      :src="currentMouseSvg.top" 
                      class="mouse-svg top-view" 
                      @load="onSvgLoad"
                      alt="Mouse top view"
                    />
                  </template>
                  <template #fallback>
                    <div class="loading-placeholder">
                      <el-skeleton :rows="3" animated />
                      <p class="text-center text-gray-500">Loading SVG...</p>
                    </div>
                  </template>
                </Suspense>
                
                <!-- Key feature annotations -->
                <div v-if="svgLoaded" class="feature-annotations">
                  <div v-for="(annotation, index) in currentAnnotations.top" 
                       :key="index" 
                       class="annotation"
                       :style="{ left: `${annotation.x}%`, top: `${annotation.y}%` }">
                    <div class="annotation-marker"></div>
                    <div class="annotation-text">{{ annotation.text }}</div>
                  </div>
                </div>

                <!-- Dimension indicators for top view -->
                <div class="dimension-indicator length" :style="{ width: `${currentMouseDimensions.length}px` }">
                  <span class="dimension-value">{{ formatDimension(currentMouseDimensions.realLength) }}</span>
                </div>
                <div class="dimension-indicator width" :style="{ width: `${currentMouseDimensions.width}px` }">
                  <span class="dimension-value">{{ formatDimension(currentMouseDimensions.realWidth) }}</span>
                </div>
                <div class="dimension-indicator grip-width" :style="{ width: `${currentMouseDimensions.gripWidth}px` }">
                  <span class="dimension-value">{{ formatDimension(currentMouseDimensions.realGripWidth) }}</span>
                </div>
              </div>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="侧视图" name="side">
            <div class="svg-container relative" ref="svgContainerSide">
              <!-- Zoom controls -->
              <div class="zoom-controls">
                <el-button-group>
                  <el-button size="small" icon="el-icon-zoom-in" @click="zoomIn">+</el-button>
                  <el-button size="small" icon="el-icon-zoom-out" @click="zoomOut">-</el-button>
                  <el-button size="small" icon="el-icon-refresh" @click="resetZoom">重置</el-button>
                </el-button-group>
              </div>
              
              <!-- SVG container with zoom and pan support -->
              <div class="svg-wrapper" 
                :style="{ transform: `scale(${zoomLevel})`, 
                        transformOrigin: 'center center' }"
                @mousedown="startPan"
                @touchstart="handleTouchStart"
                @mousemove="pan"
                @touchmove="handleTouchMove"
                @mouseup="endPan"
                @touchend="handleTouchEnd"
                @mouseleave="endPan">
                
                <!-- Dynamic SVG import with Suspense for lazy loading -->
                <Suspense>
                  <template #default>
                    <img 
                      v-if="svgLoaded" 
                      :src="currentMouseSvg.side" 
                      class="mouse-svg side-view" 
                      @load="onSvgLoad"
                      alt="Mouse side view"
                    />
                  </template>
                  <template #fallback>
                    <div class="loading-placeholder">
                      <el-skeleton :rows="3" animated />
                      <p class="text-center text-gray-500">Loading SVG...</p>
                    </div>
                  </template>
                </Suspense>
                
                <!-- Key feature annotations -->
                <div v-if="svgLoaded" class="feature-annotations">
                  <div v-for="(annotation, index) in currentAnnotations.side" 
                       :key="index" 
                       class="annotation"
                       :style="{ left: `${annotation.x}%`, top: `${annotation.y}%` }">
                    <div class="annotation-marker"></div>
                    <div class="annotation-text">{{ annotation.text }}</div>
                  </div>
                </div>

                <!-- Dimension indicators for side view -->
                <div class="dimension-indicator height" :style="{ height: `${currentMouseDimensions.height}px` }">
                  <span class="dimension-value">{{ formatDimension(currentMouseDimensions.realHeight) }}</span>
                </div>
                <div class="dimension-indicator length-side" :style="{ width: `${currentMouseDimensions.length}px` }">
                  <span class="dimension-value">{{ formatDimension(currentMouseDimensions.realLength) }}</span>
                </div>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>

        <!-- Mouse specifications -->
        <div class="mouse-specifications">
          <h4 class="font-medium mb-2">鼠标规格</h4>
          <div class="grid grid-cols-1 xs:grid-cols-2 gap-4">
            <div class="spec-item">
              <span class="spec-label">长度:</span>
              <span class="spec-value">{{ formatDimension(currentMouseDimensions.realLength) }}</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">宽度:</span>
              <span class="spec-value">{{ formatDimension(currentMouseDimensions.realWidth) }}</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">高度:</span>
              <span class="spec-value">{{ formatDimension(currentMouseDimensions.realHeight) }}</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">握宽:</span>
              <span class="spec-value">{{ formatDimension(currentMouseDimensions.realGripWidth) }}</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">重量:</span>
              <span class="spec-value">{{ currentMouse.weight }}</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">推荐握法:</span>
              <span class="spec-value">{{ currentMouse.gripStyle }}</span>
            </div>
          </div>
        </div>
        
        <el-divider></el-divider>
        
        <!-- Enhanced Ruler Tool Integration -->
        <div>
          <h4 class="font-medium mb-2">与您的手部尺寸对比</h4>
          
          <div class="ruler-comparison mb-4">
            <div class="ruler-wrapper">
              <div ref="rulerRef" class="ruler relative w-full h-16 sm:h-20 border border-gray-300 rounded cursor-pointer">
                <div class="absolute inset-0" @click="handleClick" @touchstart="handleRulerTouch">
                  <canvas ref="canvasRef" class="w-full h-full"></canvas>
                </div>
                <div v-if="measurementPoint" class="absolute h-full border-l-2 border-red-500" :style="{ left: `${measurementPoint}px` }"></div>
                
                <!-- Reference mouse shape overlays -->
                <div v-if="isCalibrated && showMouseOverlay" class="mouse-overlay" 
                    :style="{ 
                      width: `${currentMouseDimensions.length * (calibrationFactor.value || 1)}px`,
                      left: `${measurementPoint || 0}px`
                    }">
                  <div class="overlay-indicator">{{ currentMouse.name }}</div>
                </div>
                
                <!-- Comparison mouse overlay -->
                <div v-if="isCalibrated && showMouseOverlay && comparisonMouseEnabled" 
                    class="mouse-overlay comparison-overlay" 
                    :style="{ 
                      width: `${getComparisonMouseDimensions().length * (calibrationFactor.value || 1)}px`,
                      left: `${measurementPoint || 0}px`
                    }">
                  <div class="overlay-indicator comparison-indicator">{{ getComparisonMouseName() }}</div>
                </div>
              </div>
              
              <div class="ruler-controls mt-2 flex flex-wrap gap-2 justify-between">
                <div class="flex items-center gap-2">
                  <el-checkbox v-model="showMouseOverlay" size="small">显示鼠标比例</el-checkbox>
                  <el-checkbox v-if="showMouseOverlay" v-model="comparisonMouseEnabled" size="small">对比其他鼠标</el-checkbox>
                </div>
                
                <div class="flex items-center gap-2">
                  <el-select 
                    v-if="showMouseOverlay && comparisonMouseEnabled" 
                    v-model="comparisonMouseId" 
                    size="small" 
                    placeholder="选择对比鼠标"
                    style="width: 140px;"
                  >
                    <el-option 
                      v-for="mouse in availableMiceForComparison" 
                      :key="mouse.id" 
                      :label="mouse.name" 
                      :value="mouse.id"
                    ></el-option>
                  </el-select>
                  
                  <el-button @click="calibrate" size="small" type="primary">
                    {{ isCalibrated ? '重新校准' : '校准尺子' }}
                  </el-button>
                </div>
              </div>
            </div>
            
            <div class="flex flex-col xs:flex-row justify-between items-start xs:items-center my-4 space-y-2 xs:space-y-0">
              <div class="text-left">
                <div class="text-xs sm:text-sm text-gray-500">点击尺子标记位置，然后与鼠标尺寸对比</div>
                <div v-if="isCalibrated" class="text-xs sm:text-sm text-green-500">已校准 ({{ calibrationFactor.toFixed(4) }})</div>
                <div v-else class="text-xs sm:text-sm text-orange-500">未校准，请点击校准按钮</div>
              </div>
              <div class="text-right">
                <div class="text-lg sm:text-xl font-bold">{{ formattedMeasurement }}</div>
                <div class="text-xs sm:text-sm text-gray-500">当前单位: {{ unitLabel }}</div>
              </div>
            </div>
          </div>
          
          <div class="recommendation mt-4" v-if="hasHandMeasurements">
            <h4 class="font-medium mb-2">适合度分析</h4>
            <div class="compatibility-score">
              <div class="score-label">与您手型的匹配度:</div>
              <el-progress 
                :percentage="compatibilityScore" 
                :color="compatibilityColor"
                :format="compatibilityFormatter"
              ></el-progress>
            </div>
            <div class="recommendation-text mt-2">
              {{ compatibilityMessage }}
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch, nextTick, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
// 移除手型测量功能，使用内部实现替代
// import { useMeasurement } from '@/composables/useMeasurement'

// 定义鼠标模型
interface MouseModel {
  id: string
  name: string
  dimensions: {
    length: number    // 长度(mm)
    width: number     // 宽度(mm)
    height: number    // 高度(mm)
    gripWidth: number // 握宽(mm)
  }
  weight: string      // 重量(g)
  gripStyle: string   // 推荐握法
  svgScaleFactor: number // SVG比例因子
}

// 特性标注定义
interface Annotation {
  x: number  // X坐标 (百分比)
  y: number  // Y坐标 (百分比)
  text: string // 说明文本
}

// 可用的鼠标模型
const availableMice = reactive<MouseModel[]>([
  {
    id: 'vxef1pro',
    name: 'Viper V2 Pro',
    dimensions: {
      length: 126.8,
      width: 66.2,
      height: 37.8,
      gripWidth: 57.6
    },
    weight: '58g',
    gripStyle: '爪握/指尖握',
    svgScaleFactor: 1.0
  },
  {
    id: 'gpw2',
    name: 'GPW 2 Superlight',
    dimensions: {
      length: 125.0,
      width: 63.5,
      height: 40.0,
      gripWidth: 59.5
    },
    weight: '63g',
    gripStyle: '掌握/指尖握',
    svgScaleFactor: 1.0
  },
  {
    id: 'hskpro',
    name: 'Pulsar X2 Mini',
    dimensions: {
      length: 114.0,
      width: 59.0,
      height: 36.0,
      gripWidth: 53.0
    },
    weight: '52g',
    gripStyle: '爪握/指尖握',
    svgScaleFactor: 1.0
  }
])

// 各鼠标的特性标注点
const mouseAnnotations = reactive({
  vxef1pro: {
    top: [
      { x: 50, y: 30, text: '主按键 - 轻触式设计' },
      { x: 50, y: 48, text: '滚轮 - 轻量化设计' },
      { x: 75, y: 55, text: '侧面按键' },
      { x: 30, y: 70, text: '人体工学侧边曲线' }
    ],
    side: [
      { x: 30, y: 30, text: '低前端高度 - 爪握适用' },
      { x: 50, y: 15, text: '曲线设计 - 提供掌部支撑' },
      { x: 70, y: 50, text: '侧面按键位置优化' }
    ]
  },
  gpw2: {
    top: [
      { x: 50, y: 30, text: '主按键 - 轻量级开关' },
      { x: 50, y: 48, text: '无孔滚轮设计' },
      { x: 75, y: 55, text: '侧面按键 - 可编程' },
      { x: 30, y: 70, text: '对称设计 - 适合多种握法' }
    ],
    side: [
      { x: 30, y: 30, text: '中等前端高度 - 多种握法适用' },
      { x: 50, y: 15, text: '顶部曲线 - 适合掌握' },
      { x: 70, y: 50, text: '优化底座分布' }
    ]
  },
  hskpro: {
    top: [
      { x: 50, y: 35, text: '主按键 - 快速响应' },
      { x: 50, y: 50, text: '小巧滚轮' },
      { x: 75, y: 55, text: '侧键 - 微动设计' },
      { x: 25, y: 60, text: '短小设计 - 适合小手' }
    ],
    side: [
      { x: 30, y: 30, text: '低扁平设计 - 指尖握推荐' },
      { x: 50, y: 15, text: '前端低高背后高设计' },
      { x: 70, y: 50, text: '底部PTFE脚垫' }
    ]
  }
})

// 选中的鼠标和视图
const selectedMouse = ref('vxef1pro')
const activeView = ref('top')
const unit = ref('mm')

// 缩放和平移相关
const zoomLevel = ref(1)
const isPanning = ref(false)
const panStartX = ref(0)
const panStartY = ref(0)
const currentTranslateX = ref(0)
const currentTranslateY = ref(0)
const svgContainer = ref<HTMLDivElement>()
const svgContainerSide = ref<HTMLDivElement>()

// 加载状态
const svgLoaded = ref(false)

// 尺子相关属性
const rulerRef = ref<HTMLDivElement>()
const canvasRef = ref<HTMLCanvasElement>()
const measurement = ref(0)
const measurementPoint = ref<number | null>(null)
const calibrationFactor = ref(1)
const isCalibrated = ref(false)
const showMouseOverlay = ref(false)
const comparisonMouseEnabled = ref(false)
const comparisonMouseId = ref('')

// 校准对话框已在前面定义过，此处不需要重复定义

// 内部实现手型尺寸判断功能
const getHandSizeCategory = (palm: number, length: number, unit: string): string => {
  // 将输入单位转换为毫米进行计算
  const palmInMm = unit === 'mm' ? palm : (unit === 'cm' ? palm * 10 : palm * 25.4)
  const lengthInMm = unit === 'mm' ? length : (unit === 'cm' ? length * 10 : length * 25.4)
  
  // 基于手掌宽度和长度进行简单分类
  if (palmInMm >= 95) return 'large'
  if (palmInMm >= 85) return 'medium'
  return 'small'
}

const getHandSizeName = (category: string): string => {
  switch (category) {
    case 'large': return '大型手'
    case 'medium': return '中型手'
    case 'small': return '小型手'
    default: return '未知'
  }
}

const convertUnit = (value: number, fromUnit: string, toUnit: string): number => {
  // 先转换到毫米
  let valueInMm: number
  if (fromUnit === 'mm') valueInMm = value
  else if (fromUnit === 'cm') valueInMm = value * 10
  else valueInMm = value * 25.4 // inch
  
  // 从毫米转换到目标单位
  if (toUnit === 'mm') return valueInMm
  if (toUnit === 'cm') return valueInMm / 10
  return valueInMm / 25.4 // inch
}
const hasHandMeasurements = ref(false)

// 模拟用户的手部测量数据（在实际应用中，这会从用户的保存数据中获取）
const handMeasurements = ref({
  palm: 90, // 默认值 mm
  length: 180, // 默认值 mm
})

// 计算当前选中的鼠标
const currentMouse = computed(() => {
  return availableMice.find(m => m.id === selectedMouse.value) || availableMice[0]
})

// 当前鼠标的SVG路径
const currentMouseSvg = computed(() => {
  return {
    top: `/${selectedMouse.value === 'vxef1pro' ? 'VXEF1PRO' : selectedMouse.value}-top.svg`,
    side: `/${selectedMouse.value === 'vxef1pro' ? 'VXEF1PRO' : selectedMouse.value}-side.svg`
  }
})

// 当前鼠标的标注
const currentAnnotations = computed(() => {
  return mouseAnnotations[selectedMouse.value as keyof typeof mouseAnnotations] || mouseAnnotations.vxef1pro
})

// 计算鼠标的尺寸（渲染大小）
const currentMouseDimensions = computed(() => {
  const mouseModel = currentMouse.value
  const pixelScaleFactor = 2.5 // 转换真实尺寸到像素的系数
  
  // 将实际尺寸从mm转换为显示单位
  const realLength = convertDimension(mouseModel.dimensions.length)
  const realWidth = convertDimension(mouseModel.dimensions.width)
  const realHeight = convertDimension(mouseModel.dimensions.height)
  const realGripWidth = convertDimension(mouseModel.dimensions.gripWidth)
  
  return {
    // 显示尺寸（像素）
    length: mouseModel.dimensions.length * pixelScaleFactor / mouseModel.svgScaleFactor,
    width: mouseModel.dimensions.width * pixelScaleFactor / mouseModel.svgScaleFactor,
    height: mouseModel.dimensions.height * pixelScaleFactor / mouseModel.svgScaleFactor,
    gripWidth: mouseModel.dimensions.gripWidth * pixelScaleFactor / mouseModel.svgScaleFactor,
    
    // 实际尺寸（用户选择的单位）
    realLength,
    realWidth,
    realHeight,
    realGripWidth
  }
})

// 根据当前单位转换尺寸
function convertDimension(value: number): number {
  if (unit.value === 'cm') {
    return value / 10
  } else if (unit.value === 'inch') {
    return value / 25.4
  }
  return value // mm
}

// 格式化尺寸显示
function formatDimension(value: number): string {
  return `${value.toFixed(1)} ${unit.value}`
}

// 计算尺子测量值的格式化字符串
const formattedMeasurement = computed(() => {
  if (measurement.value === 0) return '0 ' + unit.value
  return measurement.value.toFixed(1) + ' ' + unit.value
})

// 计算单位标签
const unitLabel = computed(() => {
  switch (unit.value) {
    case 'cm': return '厘米 (cm)'
    case 'mm': return '毫米 (mm)'
    case 'inch': return '英寸 (inch)'
    default: return '毫米 (mm)'
  }
})

// 计算鼠标与用户手部尺寸的兼容性分数
const compatibilityScore = computed(() => {
  if (!hasHandMeasurements.value) return 0
  
  const mouse = currentMouse.value
  const userHandSize = getHandSizeCategory(handMeasurements.value.palm, handMeasurements.value.length, 'mm')
  
  // 基于手掌宽度和鼠标宽度的匹配度
  let widthCompatibility = 0
  const palmWidth = handMeasurements.value.palm // mm
  const mouseWidth = mouse.dimensions.gripWidth // mm
  
  // 基于手掌宽度和鼠标握宽的理想比例判断 (假设理想比例为 1.55-1.65)
  const ratio = palmWidth / mouseWidth
  if (ratio >= 1.55 && ratio <= 1.65) {
    widthCompatibility = 100 // 完美匹配
  } else if (ratio >= 1.5 && ratio <= 1.7) {
    widthCompatibility = 90 // 非常好
  } else if (ratio >= 1.45 && ratio <= 1.75) {
    widthCompatibility = 75 // 良好
  } else if (ratio >= 1.4 && ratio <= 1.8) {
    widthCompatibility = 60 // 可接受
  } else {
    widthCompatibility = 40 // 不理想
  }
  
  // 基于鼠标长度和手指长度
  let lengthCompatibility = 0
  const fingerLength = handMeasurements.value.length // mm
  const mouseLength = mouse.dimensions.length // mm
  
  // 握法建议与长度比较
  if (mouse.gripStyle.includes('爪握') && fingerLength >= mouseLength * 0.8) {
    lengthCompatibility = 90
  } else if (mouse.gripStyle.includes('掌握') && fingerLength >= mouseLength * 0.9) {
    lengthCompatibility = 90
  } else if (mouse.gripStyle.includes('指尖握') && fingerLength >= mouseLength * 0.7) {
    lengthCompatibility = 90
  } else {
    lengthCompatibility = 70
  }
  
  // 总体兼容性评分
  return Math.round((widthCompatibility * 0.6 + lengthCompatibility * 0.4))
})

// 兼容性评分颜色
const compatibilityColor = computed(() => {
  const score = compatibilityScore.value
  if (score >= 90) return '#67C23A' // 绿色 - 优秀
  if (score >= 70) return '#E6A23C' // 橙色 - 良好
  return '#F56C6C' // 红色 - 不佳
})

// 格式化兼容性评分显示
const compatibilityFormatter = (percentage: number) => {
  if (percentage >= 90) return '优秀'
  if (percentage >= 70) return '良好'
  if (percentage >= 50) return '一般'
  return '不佳'
}

// 兼容性建议信息
const compatibilityMessage = computed(() => {
  const score = compatibilityScore.value
  const mouse = currentMouse.value.name
  
  if (score >= 90) {
    return `${mouse} 与您的手型非常匹配，这款鼠标应该会给您带来舒适的使用体验。`
  } else if (score >= 70) {
    return `${mouse} 与您的手型匹配度良好，可能需要稍微适应，但总体来说适合您。`
  } else if (score >= 50) {
    return `${mouse} 与您的手型匹配度一般，长时间使用可能会感到些许不适。`
  } else {
    return `${mouse} 可能不太适合您的手型，建议考虑其他更匹配的鼠标型号。`
  }
})

// SVG加载完成处理
const onSvgLoad = () => {
  svgLoaded.value = true
}

// 缩放功能
const zoomIn = () => {
  if (zoomLevel.value < 3) {
    zoomLevel.value += 0.2
  }
}

const zoomOut = () => {
  if (zoomLevel.value > 0.5) {
    zoomLevel.value -= 0.2
  }
}

const resetZoom = () => {
  zoomLevel.value = 1
  currentTranslateX.value = 0
  currentTranslateY.value = 0
  updateTransform()
}

// 平移功能
const startPan = (e: MouseEvent) => {
  isPanning.value = true
  panStartX.value = e.clientX
  panStartY.value = e.clientY
}

const pan = (e: MouseEvent) => {
  if (!isPanning.value) return
  
  const dx = e.clientX - panStartX.value
  const dy = e.clientY - panStartY.value
  
  // 更新平移坐标
  currentTranslateX.value += dx
  currentTranslateY.value += dy
  
  // 更新起始位置
  panStartX.value = e.clientX
  panStartY.value = e.clientY
  
  updateTransform()
}

const endPan = () => {
  isPanning.value = false
}

// 触摸屏支持
const handleTouchStart = (e: TouchEvent) => {
  if (e.touches.length === 1) {
    isPanning.value = true
    panStartX.value = e.touches[0].clientX
    panStartY.value = e.touches[0].clientY
    e.preventDefault()
  }
}

const handleTouchMove = (e: TouchEvent) => {
  if (!isPanning.value || e.touches.length !== 1) return
  
  const dx = e.touches[0].clientX - panStartX.value
  const dy = e.touches[0].clientY - panStartY.value
  
  // 更新平移坐标
  currentTranslateX.value += dx
  currentTranslateY.value += dy
  
  // 更新起始位置
  panStartX.value = e.touches[0].clientX
  panStartY.value = e.touches[0].clientY
  
  updateTransform()
  e.preventDefault()
}

const handleTouchEnd = () => {
  isPanning.value = false
}

// 尺子触摸支持
const handleRulerTouch = (e: TouchEvent) => {
  if (e.touches.length === 1) {
    const canvas = canvasRef.value
    if (!canvas) return
    
    const rect = canvas.getBoundingClientRect()
    const x = e.touches[0].clientX - rect.left
    
    measurementPoint.value = x
    
    // 计算测量值
    if (isCalibrated.value && calibrationFactor.value > 0) {
      const rawValue = x / calibrationFactor.value
      measurement.value = rawValue
    }
    
    e.preventDefault()
  }
}

const updateTransform = () => {
  const container = activeView.value === 'top' ? svgContainer.value : svgContainerSide.value
  if (!container) return
  
  const wrapper = container.querySelector('.svg-wrapper') as HTMLDivElement
  if (!wrapper) return
  
  wrapper.style.transform = `scale(${zoomLevel.value}) translate(${currentTranslateX.value}px, ${currentTranslateY.value}px)`
}

// 窗口大小变化处理
const handleResize = () => {
  resizeCanvas()
}

// 初始化
onMounted(() => {
  initRuler()
  window.addEventListener('resize', handleResize)
  resizeCanvas()
  
  // 模拟获取用户的手部测量数据
  setTimeout(() => {
    hasHandMeasurements.value = true
  }, 500)
  
  // 预加载SVG
  preloadSvgs()
})

// 清理
onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
})

// 预加载SVG
const preloadSvgs = () => {
  const svgPaths = [
    '/VXEF1PRO-top.svg',
    '/VXEF1PRO-side.svg',
    '/gpw2-top.svg',
    '/gpw2-side.svg',
    '/hskpro-top.svg',
    '/hskpro-side.svg'
  ]
  
  // 使用 Image 对象预加载
  svgPaths.forEach(path => {
    const img = new Image()
    img.src = path
  })
}

// 当单位变化时，更新尺子和尺寸显示
watch(unit, () => {
  initRuler()
})

// 当选中的鼠标改变时，重置缩放和平移
watch(selectedMouse, () => {
  resetZoom()
  svgLoaded.value = false
  
  // 如果当前选中的鼠标与比较鼠标相同，重置比较鼠标
  if (selectedMouse.value === comparisonMouseId.value) {
    comparisonMouseId.value = ''
  }
  
  nextTick(() => {
    // 模拟延迟加载
    setTimeout(() => {
      svgLoaded.value = true
    }, 300)
  })
})

// 可比较的鼠标列表（排除当前选中的鼠标）
const availableMiceForComparison = computed(() => {
  return availableMice.filter(mouse => mouse.id !== selectedMouse.value)
})

// 获取比较鼠标的详情
const getComparisonMouse = () => {
  return availableMice.find(m => m.id === comparisonMouseId.value) || null
}

// 获取比较鼠标的名称
const getComparisonMouseName = () => {
  const mouse = getComparisonMouse()
  return mouse ? mouse.name : '未选择'
}

// 获取比较鼠标的尺寸
const getComparisonMouseDimensions = () => {
  const mouse = getComparisonMouse()
  if (!mouse) return { length: 0, width: 0, height: 0, gripWidth: 0 }
  
  const pixelScaleFactor = 2.5 // 转换真实尺寸到像素的系数
  
  return {
    length: mouse.dimensions.length * pixelScaleFactor / mouse.svgScaleFactor,
    width: mouse.dimensions.width * pixelScaleFactor / mouse.svgScaleFactor,
    height: mouse.dimensions.height * pixelScaleFactor / mouse.svgScaleFactor,
    gripWidth: mouse.dimensions.gripWidth * pixelScaleFactor / mouse.svgScaleFactor
  }
}

// 视图改变时重置缩放和平移
watch(activeView, () => {
  resetZoom()
})

// 调整画布大小
const resizeCanvas = () => {
  const canvas = canvasRef.value
  const container = rulerRef.value
  if (!canvas || !container) return
  
  canvas.width = container.clientWidth
  canvas.height = container.clientHeight
  initRuler()
}

// 初始化尺子
const initRuler = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return
  
  const width = canvas.width
  const height = canvas.height
  
  // 清除画布
  ctx.clearRect(0, 0, width, height)
  
  // 绘制背景
  ctx.fillStyle = '#f8f9fa'
  ctx.fillRect(0, 0, width, height)
  
  // 绘制主刻度线
  const majorTickInterval = unit.value === 'mm' ? 10 : 1
  const minorTicksPerMajor = unit.value === 'mm' ? 10 : 10
  const pixelsPerMajorTick = 50 * calibrationFactor.value // 校准后的每刻度像素数
  
  // 计算可以显示的刻度数量
  const numMajorTicks = Math.floor(width / pixelsPerMajorTick) + 1
  
  ctx.strokeStyle = '#333'
  ctx.lineWidth = 1
  ctx.textAlign = 'center'
  ctx.font = '10px Arial'
  
  for (let i = 0; i < numMajorTicks; i++) {
    const x = i * pixelsPerMajorTick
    
    // 绘制主刻度线
    ctx.beginPath()
    ctx.moveTo(x, height * 0.1)
    ctx.lineTo(x, height * 0.4)
    ctx.stroke()
    
    // 绘制刻度值
    ctx.fillStyle = '#333'
    ctx.fillText(String(i * majorTickInterval), x, height * 0.7)
    
    // 绘制次要刻度线
    if (i < numMajorTicks - 1) {
      const minorTickSpacing = pixelsPerMajorTick / minorTicksPerMajor
      for (let j = 1; j < minorTicksPerMajor; j++) {
        const minorX = x + j * minorTickSpacing
        const tickHeight = j % 5 === 0 ? height * 0.3 : height * 0.2 // 每5个次要刻度稍长
        
        ctx.beginPath()
        ctx.moveTo(minorX, height * 0.1)
        ctx.lineTo(minorX, tickHeight)
        ctx.stroke()
      }
    }
  }
}

// 校准逻辑
const calibrateDialogVisible = ref(false)
const calibrationValue = ref('')
const calibrationStep = ref(1) // 1: 标记位置, 2: 输入实际长度

const calibrate = () => {
  // 重置校准状态
  calibrationStep.value = 1
  calibrationValue.value = ''
  
  // 显示校准对话框
  calibrateDialogVisible.value = true
  
  // 首先提示用户标记位置
  ElMessage.info('请在尺子上标记一个已知长度')
  
  // 第一步：等待用户点击标记位置
  isCalibrated.value = false
}

// 用户输入实际长度后的处理
const handleCalibrateConfirm = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  
  if (calibrationStep.value === 1) {
    // 已标记位置，转到第二步
    if (measurementPoint.value) {
      calibrationStep.value = 2
    } else {
      ElMessage.warning('请先在尺子上标记位置')
      return
    }
  } else if (calibrationStep.value === 2) {
    // 处理用户输入的实际长度
    if (calibrationValue.value && !isNaN(Number(calibrationValue.value)) && measurementPoint.value) {
      const pixelLength = measurementPoint.value
      const actualLength = parseFloat(calibrationValue.value)
      
      if (actualLength > 0) {
        // 计算校准因子
        calibrationFactor.value = pixelLength / actualLength
        isCalibrated.value = true
        ElMessage.success('校准成功!')
        
        // 重绘尺子
        initRuler()
        
        // 关闭对话框
        calibrateDialogVisible.value = false
      } else {
        ElMessage.error('请输入有效的长度值')
      }
    } else {
      ElMessage.error('请输入有效的长度值')
    }
  }
}

// 处理点击事件
const handleClick = (e: MouseEvent) => {
  const canvas = canvasRef.value
  if (!canvas) return
  
  const rect = canvas.getBoundingClientRect()
  const x = e.clientX - rect.left
  
  measurementPoint.value = x
  
  // 计算测量值
  if (isCalibrated.value && calibrationFactor.value > 0) {
    const rawValue = x / calibrationFactor.value
    
    // 根据不同单位进行转换
    let finalValue = rawValue
    
    // 设置测量值
    measurement.value = finalValue
  }
}
</script>

<style scoped>
.mouse-shape-visualization {
  max-width: 800px;
  margin: 0 auto;
}

.svg-container {
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  border: 1px solid #eee;
  border-radius: 4px;
  padding: 10px;
  margin: 10px 0;
  background: #fafafa;
  overflow: hidden;
}

.zoom-controls {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 20;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 4px;
  padding: 2px;
}

.svg-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: transform 0.1s ease;
}

.mouse-svg {
  max-width: 80%;
  max-height: 80%;
  filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.1));
}

.loading-placeholder {
  width: 80%;
  max-width: 400px;
  padding: 20px;
  background: #f5f5f5;
  border-radius: 8px;
}

/* Feature annotations */
.feature-annotations {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.annotation {
  position: absolute;
  transform: translate(-50%, -50%);
  z-index: 10;
}

.annotation-marker {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: #409eff;
  border: 2px solid white;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.2);
  position: relative;
  z-index: 2;
}

.annotation-text {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
  background-color: rgba(255, 255, 255, 0.95);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  color: #333;
  border: 1px solid #ddd;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 5px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.annotation:hover .annotation-text {
  opacity: 1;
}

/* Enhanced dimension indicators */
.dimension-indicator {
  position: absolute;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #409eff;
  font-size: 12px;
  z-index: 5;
}

.dimension-indicator::before,
.dimension-indicator::after {
  content: '';
  position: absolute;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: #409eff;
}

.dimension-indicator.length {
  border-top: 2px solid #409eff;
  border-left: 2px solid #409eff;
  border-right: 2px solid #409eff;
  height: 24px;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
}

.dimension-indicator.length::before {
  left: -3px;
  bottom: -3px;
}

.dimension-indicator.length::after {
  right: -3px;
  bottom: -3px;
}

.dimension-indicator.width {
  border-top: 2px solid #409eff;
  border-left: 2px solid #409eff;
  border-right: 2px solid #409eff;
  height: 24px;
  top: 30px;
  left: 50%;
  transform: translateX(-50%);
}

.dimension-indicator.width::before {
  left: -3px;
  bottom: -3px;
}

.dimension-indicator.width::after {
  right: -3px;
  bottom: -3px;
}

.dimension-indicator.grip-width {
  border-top: 2px solid #409eff;
  border-left: 2px solid #409eff;
  border-right: 2px solid #409eff;
  height: 24px;
  top: 100px;
  left: 50%;
  transform: translateX(-50%);
}

.dimension-indicator.grip-width::before {
  left: -3px;
  bottom: -3px;
}

.dimension-indicator.grip-width::after {
  right: -3px;
  bottom: -3px;
}

.dimension-indicator.height {
  border-left: 2px solid #409eff;
  border-top: 2px solid #409eff;
  border-bottom: 2px solid #409eff;
  width: 24px;
  right: 30px;
  top: 50%;
  transform: translateY(-50%);
}

.dimension-indicator.height::before {
  left: -3px;
  top: -3px;
}

.dimension-indicator.height::after {
  left: -3px;
  bottom: -3px;
}

.dimension-indicator.length-side {
  border-top: 2px solid #409eff;
  border-left: 2px solid #409eff;
  border-right: 2px solid #409eff;
  height: 24px;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
}

.dimension-indicator.length-side::before {
  left: -3px;
  bottom: -3px;
}

.dimension-indicator.length-side::after {
  right: -3px;
  bottom: -3px;
}

.dimension-value {
  background-color: white;
  padding: 3px 6px;
  border: 1px solid #ddd;
  border-radius: 3px;
  font-weight: 500;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* Enhanced ruler integration */
.ruler-wrapper {
  position: relative;
}

.ruler {
  transition: all 0.3s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.ruler:hover {
  border-color: #409eff;
  box-shadow: 0 0 5px rgba(64, 158, 255, 0.3);
}

.mouse-overlay {
  position: absolute;
  height: 10px;
  bottom: 0;
  background-color: rgba(64, 158, 255, 0.2);
  border: 1px dashed #409eff;
  pointer-events: none;
}

.comparison-overlay {
  background-color: rgba(245, 108, 108, 0.2);
  border: 1px dashed #F56C6C;
  bottom: 12px;
}

.overlay-indicator {
  position: absolute;
  top: -15px;
  left: 0;
  right: 0;
  text-align: center;
  font-size: 10px;
  color: #409eff;
  background-color: white;
  padding: 1px 4px;
  border-radius: 2px;
  white-space: nowrap;
}

.comparison-indicator {
  top: -15px;
  color: #F56C6C;
}

/* Mouse specifications */
.spec-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.spec-label {
  color: #666;
}

.spec-value {
  font-weight: 500;
  color: #333;
}

/* Compatibility score */
.compatibility-score {
  margin-top: 12px;
}

.score-label {
  margin-bottom: 5px;
  color: #666;
  font-size: 14px;
}

.recommendation-text {
  background-color: #f8f9fa;
  padding: 12px;
  border-radius: 4px;
  border-left: 3px solid #409eff;
  font-size: 14px;
  line-height: 1.5;
}
</style>