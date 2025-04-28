#!/bin/bash

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}正在停止所有容器...${NC}"
docker-compose down

echo -e "${YELLOW}正在重构并启动所有容器...${NC}"
docker-compose up -d --build

echo -e "${YELLOW}等待容器启动...${NC}"
sleep 10

echo -e "${YELLOW}显示容器状态:${NC}"
docker-compose ps

echo -e "${YELLOW}显示后端日志:${NC}"
docker-compose logs backend

echo -e "${GREEN}重启完成！${NC}"
echo -e "${YELLOW}前端访问地址: http://localhost${NC}"
echo -e "${YELLOW}后端API地址: http://localhost:8081/api${NC}"
echo -e "${YELLOW}Mongo Express: http://localhost:8082 (用户名: admin, 密码: pass)${NC}"