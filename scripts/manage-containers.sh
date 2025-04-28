#!/bin/bash

# Colors for nice output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed.${NC}"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: Docker Compose is not installed.${NC}"
    exit 1
fi

# Function to display menu
show_menu() {
    clear
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}    容器管理工具    ${NC}"
    echo -e "${BLUE}================================${NC}"
    echo -e "${GREEN}1.${NC} 启动所有容器"
    echo -e "${GREEN}2.${NC} 停止所有容器"
    echo -e "${GREEN}3.${NC} 删除所有容器（保留数据）"
    echo -e "${GREEN}4.${NC} 删除所有容器和数据"
    echo -e "${GREEN}5.${NC} 查看容器状态"
    echo -e "${GREEN}6.${NC} 查看容器日志"
    echo -e "${GREEN}7.${NC} 重构并重启容器"
    echo -e "${GREEN}8.${NC} 清理无用资源"
    echo -e "${GREEN}0.${NC} 退出"
    echo -e "${BLUE}================================${NC}"
    echo
}

# Function to start all containers
start_containers() {
    echo -e "${YELLOW}正在启动所有容器...${NC}"
    docker-compose up -d
    echo -e "${GREEN}容器已启动！${NC}"
    echo -e "${YELLOW}前端访问地址: http://localhost${NC}"
    echo -e "${YELLOW}后端API地址: http://localhost/api${NC}"
    echo
    read -p "按Enter键继续..."
}

# Function to stop all containers
stop_containers() {
    echo -e "${YELLOW}正在停止所有容器...${NC}"
    docker-compose stop
    echo -e "${GREEN}容器已停止！${NC}"
    echo
    read -p "按Enter键继续..."
}

# Function to remove containers but keep data
remove_containers() {
    echo -e "${YELLOW}正在移除所有容器（保留数据）...${NC}"
    docker-compose down
    echo -e "${GREEN}容器已移除，数据卷已保留！${NC}"
    echo
    read -p "按Enter键继续..."
}

# Function to remove containers and data
remove_containers_and_data() {
    echo -e "${RED}警告: 这将删除所有容器和数据！${NC}"
    read -p "确定要删除所有数据吗？（输入 'DELETE' 确认）: " confirm
    if [[ "$confirm" == "DELETE" ]]; then
        echo -e "${YELLOW}正在移除所有容器和数据...${NC}"
        docker-compose down -v
        echo -e "${GREEN}容器和数据已移除！${NC}"
    else
        echo -e "${YELLOW}操作已取消。${NC}"
    fi
    echo
    read -p "按Enter键继续..."
}

# Function to show container status
show_status() {
    echo -e "${YELLOW}容器状态:${NC}"
    docker-compose ps
    echo
    read -p "按Enter键继续..."
}

# Function to show container logs
show_logs() {
    echo -e "${YELLOW}请选择要查看日志的服务:${NC}"
    echo -e "${GREEN}1.${NC} backend"
    echo -e "${GREEN}2.${NC} frontend"
    echo -e "${GREEN}3.${NC} mongodb"
    echo -e "${GREEN}4.${NC} redis"
    echo -e "${GREEN}5.${NC} 所有服务"
    echo -e "${GREEN}0.${NC} 返回"
    read -p "请选择: " log_choice
    
    case $log_choice in
        1) 
            echo -e "${YELLOW}后端日志（按Ctrl+C退出）:${NC}"
            docker-compose logs -f backend
            ;;
        2) 
            echo -e "${YELLOW}前端日志（按Ctrl+C退出）:${NC}"
            docker-compose logs -f frontend
            ;;
        3) 
            echo -e "${YELLOW}MongoDB日志（按Ctrl+C退出）:${NC}"
            docker-compose logs -f mongodb
            ;;
        4) 
            echo -e "${YELLOW}Redis日志（按Ctrl+C退出）:${NC}"
            docker-compose logs -f redis
            ;;
        5) 
            echo -e "${YELLOW}所有服务日志（按Ctrl+C退出）:${NC}"
            docker-compose logs -f
            ;;
        0) 
            return
            ;;
        *) 
            echo -e "${RED}无效选择${NC}"
            ;;
    esac
    echo
    read -p "按Enter键继续..."
}

# Function to rebuild and restart containers
rebuild_containers() {
    echo -e "${YELLOW}请选择要重构的服务:${NC}"
    echo -e "${GREEN}1.${NC} backend"
    echo -e "${GREEN}2.${NC} frontend"
    echo -e "${GREEN}3.${NC} 所有服务"
    echo -e "${GREEN}0.${NC} 返回"
    read -p "请选择: " rebuild_choice
    
    case $rebuild_choice in
        1) 
            echo -e "${YELLOW}重构并重启后端...${NC}"
            docker-compose up -d --build backend
            ;;
        2) 
            echo -e "${YELLOW}重构并重启前端...${NC}"
            docker-compose up -d --build frontend
            ;;
        3) 
            echo -e "${YELLOW}重构并重启所有服务...${NC}"
            docker-compose up -d --build
            ;;
        0) 
            return
            ;;
        *) 
            echo -e "${RED}无效选择${NC}"
            ;;
    esac
    echo -e "${GREEN}重构完成！${NC}"
    echo
    read -p "按Enter键继续..."
}

# Function to clean unused resources
clean_resources() {
    echo -e "${YELLOW}正在清理未使用的Docker资源...${NC}"
    
    # Show what will be removed
    echo -e "${YELLOW}将移除的容器:${NC}"
    docker container ls -a --filter "status=exited" --filter "status=created" --format "{{.Names}}"
    
    echo -e "${YELLOW}将移除的镜像:${NC}"
    docker images -f "dangling=true" --format "{{.Repository}}:{{.Tag}}"
    
    echo -e "${YELLOW}将移除的网络:${NC}"
    docker network ls --filter "dangling=true" --format "{{.Name}}"
    
    echo -e "${YELLOW}将移除的卷:${NC}"
    docker volume ls --filter "dangling=true" --format "{{.Name}}"
    
    # Ask for confirmation
    read -p "是否继续清理？(y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}清理已取消。${NC}"
    else
        echo -e "${YELLOW}正在清理...${NC}"
        docker system prune --volumes -f
        echo -e "${GREEN}清理完成！${NC}"
    fi
    echo
    read -p "按Enter键继续..."
}

# Main loop
while true; do
    show_menu
    read -p "请选择操作: " choice
    
    case $choice in
        1) start_containers ;;
        2) stop_containers ;;
        3) remove_containers ;;
        4) remove_containers_and_data ;;
        5) show_status ;;
        6) show_logs ;;
        7) rebuild_containers ;;
        8) clean_resources ;;
        0) 
            echo -e "${GREEN}感谢使用，再见！${NC}"
            exit 0
            ;;
        *) 
            echo -e "${RED}无效选择，请重试${NC}"
            read -p "按Enter键继续..."
            ;;
    esac
done