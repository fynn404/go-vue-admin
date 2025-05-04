# Go-Vue-Admin

一个基于 Go + Vue3 的课程管理系统，提供课程管理、选课、成绩管理等功能。

## 功能特点

### 用户管理
- 多角色支持（学生、教师、管理员）
- 用户认证与授权
- 个人信息管理

### 课程管理
- 课程创建与编辑
- 课程状态管理（开放/关闭）
- 课程容量控制
- 课程列表查看与搜索

### 选课系统
- 学生选课功能
- 课程容量限制
- 选课状态追踪
- 已选课程管理

### 成绩管理
- 教师成绩录入
- 成绩统计分析
  - 课程成绩统计（平均分、最高分、及格率等）
  - 学生成绩统计（GPA、学分统计等）
- 成绩分布可视化
- 成绩历史记录

## 技术栈

### 后端
- Go
- Gin Web Framework
- GORM
- JWT Authentication
- MySQL

### 前端
- Vue 3
- Vue Router
- Pinia
- Element Plus
- ECharts
- Axios

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
├── pkg/            # 公共包
│   ├── auth/       # 认证相关
│   ├── database/   # 数据库工具
│   └── utils/      # 工具函数
└── frontend/       # 前端项目
    ├── src/
    │   ├── api/        # API 调用
    │   ├── assets/     # 静态资源
    │   ├── components/ # 组件
    │   ├── composables/# 组合式函数
    │   ├── constants/  # 常量定义
    │   ├── router/     # 路由配置
    │   ├── stores/     # 状态管理
    │   ├── utils/      # 工具函数
    │   └── views/      # 页面视图
    └── public/         # 公共静态资源
```

## 开发环境要求

- Go 1.16+
- Node.js 14+
- MySQL 5.7+

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/yourusername/go-vue-admin.git
cd go-vue-admin
```

2. 安装依赖
```bash
# 后端依赖
go mod download

# 前端依赖
cd frontend
npm install
```

3. 配置数据库
- 创建 MySQL 数据库
- 修改 `configs/config.ini` 中的数据库配置

4. 运行项目
```bash
# 后端（在项目根目录下）
go run cmd/main.go

# 前端（在 frontend 目录下）
npm run dev
```

5. 访问系统
- 前端开发服务器：http://localhost:5173
- 后端 API 服务器：http://localhost:8080

## 贡献指南

欢迎提交 Issue 和 Pull Request。

## 许可证

MIT License