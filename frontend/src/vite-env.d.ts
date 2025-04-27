/// <reference types="vite/client" />

// 声明Vue文件模块
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// 导入ElementPlus类型
import { ElMessage, ElMessageBox, ElNotification, ElLoading } from 'element-plus';

// 扩展Window接口
declare global {
  interface Window {
    ElMessage: typeof ElMessage;
    ElMessageBox: typeof ElMessageBox;
    ElNotification: typeof ElNotification;
    ElLoading: typeof ElLoading;
  }
}

// 声明Element Plus组件库类型
declare module 'element-plus/dist/locale/zh-cn.mjs';
declare module 'element-plus/dist/locale/en.mjs';

// 声明第三方库
declare module '@element-plus/icons-vue';
