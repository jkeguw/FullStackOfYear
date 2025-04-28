#!/bin/bash

# 停止并移除所有容器
echo "停止并移除所有容器..."
docker-compose down

# 移除所有卷
echo "移除所有卷..."
docker-compose down -v

# 重新构建前端
echo "重新构建前端容器..."
docker-compose build frontend

# 启动容器
echo "启动容器..."
docker-compose up -d

# 等待容器启动
echo "等待容器启动..."
sleep 15

# 查看容器状态
echo "容器状态:"
docker-compose ps

# 执行数据导入脚本
echo "初始化数据..."
docker-compose exec backend /bin/bash -c "cd /app/scripts && go build -o init_mongo init_mongo.go && ./init_mongo"

echo "服务准备完成!"
echo "前端地址: http://localhost:80"
echo "后端地址: http://localhost:8081"
echo "MongoDB管理界面: http://localhost:8082"