#!/bin/bash

# 脚本路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "Initializing admin account..."

# 编译Go脚本
go build -o init_admin init_admin.go

# 检查编译是否成功
if [ $? -ne 0 ]; then
    echo "Error: Failed to compile init_admin.go"
    exit 1
fi

# 运行脚本
./init_admin

# 检查运行是否成功
if [ $? -ne 0 ]; then
    echo "Error: Failed to initialize admin account"
    exit 1
fi

# 清理编译的二进制文件
rm init_admin

echo "Admin account initialization completed successfully!"