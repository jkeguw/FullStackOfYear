#!/bin/bash

# 脚本路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "Initializing MongoDB collections and indexes..."

# 编译Go脚本
go build -o init_collections init_collections.go

# 检查编译是否成功
if [ $? -ne 0 ]; then
    echo "Error: Failed to compile init_collections.go"
    exit 1
fi

# 运行脚本
./init_collections

# 检查运行是否成功
if [ $? -ne 0 ]; then
    echo "Error: Failed to initialize collections and indexes"
    exit 1
fi

# 清理编译的二进制文件
rm init_collections

echo "MongoDB collections and indexes initialization completed successfully!"