<template>
  <!-- Language switcher flags hidden since we only support English now -->
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { locale } = useI18n();

const availableLanguages = [
  { code: 'en-US', name: 'English', flag: '🇺🇸' }
];

const currentLanguage = computed(() => locale.value);

function switchLanguage(lang: string) {
  locale.value = lang;
  localStorage.setItem('language', lang);
}
</script>

<style scoped>
.language-switcher {
  display: flex;
  gap: 0.5rem;
}

.lang-button {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  border-radius: 0.375rem;
  background-color: var(--claude-bg-light);
  border: 1px solid var(--claude-border-light);
  cursor: pointer;
  transition: all 0.2s ease;
}

.lang-button.active {
  border-color: var(--claude-primary-purple);
  background-color: rgba(125, 90, 243, 0.1);
}

.lang-button:hover {
  background-color: var(--claude-bg-medium);
}

.flag {
  font-size: 1.25rem;
}
</style>
