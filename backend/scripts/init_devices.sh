#!/bin/bash

# 设置工作目录
cd /app/scripts

# 编译并运行初始化脚本
echo "开始编译初始化脚本..."
go build -o init_mongo init_mongo.go

if [ $? -eq 0 ]; then
    echo "编译成功，开始初始化数据..."
    ./init_mongo
else
    echo "编译失败，请检查错误"
    exit 1
fi