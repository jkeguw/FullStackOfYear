#\!/bin/bash

# Colors for nice output
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
NC="\033[0m" # No Color

echo -e "${YELLOW}Setting up containers...${NC}"
docker-compose up -d
echo -e "${GREEN}Containers started\!${NC}"
echo -e "${YELLOW}Frontend URL: http://localhost${NC}"
echo -e "${YELLOW}Backend API URL: http://localhost/api${NC}"
