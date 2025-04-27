<template>
  <div class="i18n-demo container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6">{{ $t('home.title') }}</h1>
    
    <div class="mb-8">
      <h2 class="text-xl font-semibold mb-4">{{ $t('common.language') }}</h2>
      <el-radio-group v-model="currentLanguage" @change="changeLanguage">
        <el-radio label="en">{{ $t('common.english') }}</el-radio>
        <el-radio label="zh">{{ $t('common.chinese') }}</el-radio>
      </el-radio-group>
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <!-- 表单示例 -->
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-semibold mb-4">{{ $t('profile.personal_info') }}</h3>
        <el-form :model="form" label-position="top">
          <el-form-item :label="$t('profile.username')">
            <el-input v-model="form.username" />
          </el-form-item>
          <el-form-item :label="$t('profile.email')">
            <el-input v-model="form.email" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary">{{ $t('common.save') }}</el-button>
            <el-button>{{ $t('common.cancel') }}</el-button>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 信息展示 -->
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-semibold mb-4">{{ $t('device.details') }}</h3>
        <div class="space-y-4">
          <div class="flex justify-between">
            <span class="text-gray-600">{{ $t('device.dimensions') }}:</span>
            <span>120 x 66 x 42 mm</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">{{ $t('device.weight') }}:</span>
            <span>85g</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">{{ $t('device.sensor') }}:</span>
            <span>PAW3395</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">{{ $t('device.buttons') }}:</span>
            <span>5</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">{{ $t('device.dpi') }}:</span>
            <span>400-26000 DPI</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 错误消息示例 -->
    <div class="mt-8 bg-white p-6 rounded-lg shadow">
      <h3 class="text-lg font-semibold mb-4">{{ $t('errors.required_field') }} ({{ $t('errors.server_error') }})</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <el-alert
          v-for="(error, index) in errorMessages"
          :key="index"
          :title="$t(error)"
          type="error"
          :closable="false"
          show-icon
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLanguage } from '@/i18n'

const { t, locale } = useI18n()
const currentLanguage = ref(locale.value)

// 表单数据
const form = ref({
  username: 'John Doe',
  email: 'john@example.com'
})

// 错误消息示例
const errorMessages = [
  'errors.required_field',
  'errors.invalid_email',
  'errors.password_mismatch',
  'errors.login_failed',
  'errors.server_error'
]

// 切换语言
const changeLanguage = (lang: 'en' | 'zh') => {
  setLanguage(lang)
}
</script>

<style scoped>
.i18n-demo {
  margin-bottom: 3rem;
}
</style>