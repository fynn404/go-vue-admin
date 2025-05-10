# 贡献指南

感谢你对 Go-Vue-Admin 项目的关注！本文档将指导你如何为项目做出贡献。

## 行为准则

### 基本原则
1. 尊重所有贡献者
2. 保持专业和友善
3. 接受建设性批评
4. 关注社区利益

### 不当行为
1. 人身攻击
2. 骚扰或歧视
3. 垃圾信息
4. 恶意破坏

## 贡献流程

### 1. 准备工作

1. Fork 项目仓库
```bash
# 克隆你的 fork
git clone https://github.com/your-username/go-vue-admin.git
cd go-vue-admin

# 添加上游仓库
git remote add upstream https://github.com/original-owner/go-vue-admin.git
```

2. 创建分支
```bash
# 更新主分支
git checkout main
git pull upstream main

# 创建功能分支
git checkout -b feature/your-feature
```

### 2. 开发规范

1. 代码风格
   - 遵循项目既定的代码规范
   - 使用工具进行格式化
   - 保持代码整洁

2. 提交规范
```bash
# 提交格式
<type>(<scope>): <subject>

# 示例
feat(user): add user registration
fix(course): fix course selection bug
docs(api): update API documentation
```

类型说明：
- feat: 新功能
- fix: 修复问题
- docs: 文档更新
- style: 代码格式
- refactor: 代码重构
- test: 测试相关
- chore: 构建过程或辅助工具的变动

### 3. 测试要求

1. 单元测试
```go
// 示例：用户服务测试
func TestUserService_Create(t *testing.T) {
    // 准备测试数据
    // 执行测试
    // 验证结果
}
```

2. 集成测试
```go
func TestUserAPI(t *testing.T) {
    // 设置测试环境
    // 执行 HTTP 请求
    // 验证响应
}
```

3. 前端测试
```javascript
describe('UserList', () => {
    test('renders user list correctly', () => {
        // 准备测试数据
        // 渲染组件
        // 验证渲染结果
    })
})
```

### 4. 文档要求

1. 代码注释
```go
// UserService 处理用户相关的业务逻辑
type UserService interface {
    // Create 创建新用户
    // 参数：
    //   - ctx: 上下文
    //   - user: 用户信息
    // 返回：
    //   - 创建的用户
    //   - 错误信息
    Create(ctx context.Context, user *User) (*User, error)
}
```

2. API 文档
```markdown
### 创建用户

POST /api/v1/users

请求参数：
- username: 用户名
- password: 密码
- role: 角色

响应：
- 200: 成功
- 400: 参数错误
- 500: 服务器错误
```

3. 更新文档
- 更新 README.md
- 更新 API 文档
- 更新使用说明

### 5. 提交 Pull Request

1. 准备工作
```bash
# 确保代码最新
git fetch upstream
git rebase upstream/main

# 提交更改
git add .
git commit -m "feat(module): your feature description"
git push origin feature/your-feature
```

2. PR 描述模板
```markdown
## 描述
简要描述你的更改

## 类型
- [ ] 功能新增
- [ ] Bug 修复
- [ ] 文档更新
- [ ] 代码重构
- [ ] 其他

## 测试
- [ ] 单元测试
- [ ] 集成测试
- [ ] 手动测试

## 检查清单
- [ ] 代码符合规范
- [ ] 添加测试用例
- [ ] 更新文档
- [ ] 本地测试通过
```

## 开发环境

### 1. 必要软件
- Go 1.16+
- Node.js 14+
- MySQL 5.7+
- Git
- Docker (可选)
- Make (可选)

### 2. 开发工具
- GoLand/VSCode
- MySQL Workbench
- Postman
- Git 客户端

### 3. 环境设置
```bash
# 后端依赖
go mod download

# 前端依赖
cd frontend
npm install
```

## 开发流程

### 1. 功能开发

1. 需求分析
   - 理解需求
   - 设计方案
   - 讨论反馈

2. 开发实现
   - 编写代码
   - 添加测试
   - 本地验证

3. 代码审查
   - 提交 PR
   - 响应反馈
   - 修改完善

### 2. Bug 修复

1. 问题报告
   - 描述问题
   - 复现步骤
   - 期望结果

2. 修复流程
   - 定位问题
   - 修复验证
   - 添加测试

3. 提交修复
   - 创建 PR
   - 更新文档
   - 等待审查

## 发布流程

### 1. 版本规范

1. 版本号格式
   - 主版本号：不兼容的 API 修改
   - 次版本号：向下兼容的功能新增
   - 修订号：向下兼容的问题修复

2. 发布标签
```bash
git tag -a v1.0.0 -m "version 1.0.0"
git push origin v1.0.0
```

### 2. 发布检查

1. 代码检查
   - 代码质量
   - 测试覆盖
   - 性能测试

2. 文档检查
   - API 文档
   - 使用说明
   - 更新日志

3. 依赖检查
   - 版本兼容
   - 安全漏洞
   - 许可证合规

## 社区参与

### 1. 问题反馈

1. 提交 Issue
   - 使用模板
   - 详细描述
   - 提供信息

2. 问题讨论
   - 积极参与
   - 提供建议
   - 分享经验

### 2. 功能建议

1. 提案流程
   - 描述需求
   - 设计方案
   - 收集反馈

2. 实现计划
   - 任务分解
   - 时间估计
   - 资源协调

## 其他说明

### 1. 许可证

本项目采用 MIT 许可证，详见 [LICENSE](./LICENSE) 文件。

### 2. 联系方式

- Issue 系统：问题报告和功能建议
- 邮件列表：技术讨论和公告
- 社区论坛：经验分享和交流

### 3. 资源链接

- [项目文档](./docs)
- [开发指南](./docs/development)
- [API 文档](./docs/api)
- [常见问题](./docs/faq.md) 