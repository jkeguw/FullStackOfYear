<template>
  <div class="mouse-selector">
    <div class="grid grid-cols-1 gap-4">
      <!-- Search bar -->
      <div class="search-container mb-4">
        <el-input v-model="searchQuery" placeholder="Search mice..." :prefix-icon="Search" clearable />
      </div>

      <!-- Mouse list -->
      <div
        class="mouse-list bg-[#1A1A1A] border border-[#333333] rounded-lg p-4 max-h-96 overflow-y-auto"
      >
        <div v-if="loading" class="text-center py-4">
          <el-skeleton :rows="5" animated />
        </div>

        <div v-else-if="filteredMice.length === 0" class="text-center py-8 text-gray-400">
          No matching mice found
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="mouse in filteredMice"
            :key="mouse.id"
            class="mouse-item p-3 rounded-lg bg-[#242424] hover:bg-[#333333] cursor-pointer border border-[#333333] flex items-center"
            @click="toggleMouse(mouse)"
          >
            <div
              class="flex-shrink-0 w-16 h-16 bg-[#1A1A1A] rounded-lg flex items-center justify-center mr-3 overflow-hidden"
            >
              <img
                v-if="mouse.imageUrl"
                :src="mouse.imageUrl"
                :alt="mouse.name"
                class="max-w-full max-h-full object-contain"
              />
              <el-icon v-else class="text-3xl text-gray-500"><Mouse /></el-icon>
            </div>

            <div class="flex-grow">
              <h4 class="text-white font-medium">{{ mouse.name }}</h4>
              <div class="text-sm text-gray-400">{{ mouse.brand }}</div>
              <div class="text-xs text-gray-500 mt-1">
                {{ mouse.weight }}g | {{ mouse.connection_type }} | {{ mouse.sensor }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recently viewed mice -->
      <div
        v-if="recentlyViewedMice.length > 0"
        class="recently-viewed bg-[#1A1A1A] border border-[#333333] rounded-lg p-4"
      >
        <h3 class="text-lg font-medium text-white mb-4">Recently Viewed</h3>

        <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
          <div
            v-for="mouse in recentlyViewedMice"
            :key="mouse.id"
            class="recent-mouse-item p-2 rounded-lg bg-[#242424] hover:bg-[#333333] cursor-pointer border border-[#333333] flex flex-col items-center"
            @click="toggleMouse(mouse)"
          >
            <div
              class="w-12 h-12 bg-[#1A1A1A] rounded-lg flex items-center justify-center mb-2 overflow-hidden"
            >
              <img
                v-if="mouse.imageUrl"
                :src="mouse.imageUrl"
                :alt="mouse.name"
                class="max-w-full max-h-full object-contain"
              />
              <el-icon v-else class="text-2xl text-gray-500"><Mouse /></el-icon>
            </div>

            <div class="text-sm text-white text-center truncate w-full">{{ mouse.name }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, defineProps, defineEmits } from 'vue';
import { Search, Plus, Minus, Mouse } from '@element-plus/icons-vue';
import { useComparisonStore, type MouseDevice } from '@/stores';
import { getSVGMouseList } from '@/api/device';
import { hardcodedMice } from '@/data/hardcodedMice'; // Import hardcoded mouse data

// Define props
const props = defineProps({
  initialSelectedMice: {
    type: Array as () => MouseDevice[],
    default: () => []
  },
  maxSelection: {
    type: Number,
    default: 3
  }
});

// Define events
const emit = defineEmits(['select', 'cancel']);

// Use store
const comparisonStore = useComparisonStore();

// Search query
const searchQuery = ref('');
const loading = ref(true);

// Get mouse data from API
const mice = ref([]);

// Load data
onMounted(async () => {
  // Initialize selected mice
  if (props.initialSelectedMice && props.initialSelectedMice.length > 0) {
    // If there are initially selected mice, make sure they are marked as selected in the store
    props.initialSelectedMice.forEach((mouse) => {
      if (!comparisonStore.isMouseSelected(mouse.id)) {
        comparisonStore.addMouse(mouse);
      }
    });
  }

  try {
    loading.value = true;
    
    try {
      // Try to get data from API
      console.log('Requesting SVG mouse list with views: top,side');
      const response = await getSVGMouseList({ views: ['top', 'side'] });
      
      // Print response for debugging
      console.log('API response:', response);
      
      // Check if the response format is correct (even if it returns an empty list)
      if (response && response.code === 0 && response.data) {
        // Ensure devices exists, even if it's not an array or is an empty array
        const devices = response.data.devices || [];
        
        // Check if it's an array and log
        if (!Array.isArray(devices)) {
          console.warn('API returned devices is not an array:', devices);
        }
        
        // API returned valid format data
        if (Array.isArray(devices) && devices.length > 0) {
          // Has mouse data, use API data
          console.log('Mapping API data to mouse objects, device count:', devices.length);
          
          mice.value = devices.map(device => {
            // 处理SVG数据供预览使用（如果有）
            const svgData = {
              top: device.svgData?.topView || '',
              side: device.svgData?.sideView || ''
            };
            
            return {
              id: device.id,
              name: device.name,
              brand: device.brand,
              weight: device.dimensions?.weight || 0,
              connection_type: device.technical?.connectivity?.join(', ') || '',
              sensor: device.technical?.sensor || '',
              imageUrl: device.imageUrl || '',
              svgData: svgData
            };
          });
          console.log('Loaded mouse list using API data, retrieved', mice.value.length, 'mice');
        } else {
          // Data format is correct but the list is empty, also use hardcoded data as a supplement
          console.warn('API returned device list is empty or incorrectly formatted, using hardcoded data as a fallback');
          
          // 使用硬编码数据并添加一个标志来区分
          mice.value = hardcodedMice.map(device => ({
            id: device.id,
            name: device.name,
            brand: device.brand,
            weight: device.dimensions?.weight || 0,
            connection_type: device.technical?.connectivity?.join(', ') || '',
            sensor: device.technical?.sensor || '',
            imageUrl: device.imageUrl || '',
            isHardcoded: true // 添加标志
          }));
          
          console.log('Loaded mouse list from hardcoded data, count:', mice.value.length);
        }
      } else {
        // API response format is incorrect, fallback to hardcoded data
        console.warn('API response format error, using hardcoded data as fallback');
        console.warn('Response details:', response);
        
        mice.value = hardcodedMice.map(device => ({
          id: device.id,
          name: device.name,
          brand: device.brand,
          weight: device.dimensions?.weight || 0,
          connection_type: device.technical?.connectivity?.join(', ') || '',
          sensor: device.technical?.sensor || '',
          imageUrl: device.imageUrl || '',
          isHardcoded: true // 添加标志
        }));
      }
    } catch (error) {
      console.warn('Network error: Unable to connect to API, using hardcoded data as fallback', error);
      
      // Use hardcoded data
      mice.value = hardcodedMice.map(device => ({
        id: device.id,
        name: device.name,
        brand: device.brand,
        weight: device.dimensions?.weight || 0,
        connection_type: device.technical?.connectivity?.join(', ') || '',
        sensor: device.technical?.sensor || '',
        imageUrl: device.imageUrl || '',
        isHardcoded: true // 添加标志
      }));
      
      console.log('Loaded mouse list using hardcoded data, retrieved', mice.value.length, 'mice');
    }
  } catch (error) {
    console.error('Failed to load mouse data:', error);
    mice.value = []; // Set to empty array to avoid referencing null
  } finally {
    loading.value = false;
  }
});

// Filter mouse data
const filteredMice = computed(() => {
  if (!searchQuery.value) return mice.value;

  const query = searchQuery.value.toLowerCase();
  return mice.value.filter(
    (mouse) => mouse.name.toLowerCase().includes(query) || mouse.brand.toLowerCase().includes(query)
  );
});

// Recently viewed mice
const recentlyViewedMice = computed(() => {
  // Get recently viewed mice from store, instead of returning the first 3
  return comparisonStore.recentlyViewedMice;
});

// Select mouse
const selectMouse = (mouse) => {
  // Only select the mouse, but don't automatically add it to the comparison list to avoid adding it when clicking for details
  // Add mouse to recently viewed list
  comparisonStore.addToRecentlyViewed(mouse);
  // Notify parent component of selection change
  emit('select', mouse);
};

// Toggle mouse selection state
const toggleMouse = (mouse) => {
  if (isSelected(mouse)) {
    // If the mouse is already selected, but we don't want to remove it
    // Instead, display its unselectable state and keep it selected
    // So we don't remove it here, just add it to the recently viewed list
    comparisonStore.addToRecentlyViewed(mouse);
  } else {
    // Check if maximum selection count has been reached
    if (comparisonStore.selectedMice.length < props.maxSelection) {
      comparisonStore.addMouse(mouse);
      // Add mouse to recently viewed list
      comparisonStore.addToRecentlyViewed(mouse);
    }
  }
  // Notify parent component of selection change
  emit('select', comparisonStore.selectedMice);
};

// Check if the mouse is already selected
const isSelected = (mouse) => {
  return comparisonStore.isMouseSelected(mouse.id);
};
</script>

<style scoped>
.mouse-selector {
  color: var(--claude-text-light);
}

.search-container :deep(.el-input__inner) {
  background-color: var(--claude-bg-medium);
  border-color: var(--claude-border-light);
  color: var(--claude-text-light);
}

:deep(.el-skeleton__item) {
  background: var(--claude-bg-light);
}

/* Beautify scrollbar */
.mouse-list {
  scrollbar-width: thin;
  scrollbar-color: var(--claude-border-light) var(--claude-bg-dark);
}

.mouse-list::-webkit-scrollbar {
  width: 6px;
}

.mouse-list::-webkit-scrollbar-track {
  background: var(--claude-bg-dark);
}

.mouse-list::-webkit-scrollbar-thumb {
  background-color: var(--claude-border-light);
  border-radius: 3px;
}
</style>
