<template>
  <div class="dpi-calculator p-4 sm:p-6 bg-gray-800 rounded-lg shadow-sm text-white">
    <h3 class="text-lg font-medium mb-4">DPI转换计算器</h3>
    <el-form :model="form" label-position="top" class="max-w-md mx-auto">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <el-form-item label="原始DPI" class="mb-3">
          <el-input-number 
            v-model="form.sourceDPI" 
            :min="100" 
            :max="32000" 
            class="w-full" 
            controls-position="right" 
          />
        </el-form-item>
        <el-form-item label="目标DPI" class="mb-3">
          <el-input-number 
            v-model="form.targetDPI" 
            :min="100" 
            :max="32000" 
            class="w-full" 
            controls-position="right" 
          />
        </el-form-item>
      </div>
      <el-form-item label="原始灵敏度" class="mb-3">
        <el-input-number 
          v-model="form.sensitivity" 
          :min="0.1" 
          :max="10" 
          :step="0.1" 
          class="w-full" 
          controls-position="right" 
        />
      </el-form-item>
      <el-form-item class="mt-6 text-center">
        <div class="p-3 bg-blue-900 rounded-md">
          <span class="text-white">转换后灵敏度: </span>
          <span class="text-xl font-bold text-white">{{ convertedSensitivity }}</span>
        </div>
      </el-form-item>
      <div class="mt-4 p-3 bg-gray-700 rounded-md">
        <p class="text-sm text-gray-300">cm/360° = <span class="font-semibold text-white">{{ cmPer360 }}</span></p>
      </div>
    </el-form>

    <div class="bg-gray-700 rounded-lg p-4 mt-8">
      <h4 class="text-lg font-medium mb-3">DPI计算公式</h4>
      <div class="text-gray-300 text-sm space-y-2">
        <p>此计算器使用以下公式进行DPI和灵敏度转换：</p>
        <div class="p-2 bg-gray-800 rounded">
          <p>新灵敏度 = (原灵敏度 × 原DPI) ÷ 新DPI</p>
        </div>
        <p>此公式确保在DPI变化时保持相同的实际鼠标移动距离，从而保持一致的游戏体验。</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { convertDPI } from '@/utils/dpi'

const form = ref({
  sourceDPI: 800,
  targetDPI: 800,
  sensitivity: 1
})

const convertedSensitivity = computed(() => {
  return convertDPI(
    form.value.sourceDPI,
    form.value.targetDPI,
    form.value.sensitivity
  ).toFixed(3)
})

const cmPer360 = computed(() => {
  // cm/360° = 360 × 2.54 / (灵敏度 × DPI × 游戏系数)
  // 使用CSGO的游戏系数0.022作为默认值
  const gameCoefficient = 0.022
  const cm = (360 * 2.54) / (form.value.sensitivity * form.value.sourceDPI * gameCoefficient)
  return cm.toFixed(2) + ' cm'
})
</script>

<style scoped>
.el-input-number :deep(.el-input__inner) {
  color: white !important;
  background-color: #4a4a4a !important;
}

.el-form-item :deep(.el-form-item__label) {
  color: #e0e0e0 !important;
}
</style>