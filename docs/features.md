# 成绩管理系统功能文档

## 1. 系统概述

本系统是一个完整的教育机构成绩管理系统，支持教师录入成绩、学生查看成绩，并具有完整的成绩历史记录追踪功能。系统采用基于角色的访问控制（RBAC），支持管理员、教师和学生三种角色。

## 2. 用户角色

### 2.1 管理员（Admin）
- 系统管理和监控
- 查看所有用户信息
- 查看所有成绩记录
- 查看系统操作历史

### 2.2 教师（Teacher）
- 创建和管理课程
- 录入和修改成绩
- 发布成绩
- 查看课程统计信息
- 查看成绩修改历史

### 2.3 学生（Student）
- 查看个人成绩
- 查看成绩历史记录
- 选课和退课
- 查看个人课程表

## 3. 核心功能

### 3.1 用户认证
- 用户注册
- 用户登录
- 个人信息管理

### 3.2 课程管理
- 课程创建和修改
- 课程列表查询
- 课程详情查看
- 课程统计信息

### 3.3 成绩管理
- 成绩录入
- 成绩修改
- 成绩发布
- 成绩查询
- 成绩历史记录

### 3.4 选课管理
- 学生选课
- 学生退课
- 选课记录查询
- 课程成绩统计

## 4. API 接口

### 4.1 认证接口
```
POST /api/v1/auth/login          # 用户登录
POST /api/v1/auth/register      # 用户注册
```

### 4.2 用户接口
```
GET  /api/v1/users/profile      # 获取个人信息
PUT  /api/v1/users/profile      # 更新个人信息
```

### 4.3 课程接口
```
# 公共接口
GET  /api/v1/courses           # 获取课程列表
GET  /api/v1/courses/:id       # 获取课程详情

# 教师接口
POST   /api/v1/courses         # 创建课程
PUT    /api/v1/courses/:id     # 更新课程
DELETE /api/v1/courses/:id     # 删除课程

# 学生接口
POST   /api/v1/courses/:id/enroll  # 选课
POST   /api/v1/courses/:id/drop    # 退课
```

### 4.4 成绩接口
```
# 学生接口
GET  /api/v1/grades/my         # 获取我的成绩列表

# 教师接口
POST /api/v1/grades/courses/:course_id/students/:student_id  # 创建成绩
GET  /api/v1/grades/courses/:course_id                      # 获取课程成绩列表
PUT  /api/v1/grades/:id                                     # 更新成绩
POST /api/v1/grades/:id/publish                             # 发布成绩

# 通用接口
GET  /api/v1/grades/:id/history                             # 获取成绩历史记录
```

### 4.5 成绩历史记录接口
```
# 学生接口
GET /api/v1/grade-histories/my                    # 获取我的成绩历史记录

# 教师接口
GET /api/v1/grade-histories/courses/:course_id    # 获取课程成绩历史记录
GET /api/v1/grade-histories/grades/:grade_id      # 获取单个成绩的历史记录
GET /api/v1/grade-histories/my-operations         # 获取教师的操作历史

# 管理员接口
GET /api/v1/grade-histories/date-range            # 按日期范围查询历史记录
    参数：
    - start_date: 开始日期 (YYYY-MM-DD)
    - end_date: 结束日期 (YYYY-MM-DD)

# 通用接口
GET /api/v1/grade-histories/:id                   # 获取历史记录详情
```

### 4.6 选课记录接口
```
GET  /api/v1/enrollments                          # 获取选课记录列表

# 教师接口
PUT  /api/v1/enrollments/:id/grade                # 更新选课成绩
POST /api/v1/enrollments/grades/batch             # 批量更新成绩
GET  /api/v1/enrollments/courses/:id/stats        # 获取课程统计信息

# 学生接口
GET  /api/v1/enrollments/grades                   # 获取学生成绩列表
```

## 5. 数据模型

### 5.1 用户（User）
- ID：用户唯一标识
- 用户名：登录账号
- 密码：加密存储
- 角色：admin/teacher/student
- 姓名：真实姓名
- 其他基本信息

### 5.2 课程（Course）
- ID：课程唯一标识
- 名称：课程名称
- 教师ID：授课教师
- 学分：课程学分
- 描述：课程描述
- 状态：课程状态

### 5.3 成绩（Grade）
- ID：成绩唯一标识
- 课程ID：关联课程
- 学生ID：关联学生
- 教师ID：操作教师
- 分数：课程得分
- 评语：教师评语
- 状态：成绩状态（草稿/已发布）
- 发布时间：成绩发布时间

### 5.4 成绩历史（GradeHistory）
- ID：历史记录唯一标识
- 成绩ID：关联成绩
- 课程ID：关联课程
- 学生ID：关联学生
- 教师ID：操作教师
- 变更类型：create/update/delete
- 原分数：变更前分数
- 新分数：变更后分数
- 原评语：变更前评语
- 新评语：变更后评语
- 修改原因：变更原因
- 操作时间：记录创建时间

### 5.5 选课记录（Enrollment）
- ID：选课记录唯一标识
- 课程ID：关联课程
- 学生ID：关联学生
- 状态：选课状态
- 选课时间：记录创建时间

## 6. 安全特性

### 6.1 认证机制
- JWT Token 认证
- Token 过期机制
- 刷新 Token 机制

### 6.2 授权机制
- 基于角色的访问控制（RBAC）
- 细粒度的权限控制
- API 级别的权限校验

### 6.3 数据安全
- 密码加密存储
- 敏感数据加密
- SQL 注入防护
- XSS 防护

## 7. 其他特性

### 7.1 审计日志
- 操作日志记录
- 成绩变更历史
- 用户行为追踪

### 7.2 数据验证
- 输入数据验证
- 业务规则验证
- 数据一致性检查

### 7.3 错误处理
- 统一错误响应格式
- 详细错误信息
- 友好的错误提示 