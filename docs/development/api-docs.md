# API 文档

本文档详细说明 Go-Vue-Admin 项目的 API 接口规范和使用说明。

## API 规范

### 基础信息
- 基础路径：`/api/v1`
- 请求方式：REST
- 数据格式：JSON
- 字符编码：UTF-8

### 认证方式
- 使用 JWT Token
- Token 在 Header 中通过 `Authorization: Bearer <token>` 传递

### 响应格式

```json
{
    "code": 0,        // 状态码，0 表示成功
    "message": "",    // 提示信息
    "data": null      // 响应数据
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40001 | 参数错误 |
| 40100 | 未授权 |
| 40300 | 禁止访问 |
| 40400 | 资源不存在 |
| 50000 | 服务器错误 |

## 认证接口

### 登录

```http
POST /api/v1/auth/login
```

请求参数：
```json
{
    "username": "admin",
    "password": "password"
}
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "user": {
            "id": 1,
            "username": "admin",
            "role": "admin"
        }
    }
}
```

### 登出

```http
POST /api/v1/auth/logout
```

响应示例：
```json
{
    "code": 0,
    "message": "success"
}
```

### 获取当前用户信息

```http
GET /api/v1/users/me
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "username": "admin",
        "role": "admin",
        "email": "admin@example.com",
        "created_at": "2024-03-21T10:00:00Z"
    }
}
```

## 用户管理

### 创建用户

```http
POST /api/v1/users
```

请求参数：
```json
{
    "username": "teacher1",
    "password": "password",
    "role": "teacher",
    "email": "teacher1@example.com"
}
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 2,
        "username": "teacher1",
        "role": "teacher",
        "email": "teacher1@example.com",
        "created_at": "2024-03-21T10:00:00Z"
    }
}
```

### 更新用户

```http
PUT /api/v1/users/:id
```

请求参数：
```json
{
    "email": "new.email@example.com",
    "password": "new_password"
}
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 2,
        "username": "teacher1",
        "role": "teacher",
        "email": "new.email@example.com",
        "updated_at": "2024-03-21T11:00:00Z"
    }
}
```

### 删除用户

```http
DELETE /api/v1/users/:id
```

响应示例：
```json
{
    "code": 0,
    "message": "success"
}
```

### 获取用户列表

```http
GET /api/v1/users?page=1&size=10&role=teacher
```

查询参数：
- page: 页码，默认 1
- size: 每页数量，默认 10
- role: 用户角色，可选

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "total": 100,
        "items": [
            {
                "id": 1,
                "username": "teacher1",
                "role": "teacher",
                "email": "teacher1@example.com",
                "created_at": "2024-03-21T10:00:00Z"
            }
        ]
    }
}
```

## 课程管理

### 创建课程

```http
POST /api/v1/courses
```

请求参数：
```json
{
    "name": "数据结构",
    "description": "课程描述",
    "credits": 3,
    "capacity": 100,
    "teacher_id": 2
}
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "name": "数据结构",
        "description": "课程描述",
        "credits": 3,
        "capacity": 100,
        "current_enrolled": 0,
        "status": "open",
        "teacher_id": 2,
        "created_at": "2024-03-21T10:00:00Z"
    }
}
```

### 更新课程

```http
PUT /api/v1/courses/:id
```

请求参数：
```json
{
    "name": "高级数据结构",
    "description": "新的课程描述",
    "capacity": 150
}
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "name": "高级数据结构",
        "description": "新的课程描述",
        "credits": 3,
        "capacity": 150,
        "current_enrolled": 0,
        "status": "open",
        "teacher_id": 2,
        "updated_at": "2024-03-21T11:00:00Z"
    }
}
```

### 获取课程列表

```http
GET /api/v1/courses?page=1&size=10&status=open
```

