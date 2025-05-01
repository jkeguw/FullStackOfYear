# TypeScript自动修复脚本

这组脚本用于自动修复前端TypeScript类型错误。

## 脚本说明

### 1. 一键修复（推荐）

执行以下命令一键修复所有TypeScript错误：

```bash
./scripts/run-all-fixes.sh
```

这会依次执行以下所有步骤，并自动备份原始代码。

### 2. 分步执行

如果需要更细粒度的控制，可以分步执行：

#### 基础TypeScript错误修复

```bash
./scripts/fix-ts-errors.sh
```

此脚本执行基本的文本替换和修复，解决简单的类型错误。

#### 安装并使用ts-morph

```bash
./scripts/install-ts-morph.sh
cd frontend && npx node ../scripts/fix-with-ts-morph.js
```

ts-morph提供高级AST解析能力，可以更智能地修复复杂类型问题。

#### 增强ESLint配置

```bash
./scripts/enhanced-eslint.sh
```

设置增强的ESLint规则，帮助统一代码风格并防止类型错误。

#### 运行ESLint自动修复

```bash
cd frontend && npm run lint:fix
```

让ESLint自动修复可修复的问题。

## 常见问题

1. **部分错误仍需手动修复**：某些复杂的类型错误可能需要手动修复，特别是涉及自定义类型的情况。

2. **类型不一致**：如果API返回格式与前端类型定义不一致，需要修改API响应处理代码或更新类型定义。

3. **命名约定**：所有属性名应使用camelCase风格（如`imageUrl`而非`image_url`）。