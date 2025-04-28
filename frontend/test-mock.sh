#!/bin/bash

# 测试前端鼠标对比功能的脚本
echo "开始准备测试环境..."

# 备份原始main.ts
if [ ! -f src/main.ts.original ]; then
  cp src/main.ts src/main.ts.original
  echo "已备份原始main.ts"
fi

# 使用模拟版本替换main.ts
cp src/main.ts.mock src/main.ts
echo "已替换main.ts为模拟数据版本"

# 创建.env.local文件启用模拟数据
echo "VITE_USE_MOCK_DATA=true" > .env.local
echo "已创建.env.local文件启用模拟数据"

# 启动开发服务器
echo "启动开发服务器..."
npm run dev

# 当服务器停止时，恢复原始文件
trap cleanup EXIT
function cleanup {
  echo "恢复原始文件..."
  cp src/main.ts.original src/main.ts
  rm .env.local
  echo "清理完成！"
}