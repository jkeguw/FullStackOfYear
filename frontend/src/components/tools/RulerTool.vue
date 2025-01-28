<template>
  <div class="ruler-tool">
    <el-card>
      <div ref="rulerRef" class="ruler relative w-96 h-12 border">
        <div class="absolute inset-0" @click="handleClick">
          <canvas ref="canvasRef" class="w-full h-full"></canvas>
        </div>
        <div class="text-center mt-4">
          <span>{{ measurement }}cm</span>
          <el-button @click="calibrate" size="small" class="ml-2">校准</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const rulerRef = ref<HTMLDivElement>()
const canvasRef = ref<HTMLCanvasElement>()
const measurement = ref(0)

onMounted(() => {
  initRuler()
})

const initRuler = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return
  // 绘制刻度
}

const calibrate = () => {
  // 校准逻辑
}

const handleClick = (e: MouseEvent) => {
  // 测量逻辑
}
</script>