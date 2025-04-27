<template>
  <div class="device-preferences-tab">
    <h2 class="text-2xl font-semibold mb-4">我的设备配置</h2>
    
    <div class="flex justify-between items-center mb-4">
      <el-input 
        v-model="searchQuery" 
        placeholder="搜索设备配置" 
        class="max-w-md"
        clearable
        @clear="searchQuery = ''"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      
      <el-button type="primary" @click="showAddDeviceDialog = true">
        <el-icon><Plus /></el-icon> 添加设备配置
      </el-button>
    </div>
    
    <el-table
      v-loading="loading"
      :data="filteredDevicePreferences"
      stripe
      border
      class="w-full mb-4"
      empty-text="暂无设备配置"
    >
      <el-table-column label="配置名称" prop="name" min-width="150" />
      <el-table-column label="设备数量" min-width="100">
        <template #default="{ row }">
          {{ row.devices?.length || 0 }} 个设备
        </template>
      </el-table-column>
      <el-table-column label="创建时间" min-width="150">
        <template #default="{ row }">
          {{ formatDate(row.createdAt) }}
        </template>
      </el-table-column>
      <el-table-column label="公开状态" min-width="120">
        <template #default="{ row }">
          <el-tag :type="row.isPublic ? 'success' : 'info'">
            {{ row.isPublic ? '公开' : '私有' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button size="small" @click="editDevicePreference(row)">
              编辑
            </el-button>
            <el-button 
              size="small" 
              :type="row.isPublic ? 'warning' : 'success'"
              @click="toggleVisibility(row)"
            >
              {{ row.isPublic ? '设为私有' : '设为公开' }}
            </el-button>
            <el-button size="small" type="danger" @click="confirmDelete(row)">
              删除
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    
    <!-- 添加设备配置对话框 -->
    <el-dialog
      v-model="showAddDeviceDialog"
      title="添加设备配置"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="deviceForm" label-position="top">
        <el-form-item label="配置名称" required>
          <el-input v-model="deviceForm.name" placeholder="例如：游戏设置" />
        </el-form-item>
        
        <el-form-item label="配置描述">
          <el-input 
            v-model="deviceForm.description" 
            type="textarea" 
            :rows="3"
            placeholder="描述这个设备配置的用途" 
          />
        </el-form-item>
        
        <el-form-item label="公开状态">
          <el-switch
            v-model="deviceForm.isPublic"
            :active-text="deviceForm.isPublic ? '公开' : '私有'"
            :inactive-text="deviceForm.isPublic ? '公开' : '私有'"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDeviceDialog = false">取消</el-button>
          <el-button type="primary" @click="addOrUpdateDevicePreference">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 删除确认对话框 -->
    <el-dialog
      v-model="showDeleteDialog"
      title="确认删除"
      width="400px"
    >
      <p>确定要删除设备配置 "{{ deviceToDelete?.name }}" 吗？此操作无法撤销。</p>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDeleteDialog = false">取消</el-button>
          <el-button type="danger" @click="deleteDevicePreference">
            确认删除
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import type { AddDevicePreferenceRequest, UpdateDevicePreferenceRequest } from '@/api/user'

interface DevicePreference {
  id: string
  name: string
  description?: string
  devices: Array<{
    deviceId: string
    deviceType: string
    settings?: Record<string, any>
  }>
  isPublic: boolean
  createdAt: string
  updatedAt: string
}

interface Props {
  loading: boolean
  preferences: DevicePreference[]
}

const props = withDefaults(defineProps<Props>(), {
  preferences: () => []
})

const emit = defineEmits<{
  (e: 'add', data: AddDevicePreferenceRequest): void
  (e: 'update', id: string, data: UpdateDevicePreferenceRequest): void
  (e: 'delete', id: string): void
}>()

// 状态
const searchQuery = ref('')
const showAddDeviceDialog = ref(false)
const showDeleteDialog = ref(false)
const isEditMode = ref(false)
const deviceToDelete = ref<DevicePreference | null>(null)

// 表单数据
const deviceForm = ref<{
  id?: string
  name: string
  description: string
  isPublic: boolean
}>({ 
  name: '', 
  description: '', 
  isPublic: false 
})

// 过滤后的设备配置列表
const filteredDevicePreferences = computed(() => {
  if (!searchQuery.value) return props.preferences
  
  const query = searchQuery.value.toLowerCase()
  return props.preferences.filter(item => 
    item.name.toLowerCase().includes(query) || 
    item.description?.toLowerCase().includes(query)
  )
})

// 编辑设备配置
const editDevicePreference = (item: DevicePreference) => {
  isEditMode.value = true
  deviceForm.value = {
    id: item.id,
    name: item.name,
    description: item.description || '',
    isPublic: item.isPublic
  }
  showAddDeviceDialog.value = true
}

// 添加或更新设备配置
const addOrUpdateDevicePreference = () => {
  if (!deviceForm.value.name.trim()) {
    ElMessage.warning('配置名称不能为空')
    return
  }
  
  if (isEditMode.value && deviceForm.value.id) {
    // 更新
    emit('update', deviceForm.value.id, {
      name: deviceForm.value.name,
      description: deviceForm.value.description,
      isPublic: deviceForm.value.isPublic
    })
  } else {
    // 添加
    emit('add', {
      name: deviceForm.value.name,
      description: deviceForm.value.description,
      isPublic: deviceForm.value.isPublic,
      devices: []
    })
  }
  
  showAddDeviceDialog.value = false
  resetForm()
}

// 切换可见性
const toggleVisibility = (item: DevicePreference) => {
  emit('update', item.id, {
    isPublic: !item.isPublic
  })
}

// 确认删除
const confirmDelete = (item: DevicePreference) => {
  deviceToDelete.value = item
  showDeleteDialog.value = true
}

// 删除设备配置
const deleteDevicePreference = () => {
  if (deviceToDelete.value) {
    emit('delete', deviceToDelete.value.id)
    showDeleteDialog.value = false
    deviceToDelete.value = null
  }
}

// 重置表单
const resetForm = () => {
  deviceForm.value = { 
    name: '', 
    description: '', 
    isPublic: false 
  }
  isEditMode.value = false
}

// 对话框关闭时重置表单
watch(showAddDeviceDialog, (newVal) => {
  if (!newVal) resetForm()
})
</script>

<style scoped>
.device-preferences-tab {
  padding: 20px 0;
}
</style>