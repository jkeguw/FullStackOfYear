<template>
  <div class="settings-tab">
    <h2 class="text-2xl font-semibold mb-4">系统设置</h2>
    
    <el-form v-loading="loading" :model="form" label-position="top">
      <el-card class="mb-4">
        <template #header>
          <div class="flex justify-between items-center">
            <span>界面设置</span>
          </div>
        </template>
        
        <el-form-item label="主题">
          <el-select v-model="form.theme" class="w-full">
            <el-option label="浅色" value="light" />
            <el-option label="深色" value="dark" />
            <el-option label="跟随系统" value="system" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="语言">
          <el-select v-model="form.language" class="w-full">
            <el-option label="简体中文" value="zh-CN" />
            <el-option label="English" value="en-US" />
          </el-select>
        </el-form-item>
      </el-card>
      
      <el-card class="mb-4">
        <template #header>
          <div class="flex justify-between items-center">
            <span>通知设置</span>
          </div>
        </template>
        
        <el-form-item>
          <el-checkbox v-model="form.emailNotifications">接收邮件通知</el-checkbox>
        </el-form-item>
        
        <el-form-item>
          <el-checkbox v-model="form.pushNotifications">接收推送通知</el-checkbox>
        </el-form-item>
      </el-card>
      
      <el-card class="mb-4">
        <template #header>
          <div class="flex justify-between items-center">
            <span>内容显示设置</span>
          </div>
        </template>
        
        <el-form-item label="每页显示项目数">
          <el-select v-model="form.itemsPerPage" class="w-full">
            <el-option label="10" :value="10" />
            <el-option label="20" :value="20" />
            <el-option label="50" :value="50" />
            <el-option label="100" :value="100" />
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-checkbox v-model="form.showAdvancedOptions">显示高级选项</el-checkbox>
        </el-form-item>
      </el-card>
      
      <div class="flex justify-end mt-4">
        <el-button type="primary" @click="saveSettings">保存设置</el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { UpdateSettingsRequest } from '@/api/user'

interface Props {
  loading: boolean
  settings?: {
    theme: string
    language: string
    emailNotifications: boolean
    pushNotifications: boolean
    itemsPerPage: number
    showAdvancedOptions: boolean
  }
}

const props = withDefaults(defineProps<Props>(), {
  settings: () => ({
    theme: 'light',
    language: 'zh-CN',
    emailNotifications: true,
    pushNotifications: true,
    itemsPerPage: 20,
    showAdvancedOptions: false
  })
})

const emit = defineEmits<{
  (e: 'update', settings: UpdateSettingsRequest): void
}>()

const form = ref({
  theme: 'light',
  language: 'zh-CN',
  emailNotifications: true,
  pushNotifications: true,
  itemsPerPage: 20,
  showAdvancedOptions: false
})

// 监听props变化，更新表单
watch(
  () => props.settings,
  (newSettings) => {
    if (newSettings) {
      form.value = { ...newSettings }
    }
  },
  { immediate: true }
)

// 保存设置
const saveSettings = async () => {
  try {
    emit('update', form.value)
    ElMessage.success('设置已保存')
  } catch (error) {
    ElMessage.error('保存设置失败')
    console.error(error)
  }
}
</script>

<style scoped>
.settings-tab {
  padding: 20px 0;
}
</style>