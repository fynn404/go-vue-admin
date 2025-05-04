# 后端技术文档

## 架构设计

### 目录结构
```
.
├── api/            # API 接口定义
│   └── v1/         # V1 版本接口
├── cmd/            # 主程序入口
│   └── main.go     # 主程序
├── configs/        # 配置文件
│   ├── config.go   # 配置加载
│   └── config.ini  # 配置文件
├── internal/       # 内部包
│   ├── handler/    # 请求处理器
│   ├── middleware/ # 中间件
│   ├── model/      # 数据模型
│   ├── repository/ # 数据访问层
│   └── service/    # 业务逻辑层
├── pkg/            # 公共包
│   ├── auth/       # 认证相关
│   ├── database/   # 数据库工具
│   └── utils/      # 工具函数
└── resource/       # 资源文件
```

### 分层架构
- Handler Layer: 请求处理和参数验证
- Service Layer: 业务逻辑实现
- Repository Layer: 数据访问和持久化
- Model Layer: 数据模型定义

## 核心功能实现

### 用户管理
- 用户注册与登录
- JWT 认证
- 角色权限控制
- 用户信息管理

### 课程管理
- 课程 CRUD 操作
- 课程状态管理
- 课程容量控制
- 课程查询和筛选

### 选课系统
- 选课操作
- 退课操作
- 选课状态管理
- 选课记录查询

### 成绩管理
- 成绩录入与修改
- 成绩统计分析
  - 课程维度统计
  - 学生维度统计
- 成绩历史记录
- GPA 计算

## 数据模型

### User 模型
```go
type User struct {
    ID        uint      `gorm:"primarykey"`
    Username  string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    Role      string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Course 模型
```go
type Course struct {
    ID              uint      `gorm:"primarykey"`
    Name            string    `gorm:"not null"`
    Description     string
    Credits         float32   `gorm:"not null"`
    Capacity        int       `gorm:"not null"`
    CurrentEnrolled int       `gorm:"not null;default:0"`
    Status          string    `gorm:"not null;default:'open'"`
    TeacherID       uint      `gorm:"not null"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### Enrollment 模型
```go
type Enrollment struct {
    ID        uint      `gorm:"primarykey"`
    StudentID uint      `gorm:"not null"`
    CourseID  uint      `gorm:"not null"`
    Status    string    `gorm:"not null;default:'active'"`
    Grade     float32
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Grade 模型
```go
type Grade struct {
    ID         uint      `gorm:"primarykey"`
    StudentID  uint      `gorm:"not null"`
    CourseID   uint      `gorm:"not null"`
    Score      float32   `gorm:"not null"`
    Comment    string
    CreatedAt  time.Time
    UpdatedAt  time.Time
    ModifiedBy uint      `gorm:"not null"`
}
```

### GradeHistory 模型
```go
type GradeHistory struct {
    ID         uint      `gorm:"primarykey"`
    GradeID    uint      `gorm:"not null"`
    OldScore   float32   `gorm:"not null"`
    NewScore   float32   `gorm:"not null"`
    Comment    string
    ModifiedBy uint      `gorm:"not null"`
    CreatedAt  time.Time
}
```

## API 接口

### 用户相关
- POST /api/v1/auth/register - 用户注册
- POST /api/v1/auth/login - 用户登录
- GET /api/v1/users/me - 获取当前用户信息
- PUT /api/v1/users/me - 更新用户信息

### 课程相关
- GET /api/v1/courses - 获取课程列表
- POST /api/v1/courses - 创建课程
- PUT /api/v1/courses/:id - 更新课程
- DELETE /api/v1/courses/:id - 删除课程
- PUT /api/v1/courses/:id/status - 更新课程状态

### 选课相关
- POST /api/v1/courses/:id/enroll - 选课
- POST /api/v1/courses/:id/drop - 退课
- GET /api/v1/courses/enrolled - 获取已选课程

### 成绩相关
- GET /api/v1/grades - 获取成绩列表
- POST /api/v1/grades - 录入成绩
- PUT /api/v1/grades/:id - 修改成绩
- GET /api/v1/grades/:id/history - 获取成绩历史
- GET /api/v1/grades/stats/course/:id - 获取课程成绩统计
- GET /api/v1/grades/stats/student/:id - 获取学生成绩统计

## 中间件

### 认证中间件
- JWT 令牌验证
- 用户身份解析
- 权限检查

### 日志中间件
- 请求日志记录
- 错误日志记录

### 错误处理中间件
- 统一错误响应
- 错误码管理

## 工具函数

### 认证相关
- JWT 生成与验证
- 密码加密与验证

### 数据库工具
- 数据库连接管理
- 事务处理

### 通用工具
- 分页处理
- 响应封装
- 参数验证 