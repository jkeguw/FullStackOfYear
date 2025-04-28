module.exports = {
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
}
