# 成绩管理系统数据流说明

## 1. 系统架构层次

系统采用经典的三层架构设计：

```
┌─────────────┐
│   Handler   │ 处理HTTP请求，参数验证，响应封装
├─────────────┤
│   Service   │ 业务逻辑处理，事务管理
├─────────────┤
│    Model    │ 数据访问，ORM映射
└─────────────┘
```

## 2. 数据流转过程

### 2.1 请求处理流程

```
Client Request
     ↓
[Middleware Layer]
 ├── CORS 中间件
 ├── JWT 认证中间件
 └── 角色权限中间件
     ↓
[Handler Layer]
 ├── 参数解析
 ├── 参数验证
 └── 权限检查
     ↓
[Service Layer]
 ├── 业务逻辑处理
 ├── 事务管理
 └── 数据组装
     ↓
[Model Layer]
 ├── 数据库操作
 └── ORM映射
     ↓
Database
```

### 2.2 主要业务流程

#### 2.2.1 成绩录入流程

```
教师录入成绩
     ↓
Handler.Grade.Create
     ↓
Service.Grade.CreateGrade
 ├── 检查课程存在性
 ├── 检查教师权限
 ├── 检查学生存在性
 └── 开启事务
     ↓
创建成绩记录 (Grade)
     ↓
创建历史记录 (GradeHistory)
     ↓
提交事务
```

#### 2.2.2 成绩修改流程

```
教师修改成绩
     ↓
Handler.Grade.Update
     ↓
Service.Grade.UpdateGrade
 ├── 检查成绩存在性
 ├── 检查教师权限
 └── 开启事务
     ↓
更新成绩记录 (Grade)
     ↓
创建历史记录 (GradeHistory)
     ↓
提交事务
```

#### 2.2.3 成绩查询流程

```
用户查询成绩
     ↓
Handler.Grade.GetGrades
 ├── 学生：GetStudentGrades
 └── 教师：GetCourseGrades
     ↓
Service.Grade.GetGrades
 ├── 构建查询条件
 └── 权限过滤
     ↓
Model.Grade.Find
 ├── 关联查询
 └── 数据过滤
```

## 3. 数据一致性保证

### 3.1 事务管理

```
Service层事务示例：

func (s *GradeService) UpdateGrade(...) {
    tx := s.db.Begin()
    
    // 更新成绩
    if err := tx.Save(&grade).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // 创建历史记录
    if err := tx.Create(&history).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit().Error
}
```

### 3.2 并发控制

```
乐观锁示例：

type Grade struct {
    gorm.Model
    Version   uint    // 版本号
    Score     float32 // 分数
    // ... 其他字段
}

// 更新时检查版本号
tx.Model(&grade).Where("version = ?", currentVersion).
    Updates(map[string]interface{}{
        "score": newScore,
        "version": currentVersion + 1,
    })
```

## 4. 数据缓存策略

### 4.1 缓存层次

```
┌─────────────┐
│   Memory    │ 进程内缓存（短期）
├─────────────┤
│    Redis    │ 分布式缓存（中期）
├─────────────┤
│  Database   │ 持久化存储（长期）
└─────────────┘
```

### 4.2 缓存更新策略

```
写操作流程：
1. 更新数据库
2. 删除相关缓存
3. 异步更新相关缓存

读操作流程：
1. 查询缓存
2. 缓存未命中则查询数据库
3. 更新缓存
```

## 5. 错误处理流程

```
错误产生
     ↓
Service层错误包装
     ↓
Handler层错误转换
 ├── 业务错误 → HTTP 4xx
 └── 系统错误 → HTTP 5xx
     ↓
统一错误响应
{
    "code": xxx,
    "message": "xxx",
    "details": {}
}
```

## 6. 数据验证层次

### 6.1 请求验证

```go
type GradeRequest struct {
    Score   float32 `json:"score" binding:"required,min=0,max=100"`
    Comment string  `json:"comment"`
}
```

### 6.2 业务验证

```go
func (s *GradeService) CreateGrade(...) error {
    // 检查课程是否存在
    if err := s.checkCourseExists(...); err != nil {
        return err
    }
    
    // 检查教师权限
    if err := s.checkTeacherPermission(...); err != nil {
        return err
    }
    
    // 检查学生是否存在
    if err := s.checkStudentExists(...); err != nil {
        return err
    }
    
    // ... 其他业务验证
}
```

### 6.3 数据库约束

```sql
-- 唯一索引
CREATE UNIQUE INDEX idx_course_student ON grades (course_id, student_id);

-- 外键约束
ALTER TABLE grades
ADD CONSTRAINT fk_grades_course
FOREIGN KEY (course_id) REFERENCES courses(id);
```

## 7. 性能优化策略

### 7.1 数据库优化

```
1. 索引优化
   - 单列索引：grade_id, course_id, student_id
   - 复合索引：(course_id, student_id)

2. 查询优化
   - 预加载关联数据
   - 分页查询
   - 延迟加载
```

### 7.2 API 性能优化

```
1. 数据压缩
2. 响应缓存
3. 批量操作接口
4. 异步处理
``` 