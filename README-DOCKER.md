# Docker 容器化部署指南

本文档提供了如何使用 Docker 和 Docker Compose 部署完整应用栈的指南。

## 包含的服务

- **前端**: Vue.js 应用，通过 Nginx 提供服务
- **后端**: Go API 服务
- **MongoDB**: 数据存储
- **Redis**: 缓存和会话管理

## 先决条件

确保已安装以下软件:

- Docker
- Docker Compose

## 快速启动

使用提供的脚本快速部署所有服务:

```bash
./setup-containers.sh
```

该脚本会:
1. 检查必要的依赖
2. 创建默认的环境变量文件 (如果不存在)
3. 构建并启动所有容器
4. 验证所有服务是否正常运行

## 手动设置

如果需要手动设置，请按照以下步骤操作:

### 1. 设置环境变量

复制示例环境变量文件:

```bash
cp .env.example .env
```

编辑 `.env` 文件以设置实际凭证:

```
# MongoDB
MONGODB_URI=mongodb://root:example@mongodb:27017
# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
# JWT
JWT_SECRET=your-secret-key
# OAuth (set these in production)
GOOGLE_OAUTH_CLIENT_ID=your-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost/api/auth/google/callback
# SMTP (set these in production)
SMTP_USERNAME=your-username
SMTP_PASSWORD=your-password
```

### 2. 构建并启动容器

```bash
docker-compose build
docker-compose up -d
```

### 3. 访问服务

- 前端: http://localhost
- 后端 API: http://localhost/api
- MongoDB: localhost:27017
- Redis: localhost:6379

## 数据持久化

所有数据都使用 Docker 卷进行持久化:

- `mongodb_data`: MongoDB 数据
- `redis_data`: Redis 数据
- `backend_uploads`: 后端上传的文件

## 常用命令

```bash
# 查看日志
docker-compose logs -f

# 停止所有服务
docker-compose down

# 重启单个服务
docker-compose restart backend

# 仅在后台重建和重启某个服务
docker-compose up -d --build backend
```

## 生产部署注意事项

对于生产环境:

1. 修改 `docker-compose.yaml` 以移除开发专用的映射和配置
2. 使用强密码和安全的凭证
3. 考虑使用 Docker Swarm 或 Kubernetes 进行容器编排
4. 设置适当的资源限制
5. 实现监控和日志集中处理
6. 在 `.env` 文件中替换所有默认和示例凭证
7. 配置 HTTPS
8. 定期备份数据卷

## 故障排除

### 服务无法启动

查看特定服务的日志:

```bash
docker-compose logs backend
```

### 数据库连接问题

确保容器网络正常工作:

```bash
docker network ls
docker network inspect fullstackofyear_app-network
```

### 前端无法连接到后端

检查 Nginx 配置中的代理设置是否正确指向后端服务。

## 文件结构

- `docker-compose.yaml`: 定义所有服务和网络
- `backend/Dockerfile`: 构建后端服务的指令
- `frontend/Dockerfile`: 构建前端服务的指令
- `frontend/nginx.conf`: Nginx 配置文件
- `backend/config/config.docker.yaml`: 容器化环境的后端配置
- `setup-containers.sh`: 辅助设置脚本