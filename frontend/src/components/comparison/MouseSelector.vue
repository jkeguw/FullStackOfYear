<template>
  <div class="mouse-selector">
    <div class="grid grid-cols-1 gap-4">
      <!-- 搜索栏 -->
      <div class="search-container mb-4">
        <el-input
          v-model="searchQuery"
          placeholder="搜索鼠标..."
          :prefix-icon="Search"
          clearable
        />
      </div>
      
      <!-- 鼠标列表 -->
      <div class="mouse-list bg-[#1A1A1A] border border-[#333333] rounded-lg p-4 max-h-96 overflow-y-auto">
        
        <div v-if="loading" class="text-center py-4">
          <el-skeleton :rows="5" animated />
        </div>
        
        <div v-else-if="filteredMice.length === 0" class="text-center py-8 text-gray-400">
          没有找到匹配的鼠标
        </div>
        
        <div v-else class="space-y-3">
          <div 
            v-for="mouse in filteredMice" 
            :key="mouse.id"
            class="mouse-item p-3 rounded-lg bg-[#242424] hover:bg-[#333333] cursor-pointer border border-[#333333] flex items-center"
            @click="toggleMouse(mouse)"
          >
            <div class="flex-shrink-0 w-16 h-16 bg-[#1A1A1A] rounded-lg flex items-center justify-center mr-3 overflow-hidden">
              <img 
                v-if="mouse.image_url" 
                :src="mouse.image_url" 
                :alt="mouse.name" 
                class="max-w-full max-h-full object-contain"
              >
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
      
      <!-- 最近查看的鼠标 -->
      <div v-if="recentlyViewedMice.length > 0" class="recently-viewed bg-[#1A1A1A] border border-[#333333] rounded-lg p-4">
        <h3 class="text-lg font-medium text-white mb-4">最近查看</h3>
        
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
          <div 
            v-for="mouse in recentlyViewedMice" 
            :key="mouse.id"
            class="recent-mouse-item p-2 rounded-lg bg-[#242424] hover:bg-[#333333] cursor-pointer border border-[#333333] flex flex-col items-center"
            @click="toggleMouse(mouse)"
          >
            <div class="w-12 h-12 bg-[#1A1A1A] rounded-lg flex items-center justify-center mb-2 overflow-hidden">
              <img 
                v-if="mouse.image_url" 
                :src="mouse.image_url" 
                :alt="mouse.name" 
                class="max-w-full max-h-full object-contain"
              >
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
import { ref, computed, onMounted, defineProps, defineEmits } from 'vue'
import { Search, Plus, Minus, Mouse } from '@element-plus/icons-vue'
import { useComparisonStore, type MouseDevice } from '@/stores'

// 定义props
const props = defineProps({
  initialSelectedMice: {
    type: Array as () => MouseDevice[],
    default: () => []
  },
  maxSelection: {
    type: Number,
    default: 3
  }
})

// 定义事件
const emit = defineEmits(['select', 'cancel'])

// 使用store
const comparisonStore = useComparisonStore()

// 搜索查询
const searchQuery = ref('')
const loading = ref(true)

// 模拟鼠标数据
const mice = ref([
  { id: 1, name: 'Logitech G Pro X Superlight', brand: 'Logitech', weight: 63, connection_type: '无线', sensor: 'HERO', image_url: '' },
  { id: 2, name: 'Razer Viper Ultimate', brand: 'Razer', weight: 74, connection_type: '无线', sensor: 'Focus+', image_url: '' },
  { id: 3, name: 'Zowie EC2', brand: 'Zowie', weight: 90, connection_type: '有线', sensor: '3360', image_url: '' },
  { id: 4, name: 'SteelSeries Prime', brand: 'SteelSeries', weight: 69, connection_type: '有线', sensor: 'TrueMove Pro', image_url: '' },
  { id: 5, name: 'Glorious Model O', brand: 'Glorious', weight: 67, connection_type: '有线', sensor: 'Pixart PMW3360', image_url: '' },
  { id: 6, name: 'Vaxee Outset AX', brand: 'Vaxee', weight: 76, connection_type: '有线', sensor: '3389', image_url: '' },
  { id: 7, name: 'Pulsar Xlite V2', brand: 'Pulsar', weight: 59, connection_type: '无线', sensor: 'PAW3370', image_url: '' }
])

// 加载模拟数据
onMounted(() => {
  // 初始化已选择的鼠标
  if (props.initialSelectedMice && props.initialSelectedMice.length > 0) {
    // 如果有初始选择的鼠标，确保在store中也标记为已选择
    props.initialSelectedMice.forEach(mouse => {
      if (!comparisonStore.isMouseSelected(mouse.id)) {
        comparisonStore.addMouse(mouse)
      }
    })
  }
  
  setTimeout(() => {
    loading.value = false
  }, 800)
})

// 过滤鼠标数据
const filteredMice = computed(() => {
  if (!searchQuery.value) return mice.value
  
  const query = searchQuery.value.toLowerCase()
  return mice.value.filter(mouse => 
    mouse.name.toLowerCase().includes(query) || 
    mouse.brand.toLowerCase().includes(query)
  )
})

// 最近查看的鼠标
const recentlyViewedMice = computed(() => {
  // 从store中获取最近查看的鼠标，而不是固定返回前3个
  return comparisonStore.recentlyViewedMice
})

// 选择鼠标
const selectMouse = (mouse) => {
  // 只选择鼠标，但不自动添加到比较列表，避免点击看详情就添加的问题
  // 将鼠标添加到最近查看列表
  comparisonStore.addToRecentlyViewed(mouse)
  // 通知父组件选择变化
  emit('select', mouse)
}

// 切换鼠标选择状态
const toggleMouse = (mouse) => {
  if (isSelected(mouse)) {
    comparisonStore.removeMouse(mouse.id)
  } else {
    // 检查是否已达到最大选择数量
    if (comparisonStore.selectedMice.length < props.maxSelection) {
      comparisonStore.addMouse(mouse)
    }
  }
  // 通知父组件选择变化
  emit('select', comparisonStore.selectedMice)
}

// 检查鼠标是否已被选择
const isSelected = (mouse) => {
  return comparisonStore.isMouseSelected(mouse.id)
}
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

/* 美化滚动条 */
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