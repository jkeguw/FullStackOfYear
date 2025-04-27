import type { App } from 'vue';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import zhCn from 'element-plus/dist/locale/zh-cn.mjs';
import en from 'element-plus/dist/locale/en.mjs';

// 图标
import * as ElementPlusIconsVue from '@element-plus/icons-vue';

// 创建全局 ElMessage 引用
import { ElMessage, ElMessageBox, ElNotification, ElLoading } from 'element-plus';

// 安装ElementPlus插件
export const setupElementPlus = (app: App) => {
  // 注册ElementPlus
  app.use(ElementPlus, {
    locale: zhCn, // 默认使用中文
  });

  // 注册所有图标
  for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component);
  }

  // 挂载全局方法
  window.ElMessage = ElMessage;
  window.ElMessageBox = ElMessageBox;
  window.ElNotification = ElNotification;
  window.ElLoading = ElLoading;

  return app;
};

// 暴露类型定义
declare global {
  interface Window {
    ElMessage: typeof ElMessage;
    ElMessageBox: typeof ElMessageBox;
    ElNotification: typeof ElNotification;
    ElLoading: typeof ElLoading;
  }
}