# Docker 部署指南

本文档详细说明如何使用 Docker 部署 Go-Vue-Admin 项目。

## 前置条件

- Docker 20.10+
- Docker Compose 2.0+
- Git

## 部署架构

```
                    ┌─────────────┐
                    │    Nginx    │
                    │   (前端)    │
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  Go Server  │
                    │   (后端)    │
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   MySQL     │
                    │  (数据库)   │
                    └─────────────┘
```

## Docker 镜像构建

### 后端镜像

1. Dockerfile
```dockerfile
# 构建阶段
FROM golang:1.16-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 复制配置文件
COPY --from=builder /app/configs/config.ini /app/configs/
# 复制编译后的二进制文件
COPY --from=builder /app/server /app/

EXPOSE 8080

CMD ["./server"]
```

2. 构建命令
```bash
docker build -t go-vue-admin-backend .
```

### 前端镜像

1. Dockerfile
```dockerfile
# 构建阶段
FROM node:14-alpine AS builder

WORKDIR /app

# 安装依赖
COPY frontend/package*.json ./
RUN npm install

# 复制源代码
COPY frontend .

# 构建应用
RUN npm run build

# 运行阶段
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

2. Nginx 配置
```nginx
# nginx.conf
server {
    listen 80;
    server_name localhost;

    root /usr/share/nginx/html;
    index index.html;

    # 前端路由配置
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

3. 构建命令
```bash
docker build -t go-vue-admin-frontend -f frontend/Dockerfile .
```

## Docker Compose 部署

### 配置文件

```yaml
# docker-compose.yml
version: '3'

services:
  frontend:
    image: go-vue-admin-frontend
    container_name: go-vue-admin-frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - go-vue-admin-network

  backend:
    image: go-vue-admin-backend
    container_name: go-vue-admin-backend
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=your_password
      - DB_NAME=go_vue_admin
    ports:
      - "8080:8080"
    depends_on:
      - mysql
    networks:
      - go-vue-admin-network

  mysql:
    image: mysql:5.7
    container_name: go-vue-admin-mysql
    environment:
      - MYSQL_ROOT_PASSWORD=your_password
      - MYSQL_DATABASE=go_vue_admin
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    networks:
      - go-vue-admin-network

networks:
  go-vue-admin-network:
    driver: bridge

volumes:
  mysql_data:
```

### 部署步骤

1. 准备环境
```bash
# 克隆项目
git clone https://github.com/yourusername/go-vue-admin.git
cd go-vue-admin

# 创建环境变量文件
cp .env.example .env
```

2. 修改配置
- 编辑 `.env` 文件设置环境变量
- 修改 `configs/config.ini` 配置数据库连接

3. 启动服务
```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

4. 验证部署
- 访问 http://localhost 查看前端页面
- 访问 http://localhost/api/health 检查后端服务

### 数据迁移

1. 执行数据库迁移
```bash
# 进入后端容器
docker-compose exec backend sh

# 执行迁移
./server migrate
```

2. 导入初始数据
```bash
# 进入 MySQL 容器
docker-compose exec mysql mysql -u root -p go_vue_admin

# 导入数据
source /path/to/init.sql
```

## 生产环境配置

### 安全配置

1. 使用环境变量
```yaml
# docker-compose.prod.yml
services:
  backend:
    environment:
      - GO_ENV=production
      - JWT_SECRET=${JWT_SECRET}
      - DB_PASSWORD=${DB_PASSWORD}
```

2. 配置 SSL
```nginx
# nginx.prod.conf
server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;

    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
}
```

### 性能优化

1. Nginx 优化
```nginx
# nginx.conf
worker_processes auto;
worker_rlimit_nofile 65535;

events {
    worker_connections 65535;
}

http {
    # 开启 gzip
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;
    
    # 缓存配置
    proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m;
}
```

2. MySQL 优化
```yaml
# docker-compose.prod.yml
services:
  mysql:
    command: 
      - --max_connections=1000
      - --innodb_buffer_pool_size=1G
```

### 监控配置

1. 添加 Prometheus
```yaml
# docker-compose.prod.yml
services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
```

2. 添加 Grafana
```yaml
services:
  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
```

## 维护指南

### 备份策略

1. 数据库备份
```bash
#!/bin/bash
# backup.sh
DATE=$(date +%Y%m%d_%H%M%S)
docker-compose exec -T mysql mysqldump -u root -p go_vue_admin > backup_$DATE.sql
```

2. 配置备份
```bash
# 备份配置文件
cp configs/config.ini configs/config.ini.backup
```

### 更新流程

1. 拉取更新
```bash
git pull origin main
```

2. 重建服务
```bash
# 重建并更新服务
docker-compose up -d --build

# 执行数据库迁移
docker-compose exec backend ./server migrate
```

### 日志管理

1. 配置日志驱动
```yaml
# docker-compose.prod.yml
services:
  backend:
    logging:
      driver: "json-file"
      options:
        max-size: "200m"
        max-file: "10"
```

2. 日志聚合
```yaml
services:
  filebeat:
    image: docker.elastic.co/beats/filebeat:7.9.3
    volumes:
      - ./filebeat.yml:/usr/share/filebeat/filebeat.yml
```

## 故障排除

### 常见问题

1. 容器无法启动
- 检查配置文件
- 查看容器日志
- 验证端口占用

2. 数据库连接失败
- 检查网络配置
- 验证密码正确
- 确认权限设置

3. 前端访问异常
- 检查 Nginx 配置
- 验证 API 地址
- 查看浏览器控制台

### 调试命令

```bash
# 查看容器状态
docker-compose ps

# 查看容器日志
docker-compose logs -f [service]

# 进入容器
docker-compose exec [service] sh

# 检查网络
docker network inspect go-vue-admin-network
```

## 参考资源

- [Docker 文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [Nginx 文档](https://nginx.org/en/docs/)
- [MySQL Docker 文档](https://hub.docker.com/_/mysql) 