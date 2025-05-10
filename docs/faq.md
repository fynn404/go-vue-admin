# 常见问题解答 (FAQ)

本文档收集了 Go-Vue-Admin 项目的常见问题和解答。

## 目录
- [开发环境问题](#开发环境问题)
- [部署问题](#部署问题)
- [功能问题](#功能问题)
- [性能问题](#性能问题)
- [安全问题](#安全问题)

## 开发环境问题

### Q: 如何快速搭建开发环境？
A: 按照以下步骤操作：
1. 安装必要软件：
   ```bash
   # macOS
   brew install go node mysql docker
   
   # Ubuntu
   sudo apt-get install golang nodejs mysql-server docker.io
   ```

2. 克隆项目：
   ```bash
   git clone https://github.com/yourusername/go-vue-admin.git
   cd go-vue-admin
   ```

3. 安装依赖：
   ```bash
   # 后端依赖
   go mod download
   
   # 前端依赖
   cd frontend
   npm install
   ```

### Q: 为什么我的 Go 模块无法下载？
A: 可能的原因和解决方案：
1. 检查 Go 环境变量：
   ```bash
   go env
   ```

2. 设置 GOPROXY：
   ```bash
   go env -w GOPROXY=https://goproxy.cn,direct
   ```

3. 清理模块缓存：
   ```bash
   go clean -modcache
   ```

### Q: 前端开发服务器启动失败？
A: 常见原因：
1. 端口被占用：
   ```bash
   # 查看端口占用
   lsof -i :3000
   
   # 终止进程
   kill -9 <PID>
   ```

2. Node.js 版本不兼容：
   ```bash
   # 使用 nvm 安装正确版本
   nvm install 14
   nvm use 14
   ```

## 部署问题

### Q: Docker 部署时容器无法启动？
A: 检查以下几点：
1. 查看容器日志：
   ```bash
   docker-compose logs
   ```

2. 检查配置文件：
   ```bash
   # 确保配置文件存在
   ls configs/config.ini
   
   # 检查环境变量
   cat .env
   ```

3. 检查网络配置：
   ```bash
   docker network ls
   docker network inspect go-vue-admin-network
   ```

### Q: 数据库连接失败？
A: 可能的解决方案：
1. 检查数据库配置：
   ```ini
   [database]
   host = localhost
   port = 3306
   username = root
   password = your_password
   ```

2. 确认数据库服务状态：
   ```bash
   # MySQL
   sudo service mysql status
   
   # Docker
   docker-compose ps mysql
   ```

3. 检查网络连接：
   ```bash
   telnet localhost 3306
   ```

### Q: Nginx 配置后无法访问前端？
A: 检查以下配置：
1. Nginx 配置文件：
   ```nginx
   server {
       listen 80;
       root /usr/share/nginx/html;
       
       location / {
           try_files $uri $uri/ /index.html;
       }
       
       location /api {
           proxy_pass http://backend:8080;
       }
   }
   ```

2. 静态文件路径：
   ```bash
   ls -l /usr/share/nginx/html
   ```

3. 检查 Nginx 日志：
   ```bash
   tail -f /var/log/nginx/error.log
   ```

## 功能问题

### Q: 用户登录失败？
A: 常见原因：
1. 检查用户凭证：
   - 确认用户名和密码正确
   - 检查大小写敏感性
   - 确认账户状态

2. 检查 JWT 配置：
   ```ini
   [jwt]
   secret = your_secret_key
   expire = 24h
   ```

3. 检查请求头：
   ```http
   Authorization: Bearer <token>
   ```

### Q: 选课功能异常？
A: 可能的问题：
1. 课程容量限制：
   - 检查课程当前人数
   - 验证容量设置
   - 确认选课时间

2. 权限问题：
   ```go
   // 检查用户角色
   if user.Role != "student" {
       return errors.New("only students can select courses")
   }
   ```

3. 数据一致性：
   ```sql
   -- 检查选课记录
   SELECT * FROM enrollments WHERE student_id = ? AND course_id = ?;
   ```

### Q: 成绩统计不准确？
A: 排查步骤：
1. 检查计算逻辑：
   ```go
   func calculateAverage(scores []float64) float64 {
       if len(scores) == 0 {
           return 0
       }
       sum := 0.0
       for _, score := range scores {
           sum += score
       }
       return sum / float64(len(scores))
   }
   ```

2. 验证数据完整性：
   ```sql
   SELECT COUNT(*), AVG(score) FROM grades WHERE course_id = ?;
   ```

3. 检查前端展示：
   ```javascript
   // 确保数据格式化正确
   const formatScore = (score) => Number(score).toFixed(1);
   ```

## 性能问题

### Q: API 响应速度慢？
A: 优化建议：
1. 数据库优化：
   ```sql
   -- 添加索引
   CREATE INDEX idx_user_id ON users(id);
   CREATE INDEX idx_course_id ON courses(id);
   ```

2. 缓存配置：
   ```go
   // 使用 Redis 缓存
   func (s *service) GetUserByID(ctx context.Context, id uint) (*User, error) {
       key := fmt.Sprintf("user:%d", id)
       if cached, err := s.redis.Get(ctx, key).Result(); err == nil {
           // 返回缓存数据
       }
       // 查询数据库
   }
   ```

3. 查询优化：
   ```go
   // 使用预加载
   db.Preload("Courses").Find(&users)
   ```

### Q: 前端加载缓慢？
A: 优化方案：
1. 代码分割：
   ```javascript
   // 路由懒加载
   const UserList = () => import('@/views/UserList.vue')
   ```

2. 资源优化：
   ```nginx
   # 开启 gzip
   gzip on;
   gzip_types text/plain application/javascript text/css;
   ```

3. 缓存策略：
   ```nginx
   location /static/ {
       expires 7d;
       add_header Cache-Control "public, no-transform";
   }
   ```

## 安全问题

### Q: 如何防止 SQL 注入？
A: 安全措施：
1. 使用参数化查询：
   ```go
   db.Where("username = ?", username).First(&user)
   ```

2. 验证输入：
   ```go
   if !validateInput(input) {
       return errors.New("invalid input")
   }
   ```

3. 使用 ORM：
   ```go
   // GORM 默认防止 SQL 注入
   db.Create(&user)
   ```

### Q: 如何保护敏感数据？
A: 建议措施：
1. 密码加密：
   ```go
   // 使用 bcrypt
   hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
   ```

2. 数据脱敏：
   ```go
   type User struct {
       ID       uint   `json:"id"`
       Password string `json:"-"`  // 不返回密码
   }
   ```

3. HTTPS 配置：
   ```nginx
   server {
       listen 443 ssl;
       ssl_certificate /etc/nginx/ssl/cert.pem;
       ssl_certificate_key /etc/nginx/ssl/key.pem;
   }
   ```

### Q: 如何处理 CORS？
A: 配置示例：
1. 后端配置：
   ```go
   // Gin 中间件
   router.Use(cors.New(cors.Config{
       AllowOrigins:     []string{"https://example.com"},
       AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
       AllowHeaders:     []string{"Origin", "Authorization"},
       ExposeHeaders:    []string{"Content-Length"},
       AllowCredentials: true,
   }))
   ```

2. Nginx 配置：
   ```nginx
   location /api {
       add_header 'Access-Control-Allow-Origin' 'https://example.com';
       add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE';
   }
   ```

3. 前端配置：
   ```javascript
   axios.defaults.withCredentials = true;
   ```

## 其他问题

### Q: 如何贡献代码？
A: 参考流程：
1. Fork 项目
2. 创建功能分支
3. 提交变更
4. 发起 Pull Request

详见 [贡献指南](./contributing.md)

### Q: 如何报告 Bug？
A: 建议步骤：
1. 检查是否是已知问题
2. 收集错误信息
3. 创建详细的 Issue
4. 提供复现步骤

### Q: 在哪里获取帮助？
A: 支持渠道：
1. GitHub Issues
2. 技术文档
3. 社区讨论
4. 邮件支持 