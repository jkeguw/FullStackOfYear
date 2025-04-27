<template>
  <div class="ruler-tool-page">
    <h2 class="text-2xl font-bold mb-6">在线厘米尺</h2>
    
    <div class="instructions bg-gray-700 rounded-lg p-4 mb-6">
      <h3 class="text-xl font-bold mb-2">使用说明</h3>
      <ul class="list-disc pl-5 space-y-2">
        <li>提供精确的厘米和毫米标记</li>
        <li>水平和垂直两种尺子方向</li>
        <li>所有刻度都标准显示，厘米显示数字，毫米使用短线标记</li>
      </ul>
    </div>
    
    <div class="ruler-workspace relative bg-gray-800 rounded-lg border border-gray-600 p-4 min-h-[400px] mb-6">
      <div class="text-center text-gray-400 mb-4">水平测量尺</div>
      <StaticRuler mode="horizontal" :rulerLength="30" />
    </div>
    
    <div class="ruler-workspace relative bg-gray-800 rounded-lg border border-gray-600 p-4 min-h-[400px]">
      <div class="text-center text-gray-400 mb-4">垂直测量尺</div>
      <StaticRuler mode="vertical" :rulerLength="20" />
    </div>
  </div>
</template>

<script setup lang="ts">
import StaticRuler from '@/components/tools/StaticRuler.vue';
</script>

<style scoped>
.ruler-workspace {
  position: relative;
  overflow: hidden;
}
</style>