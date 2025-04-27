#!/bin/bash
set -e

# Colors for nice output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Preparing to shut down all containers...${NC}"

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

# Display running containers before shutdown
echo -e "${YELLOW}Currently running containers:${NC}"
docker-compose ps

# Ask for confirmation
read -p "Are you sure you want to shut down all containers? (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Shutdown cancelled.${NC}"
    exit 0
fi

# Option to save data
echo
read -p "Do you want to keep the data volumes? (y/n): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}Shutting down containers while preserving data volumes...${NC}"
    docker-compose down
    echo -e "${GREEN}Containers have been shut down. Data volumes are preserved.${NC}"
    echo -e "${YELLOW}To restart the application later, run: docker-compose up -d${NC}"
else
    echo -e "${RED}WARNING: This will delete ALL data including the database and uploaded files!${NC}"
    read -p "Are you REALLY sure you want to delete all data? (type 'DELETE' to confirm): " confirm
    if [[ "$confirm" != "DELETE" ]]; then
        echo -e "${YELLOW}Data deletion cancelled. Shutting down containers while preserving data...${NC}"
        docker-compose down
        echo -e "${GREEN}Containers have been shut down. Data volumes are preserved.${NC}"
    else
        echo -e "${RED}Shutting down containers and removing all data volumes...${NC}"
        docker-compose down -v
        echo -e "${GREEN}Containers and data volumes have been removed.${NC}"
    fi
fi

# Show cleanup message
echo -e "${YELLOW}Checking for any unused Docker resources that can be cleaned up:${NC}"
echo -e "${YELLOW}Unused networks:${NC}"
docker network ls --filter "dangling=true" -q | wc -l

echo -e "${YELLOW}Unused volumes:${NC}"
docker volume ls --filter "dangling=true" -q | wc -l

echo
echo -e "${YELLOW}If you want to clean up unused resources, you can run:${NC}"
echo -e "  - docker system prune --volumes -f  # Remove all unused containers, networks, volumes"
echo -e "  - docker volume prune -f            # Remove just the unused volumes"
echo

echo -e "${GREEN}Shutdown complete!${NC}"