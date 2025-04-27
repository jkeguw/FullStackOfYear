<template>
  <div class="i18n-demo-container p-4">
    <el-card class="mb-4">
      <template #header>
        <div class="card-header">
          <h1 class="text-2xl font-bold">{{ $t('common.language') }} {{ $t('common.settings') }}</h1>
          <language-switcher />
        </div>
      </template>
      
      <div class="translation-examples">
        <h2 class="text-xl font-semibold mb-4">{{ $t('common.app_name') }}</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- 通用翻译 -->
          <el-card class="mb-4">
            <template #header>
              <div class="card-header">
                <h2 class="text-lg font-semibold">{{ $t('common.common') }}</h2>
              </div>
            </template>
            <div class="space-y-2">
              <p><strong>Welcome:</strong> {{ $t('common.welcome') }}</p>
              <p><strong>Loading:</strong> {{ $t('common.loading') }}</p>
              <p><strong>Home:</strong> {{ $t('common.home') }}</p>
              <p><strong>Login:</strong> {{ $t('common.login') }}</p>
              <p><strong>Register:</strong> {{ $t('common.register') }}</p>
            </div>
          </el-card>
          
          <!-- 鼠标相关翻译 -->
          <el-card class="mb-4">
            <template #header>
              <div class="card-header">
                <h2 class="text-lg font-semibold">{{ $t('mouse.mice') }}</h2>
              </div>
            </template>
            <div class="space-y-2">
              <p><strong>Mouse:</strong> {{ $t('mouse.mouse') }}</p>
              <p><strong>Brand:</strong> {{ $t('mouse.brand') }}</p>
              <p><strong>Shape:</strong> {{ $t('mouse.shape') }}</p>
              <p><strong>Dimensions:</strong> {{ $t('mouse.dimensions') }}</p>
              <p><strong>Weight:</strong> {{ $t('mouse.weight') }}</p>
            </div>
          </el-card>
          
          <!-- 形状翻译 -->
          <el-card class="mb-4">
            <template #header>
              <div class="card-header">
                <h2 class="text-lg font-semibold">{{ $t('shapes.ergonomic') }}</h2>
              </div>
            </template>
            <div class="space-y-2">
              <p><strong>Ergonomic:</strong> {{ $t('shapes.ergonomic') }}</p>
              <p><strong>Symmetrical:</strong> {{ $t('shapes.symmetrical') }}</p>
              <p><strong>Ambidextrous:</strong> {{ $t('shapes.ambidextrous') }}</p>
              <p><strong>Asymmetric:</strong> {{ $t('shapes.asymmetric') }}</p>
            </div>
          </el-card>
          
          <!-- 手型尺寸翻译 -->
          <el-card class="mb-4">
            <template #header>
              <div class="card-header">
                <h2 class="text-lg font-semibold">{{ $t('hand_sizes.small') }}</h2>
              </div>
            </template>
            <div class="space-y-2">
              <p><strong>Small:</strong> {{ $t('hand_sizes.small') }}</p>
              <p><strong>Medium:</strong> {{ $t('hand_sizes.medium') }}</p>
              <p><strong>Large:</strong> {{ $t('hand_sizes.large') }}</p>
              <p><strong>Extra Large:</strong> {{ $t('hand_sizes.extra_large') }}</p>
            </div>
          </el-card>
        </div>
        
        <div class="mt-6">
          <h3 class="text-lg font-semibold mb-2">{{ $t('common.language') }}:</h3>
          <el-radio-group v-model="currentLang" @change="handleLanguageChange">
            <el-radio v-for="locale in SUPPORTED_LOCALES" :key="locale.code" :label="locale.code">
              {{ locale.name }}
            </el-radio>
          </el-radio-group>
        </div>
        
        <div class="mt-6">
          <h3 class="text-lg font-semibold mb-2">{{ $t('tools.tools') }}:</h3>
          <el-button type="primary" @click="showDateExample = !showDateExample">
            {{ showDateExample ? $t('common.hide') : $t('common.show') }} {{ $t('common.dates') }}
          </el-button>
          
          <div v-if="showDateExample" class="mt-4 p-4 bg-gray-50 rounded">
            <p><strong>{{ $t('common.today') }}:</strong> {{ $d(new Date(), 'long') }}</p>
            <p><strong>{{ $t('common.time') }}:</strong> {{ $d(new Date(), 'time') }}</p>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { SUPPORTED_LOCALES, setLocale } from '@/i18n'
import LanguageSwitcher from '@/components/common/LanguageSwitcher.vue'

const { locale } = useI18n()
const currentLang = ref(locale.value)
const showDateExample = ref(false)

const handleLanguageChange = (value: string) => {
  setLocale(value)
  currentLang.value = value
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>