<template>
  <div class="user-device-list-page">
    <!-- 功能介绍通知 -->
    <notifications
      v-if="showWelcomeNotification"
      type="info"
      title="欢迎使用我的设备配置管理"
      message="在这里，您可以创建和管理自己的设备配置，记录各种设备的特定设置，并可以选择公开分享给其他用户。点击"创建新配置"按钮开始。"
      :actions="[
        { label: '不再显示', onClick: dismissWelcomeNotification },
        { label: '创建配置', onClick: () => { dismissWelcomeNotification(); handleCreateUserDevice(); }, primary: true }
      ]"
      :auto-close="false"
      @close="dismissWelcomeNotification"
    />
    
    <div class="container mx-auto p-4">
      <div class="flex justify-between items-center mb-6">
        <div>
          <h1 class="text-2xl font-bold">我的设备配置</h1>
          <p class="text-gray-500 text-sm mt-1">管理您的游戏设备和设置以获得最佳体验</p>
        </div>
        <div>
          <el-button type="primary" @click="handleCreateUserDevice">
            <i class="el-icon-plus mr-1"></i> 创建新配置
          </el-button>
        </div>
      </div>
      
      <!-- 筛选器 -->
      <el-card class="mb-6">
        <div class="filters">
          <el-form :model="filters" inline>
            <el-form-item label="是否公开">
              <el-select v-model="filters.isPublic" clearable placeholder="全部" @change="fetchData">
                <el-option label="公开" :value="true"></el-option>
                <el-option label="私密" :value="false"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="排序">
              <el-select v-model="filters.sortBy" @change="fetchData">
                <el-option label="创建时间" value="createdAt"></el-option>
                <el-option label="更新时间" value="updatedAt"></el-option>
                <el-option label="名称" value="name"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="顺序">
              <el-select v-model="filters.sortOrder" @change="fetchData">
                <el-option label="降序" value="desc"></el-option>
                <el-option label="升序" value="asc"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchData">筛选</el-button>
              <el-button @click="resetFilters">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
      
      <!-- 设备配置列表 -->
      <div v-loading="loading">
        <div v-if="userDevices.length === 0 && !loading" class="empty-state text-center py-16">
          <el-empty description="暂无设备配置"></el-empty>
          <el-button type="primary" class="mt-4" @click="handleCreateUserDevice">创建第一个配置</el-button>
        </div>
        
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <el-card 
            v-for="config in userDevices" 
            :key="config.id" 
            class="user-device-card hover:shadow-lg transition-shadow"
          >
            <template #header>
              <div class="flex justify-between items-center">
                <h3 class="text-lg font-bold">{{ config.name }}</h3>
                <div>
                  <el-tag v-if="config.isPublic" type="success" size="small">公开</el-tag>
                  <el-tag v-else type="info" size="small">私密</el-tag>
                  <span class="ml-2 text-xs text-gray-500">{{ formatDate(config.updatedAt) }}</span>
                </div>
              </div>
            </template>
            
            <div v-if="config.description" class="mb-4 text-gray-600 text-sm">
              {{ config.description }}
            </div>
            
            <div class="devices-list bg-gray-50 rounded-lg p-3">
              <div class="text-sm font-medium mb-2 text-gray-700">设备清单 ({{ config.devices.length }})</div>
              <div 
                v-for="device in config.devices" 
                :key="device.deviceId"
                class="device-item mb-3 p-3 bg-white border border-gray-100 rounded-md shadow-sm hover:shadow-md transition-shadow"
              >
                <div class="flex justify-between items-center">
                  <div class="flex items-center">
                    <el-tag class="mr-2" size="small" :type="getDeviceTagType(device.deviceType)">
                      {{ getDeviceTypeName(device.deviceType) }}
                    </el-tag>
                    <span class="font-medium">{{ device.deviceName }}</span>
                  </div>
                  <span class="text-sm text-gray-500">{{ device.deviceBrand }}</span>
                </div>
                
                <div v-if="device.settings && Object.keys(device.settings).length > 0" class="mt-3 bg-gray-50 p-2 rounded">
                  <div class="text-xs font-medium text-gray-500 mb-1">主要设置:</div>
                  <div class="flex flex-wrap gap-1">
                    <el-tag 
                      v-for="(value, key) in getFilteredSettings(device.settings)" 
                      :key="key"
                      size="small"
                      effect="plain"
                      :type="getSettingTagType(key)"
                    >
                      {{ getSettingDisplayName(key) }}: {{ formatSettingValue(value) }}
                    </el-tag>
                  </div>
                </div>
              </div>
            </div>
            
            <div class="mt-4 flex justify-center gap-2">
              <el-button size="small" type="primary" plain @click="handleEditUserDevice(config.id)">
                <i class="el-icon-edit mr-1"></i> 编辑配置
              </el-button>
              <el-button size="small" type="danger" plain @click="handleDeleteUserDevice(config.id)">
                <i class="el-icon-delete mr-1"></i> 删除
              </el-button>
            </div>
          </el-card>
        </div>
        
        <!-- 分页 -->
        <div class="pagination mt-6 flex justify-center">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[9, 18, 36]"
            layout="total, sizes, prev, pager, next"
            :total="pagination.total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          ></el-pagination>
        </div>
      </div>
    </div>
    
    <!-- 创建/编辑设备配置对话框 -->
    <el-dialog 
      v-model="userDeviceDialogVisible" 
      :title="isEditMode ? '编辑设备配置' : '创建设备配置'"
      width="80%"
      :before-close="closeUserDeviceDialog"
    >
      <user-device-form 
        :user-device-id="currentUserDeviceId" 
        @saved="handleUserDeviceSaved" 
        @canceled="closeUserDeviceDialog"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { format } from 'date-fns'
