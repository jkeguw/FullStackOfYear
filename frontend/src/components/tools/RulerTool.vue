<template>
  <div class="ruler-tool">
    <el-card>
      <template #header>
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-medium">手部测量工具</h3>
          <div>
            <el-select v-model="unit" size="small" class="mr-2">
              <el-option label="厘米 (cm)" value="cm"></el-option>
              <el-option label="毫米 (mm)" value="mm"></el-option>
              <el-option label="英寸 (inch)" value="inch"></el-option>
            </el-select>
            <el-button @click="calibrate" size="small" type="primary">{{ isCalibrated ? '重新校准' : '校准' }}</el-button>
          </div>
        </div>
      </template>
      
      <div class="flex flex-col space-y-6">
        <div ref="rulerRef" class="ruler relative w-full h-16 border border-gray-300 rounded cursor-pointer">
          <div class="absolute inset-0" @click="handleClick">
            <canvas ref="canvasRef" class="w-full h-full"></canvas>
          </div>
          <div v-if="measurementPoint" class="absolute h-full border-l-2 border-red-500" :style="{ left: `${measurementPoint}px` }"></div>
        </div>
        
        <div class="flex justify-between items-center">
          <div class="text-left">
            <div class="text-sm text-gray-500">测量原理：点击尺子标记位置，然后输入实际手部尺寸</div>
            <div v-if="isCalibrated" class="text-sm text-green-500">已校准 ({{ calibrationFactor.toFixed(4) }})</div>
            <div v-else class="text-sm text-orange-500">未校准，请校准后使用</div>
          </div>
          <div class="text-right">
            <div class="text-xl font-bold">{{ formattedMeasurement }}</div>
            <div class="text-sm text-gray-500">当前单位: {{ unitLabel }}</div>
          </div>
        </div>
        
        <el-divider></el-divider>
        
        <div>
          <div class="mb-2 font-medium">手部测量数据</div>
          <el-form :model="form" label-width="100px">
            <el-form-item label="手掌宽度">
              <el-input-number v-model="form.palm" :min="50" :max="150" :precision="1" :step="0.5" controls-position="right">
                <template #suffix>{{ unit }}</template>
              </el-input-number>
            </el-form-item>
            <el-form-item label="手指长度">
              <el-input-number v-model="form.length" :min="40" :max="120" :precision="1" :step="0.5" controls-position="right">
                <template #suffix>{{ unit }}</template>
              </el-input-number>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveMeasurement" :disabled="!isCalibrated">保存测量</el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useMeasurement } from '@/composables/useMeasurement'

// 尺子相关
const rulerRef = ref<HTMLDivElement>()
const canvasRef = ref<HTMLCanvasElement>()
const measurement = ref(0)
const measurementPoint = ref<number | null>(null)
const calibrationFactor = ref(1) // 校准因子，用于将像素转换为实际单位
const isCalibrated = ref(false)
const unit = ref('cm') // 当前单位：cm, mm, inch

// 表单数据
const form = ref({
  palm: 80, // 默认值
  length: 70, // 默认值
  calibrated: false
})

// 计算属性
const formattedMeasurement = computed(() => {
  if (measurement.value === 0) return '0 ' + unit.value
  return measurement.value.toFixed(1) + ' ' + unit.value
})

const unitLabel = computed(() => {
  switch (unit.value) {
    case 'cm': return '厘米 (cm)'
    case 'mm': return '毫米 (mm)'
    case 'inch': return '英寸 (inch)'
    default: return '厘米 (cm)'
  }
})

// 初始化
onMounted(() => {
  initRuler()
  // 设置画布尺寸以匹配其容器
  window.addEventListener('resize', resizeCanvas)
  resizeCanvas()
})

// 当单位变化时重绘尺子
watch(unit, () => {
  initRuler()
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
const calibrate = () => {
  ElMessage.info('请在尺子上标记一个已知长度，然后输入该长度')
  
  const canvas = canvasRef.value
  if (!canvas) return
  
  const inputLength = ref<number | null>(null)
  
  // 标记步骤
  ElMessage.info('请单击尺子上的任意点进行标记')
  
  // 接下来会通过点击处理用户标记，然后获取实际长度
  isCalibrated.value = false
  
  // 实际上这部分可以通过一个对话框来完成，但为简化，这里使用简单提示
  setTimeout(() => {
    const knownLength = prompt('请输入标记点对应的实际长度 (单位: ' + unit.value + ')')
    
    if (knownLength && !isNaN(Number(knownLength)) && measurementPoint.value) {
      const pixelLength = measurementPoint.value
      const actualLength = parseFloat(knownLength)
      
      if (actualLength > 0) {
        // 计算校准因子
        calibrationFactor.value = pixelLength / actualLength
        isCalibrated.value = true
        ElMessage.success('校准成功!')
        
        // 重绘尺子
        initRuler()
      } else {
        ElMessage.error('请输入有效的长度值')
      }
    } else {
      ElMessage.error('校准失败，请重试')
    }
  }, 1000)
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

// 保存测量数据
const saveMeasurement = async () => {
  if (!isCalibrated.value) {
    ElMessage.warning('请先校准尺子后再保存测量结果')
    return
  }
  
  try {
    // 调用API将数据发送到后端
    const data = {
      palm: form.value.palm,
      length: form.value.length,
      unit: unit.value,
      calibrated: isCalibrated.value,
      device: navigator.userAgent // 简单记录设备信息
    }
    
    // 这里使用 useMeasurement composable 保存数据
    const { saveMeasurement: saveData } = useMeasurement()
    await saveData(data)
    
    ElMessage.success('测量数据已保存')
    
    // 重置测量点
    measurementPoint.value = null
  } catch (error) {
    console.error('保存数据失败:', error)
    ElMessage.error('保存数据失败，请重试')
  }
}

// 重置表单
const resetForm = () => {
  form.value = {
    palm: 80,
    length: 70,
    calibrated: isCalibrated.value
  }
  measurementPoint.value = null
}
</script>

<style scoped>
.ruler {
  transition: all 0.3s ease;
}

.ruler:hover {
  border-color: #409eff;
  box-shadow: 0 0 5px rgba(64, 158, 255, 0.3);
}
</style>