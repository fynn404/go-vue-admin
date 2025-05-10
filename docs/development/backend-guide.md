# 后端开发指南

本文档详细说明 Go-Vue-Admin 项目的后端开发规范和指南。

## 项目结构

```
.
├── api/            # API 接口定义
│   └── v1/         # V1 版本接口
├── cmd/            # 主程序入口
├── configs/        # 配置文件
├── internal/       # 内部包
│   ├── handler/    # 请求处理器
│   ├── middleware/ # 中间件
│   ├── model/      # 数据模型
│   ├── repository/ # 数据访问层
│   └── service/    # 业务逻辑层
└── pkg/            # 公共包
    ├── auth/       # 认证相关
    ├── database/   # 数据库工具
    └── utils/      # 工具函数
```

## 架构设计

### 分层架构
1. Handler 层（表示层）
   - 处理 HTTP 请求
   - 参数验证
   - 响应封装

2. Service 层（业务层）
   - 业务逻辑处理
   - 事务管理
   - 数据组装

3. Repository 层（数据访问层）
   - 数据库操作
   - 缓存操作
   - 外部服务调用

4. Model 层（数据模型层）
   - 数据结构定义
   - 模型关联
   - 验证规则

## 开发规范

### 命名规范

1. 包名
   - 使用小写
   - 简短有意义
   - 避免下划线

2. 文件名
   - 使用小写
   - 使用下划线分隔
   - 描述文件用途

3. 函数名
   - 驼峰命名
   - 动词开头
   - 清晰表意

4. 变量名
   - 驼峰命名
   - 避免单字母
   - 表达含义

### 注释规范

1. 包注释
```go
// Package handler 处理 HTTP 请求和响应
package handler
```

2. 函数注释
```go
// CreateUser 创建新用户
// 参数：
//   - ctx: 上下文
//   - req: 创建用户请求
// 返回：
//   - 用户信息
//   - 错误信息
func CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
```

3. 变量注释
```go
// ErrUserNotFound 用户不存在错误
var ErrUserNotFound = errors.New("user not found")
```

## 编码指南

### 错误处理

1. 错误定义
```go
// pkg/errors/errors.go
var (
    ErrInvalidParam = errors.New("invalid parameter")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
)
```

2. 错误包装
```go
if err != nil {
    return fmt.Errorf("create user: %w", err)
}
```

3. 错误响应
```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

### 中间件

1. 认证中间件
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 验证 token
        // 设置用户信息
        c.Next()
    }
}
```

2. 日志中间件
```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 记录请求信息
        c.Next()
        // 记录响应信息
    }
}
```

### 数据库操作

1. 模型定义
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

2. 仓库模式
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uint) error
    FindByID(ctx context.Context, id uint) (*User, error)
}
```

3. 事务处理
```go
func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 事务操作
        return nil
    })
}
```

## API 开发流程

1. 定义请求/响应结构
```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Role     string `json:"role" binding:"required"`
}

type CreateUserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Role     string `json:"role"`
}
```

2. 实现 Handler
```go
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }
    
    user, err := h.service.CreateUser(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }
    
    c.JSON(http.StatusOK, SuccessResponse(user))
}
```

3. 注册路由
```go
func (h *Handler) RegisterRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")
    {
        users := v1.Group("/users")
        {
            users.POST("", h.CreateUser)
            users.GET("/:id", h.GetUser)
            users.PUT("/:id", h.UpdateUser)
            users.DELETE("/:id", h.DeleteUser)
        }
    }
}
```

## 测试指南

### 单元测试

1. Handler 测试
```go
func TestCreateUser(t *testing.T) {
    // 设置测试用例
    // 创建请求
    // 验证响应
}
```

2. Service 测试
```go
func TestUserService_Create(t *testing.T) {
    // 准备测试数据
    // 执行测试
    // 验证结果
}
```

3. Repository 测试
```go
func TestUserRepository_Create(t *testing.T) {
    // 准备测试数据库
    // 执行测试
    // 清理测试数据
}
```

### 集成测试

```go
func TestUserAPI(t *testing.T) {
    // 启动测试服务器
    // 执行 HTTP 请求
    // 验证响应
}
```

## 性能优化

1. 数据库优化
   - 使用适当的索引
   - 优化查询语句
   - 使用连接池

2. 缓存策略
   - 使用 Redis 缓存
   - 合理设置过期时间
   - 缓存预热

3. 并发处理
   - 使用 goroutine
   - 控制并发数量
   - 超时处理

## 部署相关

1. 编译
```bash
go build -o bin/server cmd/main.go
```

2. 配置文件
```ini
[server]
port = 8080
mode = release

[database]
host = localhost
port = 3306
```

3. 环境变量
```bash
export GO_ENV=production
export DB_PASSWORD=secret
```

## 最佳实践

1. 代码组织
   - 遵循标准布局
   - 合理分包
   - 避免循环依赖

2. 错误处理
   - 统一错误码
   - 错误日志记录
   - 优雅降级

3. 安全考虑
   - 参数验证
   - SQL 注入防护
   - XSS 防护

4. 日志记录
   - 分级日志
   - 结构化日志
   - 日志轮转

## 常见问题

1. 数据库连接
   - 连接池配置
   - 超时处理
   - 重试机制

2. 内存泄漏
   - goroutine 泄漏
   - 资源未释放
   - 大对象复制

3. 性能问题
   - CPU 占用高
   - 内存使用大
   - 响应时间长

## 参考资源

- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Go 编程规范](https://golang.org/doc/effective_go)
- [Go 标准库文档](https://golang.org/pkg/) 