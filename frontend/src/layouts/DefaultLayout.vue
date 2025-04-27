<template>
  <div class="default-layout min-h-screen flex flex-col">
    <app-header />
    
    <main class="flex-grow">
      <router-view />
    </main>
    
    <app-footer />
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/layout/AppHeader.vue'
import AppFooter from '@/components/layout/AppFooter.vue'
</script>

<style scoped>
.default-layout {
  background-color: #f5f7fa;
}
</style>