import { useDevice } from '@/composables/useDevice'
import type { UserDevice } from '@/api/device'
import UserDeviceForm from '@/components/form/UserDeviceForm.vue'
import Notifications from '@/components/common/Notifications.vue'

// 设备钩子
const { 
  userDevices, 
  userDeviceLoading: loading, 
  userDevicePagination: pagination,
  getDeviceTypeName,
  fetchUserDevices,
  removeUserDevice
} = useDevice()

// 筛选条件
const filters = reactive({
  isPublic: undefined as boolean | undefined,
  sortBy: 'updatedAt',
  sortOrder: 'desc'
})

// 对话框控制
const userDeviceDialogVisible = ref(false)
const currentUserDeviceId = ref<string | undefined>(undefined)
const isEditMode = computed(() => !!currentUserDeviceId.value)

// 欢迎通知控制
const showWelcomeNotification = ref(localStorage.getItem('hideDeviceConfigWelcome') !== 'true')

// 关闭欢迎通知
const dismissWelcomeNotification = () => {
  showWelcomeNotification.value = false
  localStorage.setItem('hideDeviceConfigWelcome', 'true')
}

// 生命周期钩子
onMounted(() => {
  fetchData()
})

// 获取数据
const fetchData = async () => {
  await fetchUserDevices({
    page: pagination.page,
    pageSize: pagination.pageSize,
    isPublic: filters.isPublic,
    sortBy: filters.sortBy,
    sortOrder: filters.sortOrder
  })
}

