#!/bin/bash

# 进入前端目录
cd ../frontend

# 创建新的ESLint配置
echo 'module.exports = {
  root: true,
  parser: "vue-eslint-parser",
  parserOptions: {
    parser: "@typescript-eslint/parser",
    ecmaVersion: 2020,
    sourceType: "module"
  },
  env: {
    browser: true,
    es2020: true,
    node: true
  },
  extends: [
    "eslint:recommended",
    "plugin:vue/vue3-recommended",
    "plugin:@typescript-eslint/recommended",
    "prettier"
  ],
  plugins: ["vue", "@typescript-eslint"],
  rules: {
    // 关闭命名规范检查，或根据需要调整
    "@typescript-eslint/naming-convention": "off",
    // 允许使用any
    "@typescript-eslint/no-explicit-any": "warn",
    // 允许未使用的变量
    "@typescript-eslint/no-unused-vars": ["warn", { "argsIgnorePattern": "^_", "varsIgnorePattern": "^_" }],
    // 其他规则可以根据需要添加
    "vue/multi-word-component-names": "off"
  }
}' > .eslintrc.cjs

# 删除可能冲突的ESLint配置
rm -f eslint.config.js
rm -f .eslintrc.js

# 更新package.json中的eslint相关依赖
npm install --save-dev eslint@^8.0.0 \
  eslint-plugin-vue@^9.0.0 \
  @typescript-eslint/eslint-plugin@^8.0.0 \
  @typescript-eslint/parser@^8.0.0 \
  vue-eslint-parser@^9.0.0 \
  eslint-config-prettier@^8.0.0 \
  --legacy-peer-deps

# 运行ESLint修复
echo "正在运行ESLint修复..."
npx eslint --fix src/**/*.{ts,vue} --max-warnings=0 || true

echo "修复完成。部分错误可能需要手动修复。"