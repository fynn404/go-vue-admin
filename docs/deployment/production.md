# 生产环境配置指南

本文档详细说明如何配置和优化 Go-Vue-Admin 项目的生产环境。

## 服务器要求

### 硬件推荐配置
- CPU: 4核+
- 内存: 8GB+
- 磁盘: 50GB+ SSD
- 网络: 100Mbps+

### 软件要求
- 操作系统: Ubuntu 20.04 LTS
- Docker: 20.10+
- Docker Compose: 2.0+
- Nginx: 1.18+
- MySQL: 5.7+

## 系统配置

### 操作系统优化

1. 系统参数优化
```bash
# /etc/sysctl.conf
# 网络优化
net.ipv4.tcp_max_syn_backlog = 8192
net.ipv4.tcp_max_tw_buckets = 5000
net.ipv4.tcp_max_orphans = 3276800
net.ipv4.tcp_syncookies = 1
net.core.somaxconn = 65535

# 文件系统优化
fs.file-max = 655350
fs.inotify.max_user_watches = 524288

# 应用配置
vm.swappiness = 10
vm.dirty_ratio = 60
vm.dirty_background_ratio = 2
```

2. 系统限制优化
```bash
# /etc/security/limits.conf
*         soft    nofile      65535
*         hard    nofile      65535
*         soft    nproc       65535
*         hard    nproc       65535
```

### 安全配置

1. 防火墙配置
```bash
# 开启防火墙
ufw enable

# 允许 SSH
ufw allow 22

# 允许 HTTP/HTTPS
ufw allow 80
ufw allow 443

# 允许特定 IP
ufw allow from 192.168.1.0/24 to any port 3306
```

2. SSH 配置
```bash
# /etc/ssh/sshd_config
Port 22
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
```

## 应用配置

### 后端配置

1. 环境变量
```bash
# /etc/environment
export GO_ENV=production
export JWT_SECRET=your-secure-jwt-secret
export DB_PASSWORD=your-secure-password
```

2. 应用配置
```ini
# configs/config.prod.ini
[server]
port = 8080
mode = release
timeout = 30

[database]
host = mysql
port = 3306
username = root
password = ${DB_PASSWORD}
dbname = go_vue_admin
max_idle_conns = 10
max_open_conns = 100
conn_max_lifetime = 3600

[jwt]
secret = ${JWT_SECRET}
expire = 24h

[log]
level = info
file = /var/log/go-vue-admin/app.log
max_size = 100
max_backups = 10
max_age = 30
compress = true
```

### 前端配置

1. 环境变量
```bash
# frontend/.env.production
VITE_API_BASE_URL=/api/v1
VITE_APP_TITLE=Go-Vue-Admin
```

2. 构建配置
```javascript
// vite.config.ts
export default defineConfig({
  build: {
    target: 'es2015',
    minify: 'terser',
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'element-plus': ['element-plus'],
          'echarts': ['echarts']
        }
      }
    }
  }
})
```

## 数据库配置

### MySQL 优化

1. 基础配置
```ini
# /etc/mysql/mysql.conf.d/mysqld.cnf
[mysqld]
# 基础配置
max_connections = 1000
max_allowed_packet = 64M
table_open_cache = 4000
table_definition_cache = 4096

# InnoDB 配置
innodb_buffer_pool_size = 4G
innodb_log_file_size = 512M
innodb_flush_log_at_trx_commit = 2
innodb_flush_method = O_DIRECT
innodb_file_per_table = 1

# 查询缓存
query_cache_type = 1
query_cache_size = 128M
query_cache_limit = 2M
```

2. 备份配置
```bash
#!/bin/bash
# /usr/local/bin/backup-mysql.sh
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/mysql"
mysqldump -u root -p go_vue_admin > $BACKUP_DIR/backup_$DATE.sql
```

### Redis 配置（可选）

```ini
# /etc/redis/redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
appendonly yes
appendfsync everysec
```

