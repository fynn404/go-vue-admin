# 后端功能实现文档

## 项目结构

```
backend/
├── config/         # 配置文件
├── controllers/    # 控制器
├── middleware/     # 中间件
├── models/         # 数据模型
├── routes/         # 路由定义
├── services/       # 业务逻辑
├── utils/          # 工具函数
└── main.go         # 入口文件
```

## 技术栈

- Go 1.20+
- Gin Web 框架
- GORM ORM 框架
- MySQL 数据库
- Redis 缓存
- JWT 认证
- Air 热重载

## 数据库设计

### 1. 用户表 (users)
```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    role ENUM('admin', 'teacher', 'student') NOT NULL,
    status TINYINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### 2. 课程表 (courses)
```sql
CREATE TABLE courses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    teacher_id BIGINT UNSIGNED NOT NULL,
    credits DECIMAL(3,1) NOT NULL,
    capacity INT NOT NULL,
    current_enrolled INT NOT NULL DEFAULT 0,
    status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### 3. 选课记录表 (enrollments)
```sql
CREATE TABLE enrollments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    student_id BIGINT UNSIGNED NOT NULL,
    course_id BIGINT UNSIGNED NOT NULL,
    status ENUM('active', 'dropped', 'completed') NOT NULL DEFAULT 'active',
    grade DECIMAL(5,2) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

## API 接口设计

### 1. 认证接口

#### 1.1 用户登录
- 路径: POST /api/v1/auth/login
- 参数:
  - username: string
  - password: string
- 返回:
  - token: string
  - user: object

#### 1.2 用户注册
- 路径: POST /api/v1/auth/register
- 参数:
  - username: string
  - password: string
  - email: string
  - role: string
- 返回:
  - message: string
  - user: object

### 2. 用户接口

#### 2.1 获取用户信息
- 路径: GET /api/v1/users/:id
- 权限: 认证用户
- 返回: user object

#### 2.2 更新用户信息
- 路径: PUT /api/v1/users/:id
- 权限: 用户本人或管理员
- 参数: 用户信息对象
- 返回: 更新后的用户信息

### 3. 课程接口

#### 3.1 课程列表
- 路径: GET /api/v1/courses
- 参数:
  - page: int
  - size: int
  - search: string
  - status: string
- 返回: 课程列表和分页信息

#### 3.2 课程详情
- 路径: GET /api/v1/courses/:id
- 返回: 课程详细信息

#### 3.3 创建课程
- 路径: POST /api/v1/courses
- 权限: 教师或管理员
- 参数: 课程信息对象
- 返回: 创建的课程信息

### 4. 选课接口

#### 4.1 选课
- 路径: POST /api/v1/enrollments
- 权限: 学生
- 参数:
  - course_id: int
- 返回: 选课结果

#### 4.2 退课
- 路径: DELETE /api/v1/enrollments/:id
- 权限: 学生
- 返回: 退课结果

## 中间件实现

### 1. 认证中间件
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // JWT 认证逻辑
        // 用户信息验证
        // 权限检查
    }
}
```

### 2. 角色中间件
```go
func RoleMiddleware(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 角色验证逻辑
        // 权限控制
    }
}
```

### 3. 日志中间件
```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 请求日志记录
        // 响应时间统计
        // 错误日志记录
    }
}
```

## 业务逻辑实现

### 1. 用户服务
```go
type UserService interface {
    Create(user *models.User) error
    Update(user *models.User) error
    Delete(id uint) error
    FindByID(id uint) (*models.User, error)
    FindByUsername(username string) (*models.User, error)
}
```

### 2. 课程服务
```go
type CourseService interface {
    Create(course *models.Course) error
    Update(course *models.Course) error
    Delete(id uint) error
    FindByID(id uint) (*models.Course, error)
    List(params *ListParams) ([]*models.Course, int64, error)
}
```

### 3. 选课服务
```go
type EnrollmentService interface {
    Enroll(studentID, courseID uint) error
    Drop(enrollmentID uint) error
    UpdateGrade(enrollmentID uint, grade float64) error
    ListByStudent(studentID uint) ([]*models.Enrollment, error)
    ListByCourse(courseID uint) ([]*models.Enrollment, error)
}
```

## 缓存策略

### 1. Redis 缓存
- 热门课程缓存
- 用户信息缓存
- 选课记录缓存
- 课程统计信息缓存

### 2. 缓存更新策略
- 定时更新
- 写入时更新
- 过期时间设置
- 缓存预热

## 安全措施

### 1. 密码安全
- bcrypt 加密
- 密码强度验证
- 登录失败限制
- 密码重置机制

### 2. 数据安全
- SQL 注入防护
- XSS 防护
- CSRF 防护
- 输入验证

### 3. 访问控制
- JWT 认证
- 角色权限
- 资源访问控制
- API 访问限制

## 部署配置

### 1. 环境配置
```ini
[server]
port = 8080
mode = release

[database]
host = localhost
port = 3306
name = course_admin
user = root
password = root

[redis]
host = localhost
port = 6379
password = 
db = 0

[jwt]
secret = your_secret_key
expire_hours = 24
```

### 2. 监控告警
- 系统监控
- 性能监控
- 错误监控
- 告警通知

## 开发计划

### 第一阶段：基础架构
- [x] 项目初始化
- [x] 数据库设计
- [x] 基础框架搭建
- [x] 中间件实现

### 第二阶段：核心功能
- [ ] 用户认证
- [ ] 课程管理
- [ ] 选课系统
- [ ] 成绩管理

### 第三阶段：功能优化
- [ ] 缓存优化
- [ ] 性能优化
- [ ] 安全加固
- [ ] 单元测试

### 第四阶段：部署上线
- [ ] 部署脚本
- [ ] 监控配置
- [ ] 日志系统
- [ ] 性能测试 