// 重置筛选器
const resetFilters = () => {
  filters.isPublic = undefined
  filters.sortBy = 'updatedAt'
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

// 创建设备配置
const handleCreateUserDevice = () => {
  currentUserDeviceId.value = undefined
  userDeviceDialogVisible.value = true
}

// 编辑设备配置
const handleEditUserDevice = (id: string) => {
  currentUserDeviceId.value = id
  userDeviceDialogVisible.value = true
}

// 关闭设备配置对话框
const closeUserDeviceDialog = () => {
  userDeviceDialogVisible.value = false
  currentUserDeviceId.value = undefined
}

// 设备配置保存成功处理
const handleUserDeviceSaved = (userDevice: UserDevice) => {
  userDeviceDialogVisible.value = false
  ElMessage.success(`设备配置 ${userDevice.name} 已${isEditMode.value ? '更新' : '创建'}`)
  fetchData() // 刷新列表
}

// 删除设备配置
const handleDeleteUserDevice = (id: string) => {
  ElMessageBox.confirm(
    '确定要删除此设备配置吗？删除后无法恢复！',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  )
  .then(async () => {
    const success = await removeUserDevice(id)
    if (success) {
      ElMessage.success('设备配置已删除')
      fetchData() // 刷新列表
    }
  })
  .catch(() => {
    // 取消删除，不做任何操作
  })
}

// 格式化日期
const formatDate = (date: string | Date) => {
  return format(new Date(date), 'yyyy-MM-dd HH:mm')
}

// 获取设备标签类型
const getDeviceTagType = (deviceType: string): string => {
  const types: Record<string, string> = {
    'mouse': 'primary',
    'keyboard': 'success',
    'monitor': 'warning',
    'mousepad': 'info',
    'accessory': 'danger'
  }
  return types[deviceType] || 'info'
}

// 过滤和限制显示的设置项
const getFilteredSettings = (settings: Record<string, any>): Record<string, any> => {
  // 定义重要的设置项列表（每种设备类型的主要设置）
  const importantSettings: Record<string, string[]> = {
    'mouse': ['dpi', 'pollingRate', 'liftOffDistance'],
    'keyboard': ['keyDelay', 'rgbLighting'],
    'monitor': ['brightness', 'contrast', 'colorMode'],
    'mousepad': ['notes']
  }
  
  // 获取当前设备类型的重要设置
  const result: Record<string, any> = {}
  
  // 最多显示4个设置项
  let count = 0
  for (const key in settings) {
    if (count >= 4) break
    
    if (settings[key] !== undefined && settings[key] !== null && settings[key] !== '') {
      result[key] = settings[key]
      count++
    }
  }
  
  return result
}

// 获取设置项显示名称
const getSettingDisplayName = (key: string): string => {
  const names: Record<string, string> = {
    'dpi': 'DPI',
    'pollingRate': '轮询率',
    'enhancedReportRate': '提升报告率',
    'liftOffDistance': '抬离距离',
    'debounce': '防抖',
    'keyDelay': '按键延迟',
    'rgbLighting': 'RGB灯效',
    'keyMapping': '按键映射',
    'brightness': '亮度',
    'contrast': '对比度',
    'colorMode': '颜色模式',
    'responseTime': '响应时间',
    'notes': '备注',
    'customSetting': '自定义'
  }
  return names[key] || key
}

// 获取设置项标签类型
const getSettingTagType = (key: string): string => {
  const types: Record<string, string> = {
    'dpi': 'primary',
    'pollingRate': 'success',
    'enhancedReportRate': 'warning',
    'liftOffDistance': 'info',
    'debounce': 'danger',
    'brightness': 'primary',
    'contrast': 'success',
    'colorMode': 'warning',
    'rgbLighting': 'success'
  }
  return types[key] || ''
}

// 格式化设置值
const formatSettingValue = (value: any): string => {
  if (typeof value === 'boolean') {
    return value ? '是' : '否'
  } else if (typeof value === 'number') {
    // 添加单位
    if (['dpi'].includes(String(value).toLowerCase())) {
      return `${value}`
    } else if (value > 100 && value % 125 === 0) { // 可能是轮询率
      return `${value}Hz`
    } else if (value <= 100) { // 可能是百分比
      return `${value}%`
    }
    return value.toString()
  } else {
    return String(value)
  }
}
</script>

<style scoped>
.user-device-card {
  transition: all 0.3s ease;
}

.user-device-card:hover {
  transform: translateY(-5px);
}

.device-item:hover {
  background-color: #f9f9f9;
}
</style>