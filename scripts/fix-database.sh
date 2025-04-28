#!/bin/bash

echo "检查MongoDB容器状态..."
MONGO_RUNNING=$(docker ps | grep mongodb | wc -l)

if [ "$MONGO_RUNNING" -eq 0 ]; then
  echo "MongoDB容器未运行，请先执行 docker-compose up -d"
  exit 1
fi

echo "检查MongoDB数据库..."
# 确保MongoDB初始化脚本已执行
echo "执行MongoDB初始化脚本..."
docker-compose exec mongodb mongo --username root --password example --authenticationDatabase admin --eval "db = db.getSiblingDB('cpc'); printjson(db.devices.count())"

# 如果设备集合为空，执行初始化脚本
echo "导入初始设备数据..."
docker-compose exec -T mongodb mongosh --username root --password example --authenticationDatabase admin cpc < mongodb-init/init-mongo.js

# 执行API修复（禁用DefaultService）
echo "修复完成。现在重启后端容器..."
docker-compose restart backend

echo "等待服务启动..."
sleep 5

echo "检查设备API..."
curl -s http://localhost:8081/api/devices | grep -o "devices"

echo "修复过程完成！请刷新页面查看鼠标数据。"