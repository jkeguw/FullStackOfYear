<template>
  <div class="dpi-calculator">
    <el-form :model="form" label-width="120px">
      <el-form-item label="originalDPI">
        <el-input-number v-model="form.sourceDPI" :min="100" :max="32000" />
      </el-form-item>
      <el-form-item label="targetDPI">
        <el-input-number v-model="form.targetDPI" :min="100" :max="32000" />
      </el-form-item>
      <el-form-item label="originalSensitivity">
        <el-input-number v-model="form.sensitivity" :min="0.1" :max="10" :step="0.1" />
      </el-form-item>
      <el-form-item>
        <div>Sensitivity after conversion: {{ convertedSensitivity }}</div>
      </el-form-item>
    </el-form>
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
</script>