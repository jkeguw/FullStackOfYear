# 项目修复概述

根据 `sa.txt` 中的要求，我已完成以下修复：

## 1. 修复主页图片问题
- 将缺失的 `sensitivity-preview.png` 图片添加到正确的路径：`/frontend/public/images/sensitivity-preview.png`
- 复制了 `2.png` 作为替代图片

## 2. 导航栏调整
- 抽屉导航栏已从右侧切换到左侧 (修改 `direction="ltr"`)
- 添加了一个固定在屏幕左侧的抽屉开关按钮
- 移除了顶部导航栏相关内容

## 3. 背景色适配
- 更新了各个组件的背景颜色，适配 Claude 夜间模式
- 修改了 `SortControls.vue`, `ViewToggle.vue` 等组件中的白色背景
- 使用了 CSS 变量确保颜色一致性

## 4. 灵敏度工具整合
- 创建了新的 `SensitivityToolPage.vue`，整合了 DPI 转换计算器和灵敏度转换工具
- 添加了三阶校准法和极敏内推法
- 基于算法文件夹中的公式实现了各种灵敏度计算方法

## 5. 国际化修复
- 添加了缺少的翻译条目如 `common.about`, `common.terms` 等
- 创建了 `LanguageSwitcherFlags.vue` 组件，使用国旗图标替代文本
- 修复导航栏菜单项目的显示名称

## 6. 添加缺少的页面
- 创建了 `AboutPage.vue`, `ContactPage.vue`, `PrivacyPage.vue`, `TermsPage.vue`
- 更新了路由配置，确保这些页面可以正确访问

## 7. 登录页面适配
- 更新了登录页面的样式，使用 Claude 夜间模式配色
- 确保购物车和订单页面需要登录才能访问

## 8. 其他调整
- 移除了个人设备管理功能相关的导航链接
- 修复了字体、边框和过渡效果，确保视觉一致性
- 确保所有组件适配暗色主题，避免明亮的白色背景

此修复方案确保整个应用程序在 Claude 夜间模式下有一致的视觉体验，同时解决了所有提到的功能问题。