<template>
  <el-config-provider :locale="elementLocale">
    <div class="app-container dark">
      <NavBar />

      <main class="main-content">
        <router-view />
      </main>

      <el-backtop :right="20" :bottom="20" />

      <Notifications v-if="false" title="通知" />
    </div>
  </el-config-provider>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import en from 'element-plus/es/locale/lang/en';

import NavBar from '@/components/layout/NavBar.vue';
import Notifications from '@/components/common/Notifications.vue';

const { locale } = useI18n();

// 根据当前语言设置 Element Plus 的语言
const elementLocale = computed(() => {
  if (locale.value.startsWith('zh')) {
    return zhCn;
  }
  return en;
});

// 设置主题变量
watch(
  () => locale.value,
  () => {
    document.documentElement.lang = locale.value;
  },
  { immediate: true }
);
</script>

<style>
:root {
  /* Claude夜间模式配色 */
  --claude-bg-dark: #121212;
  --claude-bg-medium: #1e1e1e;
  --claude-bg-light: #2a2a2a;
  --claude-text-white: #ffffff;
  --claude-text-light: #e0e0e0;
  --claude-text-muted: #a0a0a0;
  --claude-border-dark: #333333;
  --claude-border-light: #444444;
  --claude-primary-purple: #7d5af3;
  --claude-primary-purple-darker: #6a48e0;
  --claude-focus-ring: rgba(125, 90, 243, 0.5);
}

.dark {
  color-scheme: dark;
  --el-color-primary: var(--claude-primary-purple);
  --el-color-primary-light-3: var(--claude-primary-purple-darker);
  --el-color-primary-light-7: rgba(125, 90, 243, 0.2);
  --el-bg-color: var(--claude-bg-dark);
  --el-bg-color-overlay: var(--claude-bg-medium);
  --el-text-color-primary: var(--claude-text-white);
  --el-text-color-regular: var(--claude-text-light);
  --el-text-color-secondary: var(--claude-text-muted);
  --el-border-color: var(--claude-border-dark);
  --el-border-color-light: var(--claude-border-light);
  --el-fill-color-blank: var(--claude-bg-dark);
  --el-fill-color: var(--claude-bg-medium);
  --el-mask-color: rgba(0, 0, 0, 0.8);
  --el-box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.5);
}

html,
body {
  margin: 0;
  padding: 0;
  font-family:
    Inter,
    -apple-system,
    BlinkMacSystemFont,
    'Segoe UI',
    Roboto,
    Oxygen,
    Ubuntu,
    Cantarell,
    'Fira Sans',
    'Droid Sans',
    'Helvetica Neue',
    sans-serif;
  background-color: var(--claude-bg-dark);
  color: var(--claude-text-light);
  min-height: 100vh;
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  width: 100%;
  overflow-x: hidden;
}

/* 美化滚动条 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--claude-bg-dark);
}

::-webkit-scrollbar-thumb {
  background: var(--claude-border-light);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--claude-text-muted);
}

/* 调整element-plus组件样式 */
.el-button.is-text {
  color: var(--claude-text-light);
}

.el-input__inner,
.el-textarea__inner {
  background-color: var(--claude-bg-light) !important;
  border-color: var(--claude-border-light) !important;
  color: var(--claude-text-light) !important;
}

.el-input__inner:focus,
.el-textarea__inner:focus {
  border-color: var(--claude-primary-purple) !important;
  box-shadow: 0 0 0 2px var(--claude-focus-ring) !important;
}

.el-dropdown-menu {
  background-color: var(--claude-bg-medium) !important;
  border-color: var(--claude-border-dark) !important;
}

.el-dropdown-menu__item {
  color: var(--claude-text-light) !important;
}

.el-dropdown-menu__item:hover {
  background-color: var(--claude-bg-light) !important;
}

.el-select-dropdown {
  background-color: var(--claude-bg-medium) !important;
  border-color: var(--claude-border-dark) !important;
}

.el-select-dropdown__item {
  color: var(--claude-text-light) !important;
}

.el-select-dropdown__item.selected {
  color: var(--claude-primary-purple) !important;
}

.el-select-dropdown__item:hover {
  background-color: var(--claude-bg-light) !important;
}

.el-card {
  background-color: var(--claude-bg-medium) !important;
  border-color: var(--claude-border-dark) !important;
  color: var(--claude-text-light) !important;
}

.el-table {
  background-color: var(--claude-bg-medium) !important;
  color: var(--claude-text-light) !important;
}

.el-table th.el-table__cell {
  background-color: var(--claude-bg-light) !important;
}

.el-table tr {
  background-color: var(--claude-bg-medium) !important;
}

.el-table--striped .el-table__body tr.el-table__row--striped td.el-table__cell {
  background-color: var(--claude-bg-light) !important;
}

.el-pagination {
  background-color: transparent !important;
  color: var(--claude-text-light) !important;
}

.el-pagination .el-input__inner {
  background-color: var(--claude-bg-light) !important;
}

.el-pagination button:disabled {
  background-color: var(--claude-bg-medium) !important;
}

.el-dialog {
  background-color: var(--claude-bg-medium) !important;
  border-color: var(--claude-border-dark) !important;
}

.el-dialog__title {
  color: var(--claude-text-white) !important;
}

.el-form-item__label {
  color: var(--claude-text-light) !important;
}
</style>
