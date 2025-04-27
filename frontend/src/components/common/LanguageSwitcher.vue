<template>
  <div class="language-switcher">
    <div class="claude-lang-selector">
      <button 
        v-for="locale in SUPPORTED_LOCALES" 
        :key="locale.code"
        :class="['claude-lang-btn', { active: getCurrentLanguage() === locale.code }]"
        @click="handleLanguageChange(locale.code)"
      >
        {{ locale.name }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { SUPPORTED_LOCALES, setLocale } from '@/i18n'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

// 获取当前语言
const getCurrentLanguage = () => {
  return locale.value
}

// 处理语言切换
const handleLanguageChange = (langCode: string) => {
  setLocale(langCode)
}
</script>

<style scoped>
.language-switcher {
  display: inline-block;
}

.claude-lang-selector {
  display: flex;
  background-color: var(--claude-bg-light);
  border-radius: 8px;
  padding: 2px;
  border: 1px solid var(--claude-border-dark);
}

.claude-lang-btn {
  background: transparent;
  color: var(--claude-text-muted);
  border: none;
  padding: 5px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
}

.claude-lang-btn.active {
  background-color: var(--claude-primary-purple);
  color: white;
  font-weight: 500;
}

.claude-lang-btn:not(.active):hover {
  color: var(--claude-text-white);
  background-color: var(--claude-bg-medium);
}
</style>