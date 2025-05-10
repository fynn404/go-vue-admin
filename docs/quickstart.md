# 快速开始

本指南将帮助你快速搭建和运行 Go-Vue-Admin 项目。

## 前置条件

确保你的系统已安装以下软件：

- Go 1.16+
- Node.js 14+
- MySQL 5.7+
- Git
- Make（可选）
- Docker（可选）

## 获取代码

```bash
# 克隆项目
git clone https://github.com/yourusername/go-vue-admin.git
cd go-vue-admin
```

## 后端配置

1. 安装依赖
```bash
go mod download
```

2. 配置数据库
```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE go_vue_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

3. 修改配置文件
复制配置文件模板：
```bash
cp configs/config.example.ini configs/config.ini
```

编辑 `configs/config.ini`：
```ini
[database]
host = localhost
port = 3306
username = root
password = your_password
dbname = go_vue_admin

[server]
port = 8080
mode = debug

[jwt]
secret = your_jwt_secret
expire = 24h
```

## 前端配置

1. 安装依赖
```bash
cd frontend
npm install
```

2. 配置环境变量
```bash
cp .env.example .env.local
```

编辑 `.env.local`：
```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## 启动服务

### 使用 Make（推荐）

如果你已安装 Make，可以使用以下命令：

```bash
# 启动开发环境（前后端）
make dev

# 或分别启动
make backend-dev  # 启动后端
make frontend-dev # 启动前端
```

### 手动启动

1. 启动后端服务
```bash
# 在项目根目录下
go run cmd/main.go
```

2. 启动前端服务
```bash
# 在 frontend 目录下
npm run dev
```

## 访问系统

启动成功后，可以通过以下地址访问系统：

- 前端页面：http://localhost:5173
- 后端 API：http://localhost:8080

### 默认账号

系统初始化后会创建以下默认账号：

1. 管理员账号
   - 用户名：admin
   - 密码：admin123

2. 测试教师账号
   - 用户名：teacher
   - 密码：teacher123

3. 测试学生账号
   - 用户名：student
   - 密码：student123

## 开发建议

1. 代码规范
   - 遵循项目既定的代码规范
   - 使用 ESLint 和 Prettier 进行代码格式化
   - 遵循 Git Commit 规范

2. 分支管理
   - 主分支：main
   - 开发分支：dev
   - 功能分支：feature/*
   - 修复分支：hotfix/*

3. 开发流程
   - 从 dev 分支创建功能分支
   - 开发完成后提交 Pull Request
   - 通过代码审查后合并到 dev 分支

## 常见问题

### 1. 后端启动失败
- 检查数据库配置是否正确
- 确保数据库服务已启动
- 检查端口是否被占用

### 2. 前端启动失败
- 检查 Node.js 版本是否符合要求
- 确保已安装所有依赖
- 检查环境变量配置

### 3. 跨域问题
- 确保后端 CORS 配置正确
- 检查前端 API 地址配置

## 下一步

- 阅读[开发环境配置](./development/environment-setup.md)了解详细的环境配置
- 查看[后端开发指南](./development/backend-guide.md)了解后端开发详情
- 查看[前端开发指南](./development/frontend-guide.md)了解前端开发详情
- 参考[API 文档](./development/api-docs.md)了解接口定义 