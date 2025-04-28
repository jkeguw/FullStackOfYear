<template>
  <div class="ruler-tool-page">
    <h2 class="text-2xl font-bold mb-6">在线厘米尺</h2>

    <div class="instructions bg-gray-700 rounded-lg p-4 mb-6">
      <h3 class="text-xl font-bold mb-2 text-red-500">功能已废弃</h3>
      <p class="text-gray-300 mb-4">此工具已被废弃并移除，请使用其他测量工具代替。</p>
      <el-alert
        title="此功能已不再维护"
        type="warning"
        description="厘米尺功能已被完全移除，请使用其他测量工具。"
        show-icon
      />
    </div>

    <div class="flex justify-center mt-8">
      <router-link to="/tools" class="el-button el-button--primary">
        返回工具集
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
// 尺子组件已被移除
</script>

<style scoped>
.ruler-tool-page {
  min-height: 400px;
}
</style>
