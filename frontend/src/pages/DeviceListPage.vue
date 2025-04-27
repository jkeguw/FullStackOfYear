<template>
  <div class="device-list-page">
    <div class="container mx-auto p-4">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-bold">外设数据库</h1>
        <div>
          <el-button 
            type="primary" 
            @click="handleCreateDevice" 
            v-if="hasAdminPermission"
          >
            添加设备
          </el-button>
        </div>
      </div>
      
      <!-- 筛选器 -->
      <el-card class="mb-6">
        <div class="filters">
          <el-form :model="filters" :inline="!isMobile" class="flex flex-wrap gap-4">
            <el-form-item label="设备类型" class="w-full sm:w-auto">
              <el-select v-model="filters.type" clearable placeholder="全部类型" @change="fetchData" class="w-full sm:w-auto">
                <el-option label="鼠标" value="mouse"></el-option>
                <el-option label="键盘" value="keyboard"></el-option>
                <el-option label="显示器" value="monitor"></el-option>
                <el-option label="鼠标垫" value="mousepad"></el-option>
                <el-option label="配件" value="accessory"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="品牌" class="w-full sm:w-auto">
              <el-select v-model="filters.brand" clearable placeholder="全部品牌" @change="fetchData" class="w-full sm:w-auto">
                <el-option v-for="brand in brands" :key="brand" :label="brand" :value="brand"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="排序" class="w-full xs:w-1/2 sm:w-auto">
              <el-select v-model="filters.sortBy" @change="fetchData" class="w-full">
                <el-option label="最新上架" value="createdAt"></el-option>
                <el-option label="名称" value="name"></el-option>
                <el-option label="品牌" value="brand"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="顺序" class="w-full xs:w-1/2 sm:w-auto">
              <el-select v-model="filters.sortOrder" @change="fetchData" class="w-full">
                <el-option label="降序" value="desc"></el-option>
                <el-option label="升序" value="asc"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item class="w-full sm:w-auto sm:ml-auto">
              <div class="flex gap-2">
                <el-button type="primary" @click="fetchData">筛选</el-button>
                <el-button @click="resetFilters">重置</el-button>
              </div>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
      
      <!-- 设备列表 -->
      <div v-loading="loading">
        <div v-if="devices.length === 0 && !loading" class="empty-state text-center py-16">
          <el-empty description="暂无设备数据"></el-empty>
        </div>
        
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          <el-card 
            v-for="device in devices" 
            :key="device.id" 
            class="device-card hover:shadow-lg transition-shadow"
            @click="viewDeviceDetails(device.id)"
          >
            <div class="device-card-content">
              <div class="device-image h-40 flex items-center justify-center mb-4">
                <el-image 
                  v-if="device.imageUrl" 
                  :src="device.imageUrl" 
                  fit="contain"
                  class="max-h-full"
                  :preview-src-list="[device.imageUrl]"
                ></el-image>
                <div v-else class="text-center text-gray-400">
                  <i class="el-icon-picture-outline text-3xl"></i>
                  <div>无图片</div>
                </div>
              </div>
              
              <div class="device-info">
                <div class="flex justify-between items-start">
                  <div>
                    <div class="text-lg font-bold line-clamp-1">{{ device.name }}</div>
                    <div class="text-gray-500">{{ device.brand }}</div>
                  </div>
                  <el-tag size="small">{{ getDeviceTypeName(device.type) }}</el-tag>
                </div>
                
                <div v-if="device.description" class="mt-2 text-sm text-gray-600 line-clamp-2">
                  {{ device.description }}
                </div>
                
                <div class="mt-4 flex justify-between items-center">
                  <div class="text-xs text-gray-400">
                    添加时间: {{ formatDate(device.createdAt) }}
                  </div>
                  <el-button type="primary" size="small" @click.stop="viewDeviceDetails(device.id)">
                    查看详情
                  </el-button>
                </div>
              </div>
            </div>
          </el-card>
        </div>
        
        <!-- 分页 -->
        <div class="pagination mt-6 flex justify-center overflow-x-auto py-2">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[12, 24, 48, 96]"
            :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next'"
            :small="isMobile"
            :total="pagination.total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          ></el-pagination>
        </div>
      </div>
    </div>
    
    <!-- 创建设备对话框 -->
    <el-dialog 
      v-model="deviceDialogVisible" 
      :title="isEditMode ? '编辑设备' : '添加设备'"
      width="80%"
      :before-close="closeDeviceDialog"
    >
      <device-form 
        :device-id="currentDeviceId" 
        @saved="handleDeviceSaved" 
        @deleted="handleDeviceDeleted"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { format } from 'date-fns'
