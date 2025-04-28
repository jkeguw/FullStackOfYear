FROM node:18-alpine AS build

WORKDIR /app

COPY package*.json ./
RUN npm install --legacy-peer-deps

COPY . .

RUN npm run build

FROM nginx:alpine

# 安装wget用于健康检查
RUN apk add --no-cache wget

COPY --from=build /app/dist /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

# 添加健康检查文件
RUN echo "healthy" > /usr/share/nginx/html/health

CMD ["nginx", "-g", "daemon off;"]
