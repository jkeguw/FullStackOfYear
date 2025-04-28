FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o backend_app .

FROM alpine:latest

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache bash wget curl jq

COPY --from=builder /app/backend_app .
COPY --from=builder /app/config ./config
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/i18n ./i18n
COPY --from=builder /app/scripts /app/scripts

# 创建必要的目录
RUN mkdir -p /app/uploads

# 设置脚本权限
RUN chmod +x /app/scripts/init_devices.sh

# 添加健康检查端点
RUN echo '#!/bin/sh\necho "OK"' > /app/health.sh && chmod +x /app/health.sh

EXPOSE 8081

ENV CONFIG_FILE="/app/config/config.docker.yaml"

# 为了确保服务启动时健康检查能够通过，添加并发布/api/v1/health路由
CMD ["./backend_app"]