import { useWindowSize } from '@vueuse/core'
import DeviceForm from '@/components/form/DeviceForm.vue'
import { useDevice } from '@/composables/useDevice'
import type { Device } from '@/api/device'

// 路由
const router = useRouter()

// 响应式设计
const { width } = useWindowSize()
const isMobile = computed(() => width.value < 640)

// 设备钩子
const { 
  devices, 
  deviceLoading: loading, 
  devicePagination: pagination,
  getDeviceTypeName,
  fetchDevices,
  removeDevice
} = useDevice()

// 筛选条件
const filters = reactive({
  type: '',
  brand: '',
  sortBy: 'createdAt',
  sortOrder: 'desc'
})

// 品牌列表（实际应该从API获取）
const brands = ref<string[]>(['Logitech', 'Razer', 'HyperX', 'Corsair', 'SteelSeries', 'Zowie', 'Glorious', 'Pulsar', 'Vaxee', 'Endgame Gear'])

// 对话框控制
const deviceDialogVisible = ref(false)
const currentDeviceId = ref<string | undefined>(undefined)
const isEditMode = computed(() => !!currentDeviceId.value)

// 权限检查（实际应该从用户状态获取）
const hasAdminPermission = ref(true)

// 生命周期钩子
onMounted(() => {
  fetchData()
})

// 获取设备数据
const fetchData = async () => {
  await fetchDevices({
    page: pagination.page,
    pageSize: pagination.pageSize,
    type: filters.type || undefined,
    brand: filters.brand || undefined,
    sortBy: filters.sortBy,
    sortOrder: filters.sortOrder
  })
}

// 重置筛选器
const resetFilters = () => {
  filters.type = ''
  filters.brand = ''
  filters.sortBy = 'createdAt'
  filters.sortOrder = 'desc'
  fetchData()
}

// 分页处理
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  fetchData()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchData()
}

// 查看设备详情
const viewDeviceDetails = (id: string) => {
  router.push(`/devices/${id}`)
}

// 创建设备
const handleCreateDevice = () => {
  currentDeviceId.value = undefined
  deviceDialogVisible.value = true
}

// 编辑设备
const handleEditDevice = (id: string) => {
  currentDeviceId.value = id
  deviceDialogVisible.value = true
}

// 关闭设备对话框
const closeDeviceDialog = () => {
  deviceDialogVisible.value = false
  currentDeviceId.value = undefined
}

// 设备保存成功处理
const handleDeviceSaved = (device: Device) => {
  deviceDialogVisible.value = false
  ElMessage.success(`设备 ${device.name} 已${isEditMode.value ? '更新' : '创建'}`)
  fetchData() // 刷新列表
}

// 设备删除成功处理
const handleDeviceDeleted = () => {
  deviceDialogVisible.value = false
  ElMessage.success('设备已删除')
  fetchData() // 刷新列表
}

// 格式化日期
const formatDate = (date: string | Date) => {
  return format(new Date(date), 'yyyy-MM-dd')
}
</script>

<style scoped>
.device-card {
  cursor: pointer;
  transition: all 0.3s ease;
}

.device-card:hover {
  transform: translateY(-5px);
}

.device-image {
  overflow: hidden;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>