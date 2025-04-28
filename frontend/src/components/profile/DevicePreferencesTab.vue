<template>
  <div class="device-preferences-tab">
    <h2 class="text-2xl font-semibold mb-4">功能已停用</h2>
    <p class="text-lg text-gray-600 mb-6">个人设备配置功能已被移除。请使用设备比较功能来比较不同设备。</p>
    
    <el-alert
      title="该功能已被移除"
      type="warning"
      description="根据项目计划，个人设备配置管理功能已被移除。如有任何疑问，请联系管理员。"
      show-icon
      :closable="false"
      class="mb-6"
    />
    
    <el-button type="primary" @click="navigateToComparison">
      前往设备比较页面
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();

// 导航到比较页面
const navigateToComparison = () => {
  router.push('/compare');
};
</script>

<style scoped>
.device-preferences-tab {
  padding: 20px 0;
}
</style>