查询参数：
- page: 页码，默认 1
- size: 每页数量，默认 10
- status: 课程状态，可选
- teacher_id: 教师 ID，可选

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "total": 50,
        "items": [
            {
                "id": 1,
                "name": "高级数据结构",
                "description": "课程描述",
                "credits": 3,
                "capacity": 150,
                "current_enrolled": 0,
                "status": "open",
                "teacher_id": 2,
                "created_at": "2024-03-21T10:00:00Z"
            }
        ]
    }
}
```

## 选课管理

### 学生选课

```http
POST /api/v1/enrollments
```

请求参数：
```json
{
    "course_id": 1
}
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "student_id": 3,
        "course_id": 1,
        "status": "active",
        "created_at": "2024-03-21T10:00:00Z"
    }
}
```

### 退课

```http
DELETE /api/v1/enrollments/:id
```

响应示例：
```json
{
    "code": 0,
    "message": "success"
}
```

### 获取学生选课列表

```http
GET /api/v1/enrollments?page=1&size=10&status=active
```

查询参数：
- page: 页码，默认 1
- size: 每页数量，默认 10
- status: 选课状态，可选

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "total": 20,
        "items": [
            {
                "id": 1,
                "student_id": 3,
                "course_id": 1,
                "course_name": "高级数据结构",
                "status": "active",
                "created_at": "2024-03-21T10:00:00Z"
            }
        ]
    }
}
```

## 成绩管理

### 录入成绩

```http
POST /api/v1/grades
```

请求参数：
```json
{
    "enrollment_id": 1,
    "score": 85,
    "comment": "表现良好"
}
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "enrollment_id": 1,
        "student_id": 3,
        "course_id": 1,
        "score": 85,
        "comment": "表现良好",
        "created_at": "2024-03-21T10:00:00Z"
    }
}
```

### 更新成绩

```http
PUT /api/v1/grades/:id
```

请求参数：
```json
{
    "score": 88,
    "comment": "期末考试表现优秀"
}
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "enrollment_id": 1,
        "student_id": 3,
        "course_id": 1,
        "score": 88,
        "comment": "期末考试表现优秀",
        "updated_at": "2024-03-21T11:00:00Z"
    }
}
```

### 获取成绩列表

```http
GET /api/v1/grades?page=1&size=10&course_id=1
```

查询参数：
- page: 页码，默认 1
- size: 每页数量，默认 10
- course_id: 课程 ID，可选
- student_id: 学生 ID，可选

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "total": 100,
        "items": [
            {
                "id": 1,
                "enrollment_id": 1,
                "student_id": 3,
                "student_name": "张三",
                "course_id": 1,
                "course_name": "高级数据结构",
                "score": 88,
                "comment": "期末考试表现优秀",
                "created_at": "2024-03-21T10:00:00Z",
                "updated_at": "2024-03-21T11:00:00Z"
            }
        ]
    }
}
```

### 获取成绩统计

```http
GET /api/v1/grades/stats/course/:id
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "course_id": 1,
        "course_name": "高级数据结构",
        "total_students": 100,
        "average_score": 85.5,
        "max_score": 98,
        "min_score": 60,
        "pass_rate": 0.95,
        "score_distribution": {
            "90-100": 20,
            "80-89": 40,
            "70-79": 25,
            "60-69": 10,
            "0-59": 5
        }
    }
}
```

## 错误处理

### 参数错误

```json
{
    "code": 40001,
    "message": "invalid parameter: username is required",
    "data": null
}
```

### 未授权

```json
{
    "code": 40100,
    "message": "unauthorized",
    "data": null
}
```

### 禁止访问

```json
{
    "code": 40300,
    "message": "forbidden",
    "data": null
}
```

### 资源不存在

```json
{
    "code": 40400,
    "message": "resource not found",
    "data": null
}
```

### 服务器错误

```json
{
    "code": 50000,
    "message": "internal server error",
    "data": null
}
```

## API 调用示例

### cURL 示例

1. 登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

2. 创建课程
```bash
curl -X POST http://localhost:8080/api/v1/courses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "数据结构",
    "description": "课程描述",
    "credits": 3,
    "capacity": 100,
    "teacher_id": 2
  }'
```

### JavaScript 示例

```javascript
// 登录
async function login(username, password) {
    const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
    });
    return await response.json();
}

// 获取课程列表
async function getCourses(page = 1, size = 10) {
    const response = await fetch(`/api/v1/courses?page=${page}&size=${size}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return await response.json();
}
```

## 注意事项

1. 安全性
   - 所有请求都需要进行认证（除了登录接口）
   - 敏感数据传输使用 HTTPS
   - 注意 CORS 配置

2. 性能
   - 使用分页查询
   - 合理使用缓存
   - 避免大量数据传输

3. 错误处理
   - 统一错误响应格式
   - 合适的错误码
   - 详细的错误信息

4. 版本控制
   - API 版本在 URL 中体现
   - 向后兼容性
   - 版本升级策略 