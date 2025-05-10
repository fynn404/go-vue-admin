# 开发环境配置

本文档详细说明如何配置 Go-Vue-Admin 项目的开发环境。

## 系统要求

### 操作系统
- macOS 10.15+
- Ubuntu 18.04+
- Windows 10+

### 必要软件
1. Go
   - 版本：1.16+
   - 下载：[Go 官网](https://golang.org/dl/)
   - 环境变量配置

2. Node.js
   - 版本：14+
   - 下载：[Node.js 官网](https://nodejs.org/)
   - 推荐使用 nvm 管理版本

3. MySQL
   - 版本：5.7+
   - 下载：[MySQL 官网](https://dev.mysql.com/downloads/)
   - 配置说明

4. Git
   - 最新版本
   - 下载：[Git 官网](https://git-scm.com/)
   - 基础配置

### 推荐工具
1. IDE/编辑器
   - GoLand（推荐）
   - VSCode
   - Sublime Text

2. 数据库工具
   - MySQL Workbench
   - Navicat
   - DBeaver

3. API 测试工具
   - Postman
   - Insomnia
   - curl

4. Git 客户端
   - SourceTree
   - GitKraken
   - GitHub Desktop

## 详细安装步骤

### macOS

1. 使用 Homebrew 安装
```bash
# 安装 Homebrew（如果未安装）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装必要软件
brew install go
brew install node
brew install mysql
brew install git

# 启动 MySQL
brew services start mysql
```

2. Go 环境配置
```bash
# 设置 GOPATH
echo 'export GOPATH=$HOME/go' >> ~/.zshrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
source ~/.zshrc

# 设置 Go 代理（国内开发者建议）
go env -w GOPROXY=https://goproxy.cn,direct
```

3. Node.js 环境配置
```bash
# 安装 nvm（推荐）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# 安装 Node.js
nvm install 14
nvm use 14

# 设置 npm 镜像（国内开发者建议）
npm config set registry https://registry.npmmirror.com
```

### Ubuntu

1. 安装必要软件
```bash
# 更新包列表
sudo apt update

# 安装 Go
wget https://golang.org/dl/go1.16.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_14.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装 MySQL
sudo apt install mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# 安装 Git
sudo apt install git
```

2. MySQL 安全配置
```bash
sudo mysql_secure_installation
```

### Windows

1. 使用安装包
- 下载并安装 [Go](https://golang.org/dl/)
- 下载并安装 [Node.js](https://nodejs.org/)
- 下载并安装 [MySQL](https://dev.mysql.com/downloads/installer/)
- 下载并安装 [Git](https://git-scm.com/download/win)

2. 使用 Scoop（推荐）
```powershell
# 安装 Scoop
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex

# 安装必要软件
scoop install go
scoop install nodejs
scoop install mysql
scoop install git
```

## IDE 配置

### GoLand 配置

1. Go 配置
   - 设置 GOROOT
   - 设置 GOPATH
   - 启用 Go Modules

2. 插件安装
   - File Watchers（自动格式化）
   - Go Templates
   - Go Test Explorer

3. 代码风格
   - 启用 gofmt
   - 配置 golangci-lint

### VSCode 配置

1. 必要插件
   - Go
   - Vue Language Features
   - ESLint
   - Prettier
   - GitLens

2. 推荐设置
```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.fixAll.eslint": true
    }
}
```

## 开发工具配置

### Git 配置

1. 基础配置
```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

2. SSH 密钥配置
```bash
ssh-keygen -t ed25519 -C "your.email@example.com"
```

3. Git 提交模板
```bash
git config --global commit.template ~/.gitmessage
```

### 代码规范工具

1. 后端
```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装 goimports
go install golang.org/x/tools/cmd/goimports@latest
```

2. 前端
```bash
# 安装 ESLint
npm install -g eslint

# 安装 Prettier
npm install -g prettier
```

## 常见问题

### Go 相关
1. GOPATH 配置问题
   - 检查环境变量设置
   - 确保目录存在

2. 依赖下载失败
   - 检查网络连接
   - 尝试使用代理

### Node.js 相关
1. npm 安装失败
   - 清除 npm 缓存
   - 使用镜像源

2. 版本兼容问题
   - 使用 nvm 管理版本
   - 检查 package.json

### MySQL 相关
1. 启动失败
   - 检查服务状态
   - 查看错误日志

2. 权限问题
   - 重置 root 密码
   - 配置用户权限

## 下一步

- 阅读[后端开发指南](./backend-guide.md)
- 阅读[前端开发指南](./frontend-guide.md)
- 查看[代码规范](./coding-standards.md) 