## 日志管理

### 日志配置

1. 应用日志
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

2. Nginx 日志
```nginx
# nginx.conf
access_log /var/log/nginx/access.log combined buffer=512k flush=1m;
error_log /var/log/nginx/error.log warn;
```

### 日志轮转

```conf
# /etc/logrotate.d/go-vue-admin
/var/log/go-vue-admin/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 root root
    postrotate
        systemctl reload go-vue-admin
    endscript
}
```

## 监控配置

### Prometheus 配置

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'go-vue-admin'
    static_configs:
      - targets: ['backend:8080']
```

### Grafana 配置

1. 数据源配置
- 添加 Prometheus 数据源
- 配置基础监控面板

2. 告警配置
```yaml
# alertmanager.yml
route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h
  receiver: 'web.hook'
```

## 性能优化

### Nginx 优化

```nginx
# nginx.conf
worker_processes auto;
worker_rlimit_nofile 65535;

events {
    worker_connections 65535;
    use epoll;
    multi_accept on;
}

http {
    # 开启 gzip
    gzip on;
    gzip_comp_level 5;
    gzip_min_length 256;
    gzip_types
        application/javascript
        application/json
        application/x-javascript
        text/css
        text/plain;

    # 缓存配置
    proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m;
    
    # 客户端优化
    client_max_body_size 10M;
    client_body_buffer_size 128k;
    
    # 超时设置
    keepalive_timeout 65;
    keepalive_requests 100;
}
```

### 缓存策略

1. 浏览器缓存
```nginx
# nginx.conf
location /static/ {
    expires 7d;
    add_header Cache-Control "public, no-transform";
}

location /api/ {
    add_header Cache-Control "no-cache, no-store";
}
```

2. 服务端缓存
```go
// 使用 Redis 缓存
func (s *service) GetUserByID(ctx context.Context, id uint) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // 尝试从缓存获取
    if cached, err := s.redis.Get(ctx, key).Result(); err == nil {
        var user User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    
    // 从数据库获取
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 设置缓存
    if bytes, err := json.Marshal(user); err == nil {
        s.redis.Set(ctx, key, bytes, time.Hour)
    }
    
    return user, nil
}
```

## 安全加固

### SSL 配置

1. 获取证书
```bash
# 使用 Let's Encrypt
certbot --nginx -d your-domain.com
```

2. SSL 配置
```nginx
# nginx.conf
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
ssl_prefer_server_ciphers off;
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;
```

### 安全头配置

```nginx
# nginx.conf
add_header X-Frame-Options "SAMEORIGIN";
add_header X-XSS-Protection "1; mode=block";
add_header X-Content-Type-Options "nosniff";
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

## 部署检查清单

### 上线前检查

1. 安全检查
- [ ] 所有密码和密钥都使用环境变量
- [ ] SSL 证书已配置
- [ ] 防火墙规则已设置
- [ ] 敏感信息已加密

2. 性能检查
- [ ] 数据库索引已优化
- [ ] 缓存策略已实施
- [ ] 静态资源已压缩
- [ ] 负载测试已完成

3. 监控检查
- [ ] 监控系统已部署
- [ ] 告警规则已配置
- [ ] 日志收集已配置
- [ ] 备份策略已实施

### 定期维护

1. 更新计划
- 系统更新
- 依赖更新
- 安全补丁

2. 备份验证
- 数据库备份
- 配置备份
- 恢复测试

3. 性能监控
- 系统资源
- 应用性能
- 数据库性能

## 参考资源

- [Ubuntu 服务器指南](https://ubuntu.com/server/docs)
- [Nginx 文档](https://nginx.org/en/docs/)
- [MySQL 优化指南](https://dev.mysql.com/doc/refman/5.7/en/optimization.html)
- [Docker 生产环境最佳实践](https://docs.docker.com/develop/dev-best-practices/) 