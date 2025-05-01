<template>
  <div class="test-mouse-comparison">
    <h1 class="text-xl font-bold mb-4">Mouse Comparison Test Page</h1>
    <p class="mb-4">Mouse comparison test using hardcoded SVG data</p>
    
    <div class="mb-6">
      <el-alert
        type="info"
        :closable="false"
        show-icon
      >
        <div>Select the following mice for comparison:</div>
        <div class="mt-2">
          <el-checkbox-group v-model="selectedMiceIds">
            <el-checkbox v-for="mouse in mockMice" :key="mouse.id" :label="mouse.id">
              {{ mouse.brand }} {{ mouse.name }}
            </el-checkbox>
          </el-checkbox-group>
        </div>
      </el-alert>
    </div>
    
    <div v-if="selectedMiceIds.length > 0" class="mb-6">
      <el-button type="primary" @click="addSelectedMice">
        Add selected mice to comparison view
      </el-button>
      <el-button @click="clearComparison">
        Clear comparison view
      </el-button>
    </div>
    
    <!-- Mouse comparison view component -->
    <MouseComparisonView />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useComparisonStore } from '@/stores';
import MouseComparisonView from './MouseComparisonView.vue';
import { mockMice } from '@/data/mockMice';

const comparisonStore = useComparisonStore();
const selectedMiceIds = ref<string[]>([]);

// Add selected mice to comparison view
function addSelectedMice() {
  // First clear current selection
  clearComparison();
  
  // Add newly selected mice
  for (const id of selectedMiceIds.value) {
    const mouse = mockMice.find(m => m.id === id);
    if (mouse) {
      comparisonStore.addMouse(mouse);
    }
  }
}

// Clear comparison view
function clearComparison() {
  comparisonStore.clearMice();
}

// Initially select the first mouse by default
onMounted(() => {
  if (mockMice.length > 0) {
    selectedMiceIds.value = [mockMice[0].id];
    addSelectedMice();
  }
});
</script>

<style scoped>
.test-mouse-comparison {
  padding: 1rem;
  max-width: 1200px;
  margin: 0 auto;
}
</style>