# 课程管理系统前端

基于 Vue 3 + Vite + Element Plus 的课程管理系统前端项目。

## 项目结构

```
frontend/
├── src/                    # 源代码目录
│   ├── api/               # API 接口定义
│   │   ├── auth.js        # 认证相关接口
│   │   └── course.js      # 课程相关接口
│   ├── assets/            # 静态资源
│   │   └── images/        # 图片资源
│   ├── components/        # 公共组件
│   │   ├── NavHeader.vue  # 导航头部
│   │   └── NavMenu.vue    # 导航菜单
│   ├── composables/       # 组合式函数
│   │   └── useAuth.js     # 认证相关逻辑
│   ├── constants/         # 常量定义
│   │   └── index.js       # 全局常量
│   ├── layouts/           # 布局组件
│   │   └── DefaultLayout.vue # 默认布局
│   ├── router/            # 路由配置
│   │   └── index.js       # 路由定义
│   ├── stores/            # 状态管理
│   │   ├── auth.js        # 认证状态
│   │   └── course.js      # 课程状态
│   ├── styles/            # 全局样式
│   │   ├── element.scss   # Element Plus 主题
│   │   └── index.scss     # 全局样式入口
│   ├── utils/             # 工具函数
│   │   ├── request.js     # 请求工具
│   │   └── validate.js    # 验证工具
│   ├── views/             # 页面组件
│   │   ├── LoginView.vue     # 登录页
│   │   ├── RegisterView.vue  # 注册页
│   │   └── DashboardView.vue # 仪表盘
│   ├── App.vue            # 根组件
│   └── main.js            # 入口文件
├── .env                   # 环境变量
├── .eslintrc.js          # ESLint 配置
├── .prettierrc           # Prettier 配置
├── index.html            # HTML 模板
├── package.json          # 项目依赖
└── vite.config.js        # Vite 配置
```

## 开发规范

1. 文件命名
   - 组件文件使用 PascalCase（如 `NavMenu.vue`）
   - 其他文件使用 camelCase（如 `useAuth.js`）
   - 样式文件使用 kebab-case（如 `element.scss`）

2. 代码风格
   - 使用 ESLint + Prettier 进行代码格式化
   - 组件使用 Composition API
   - 使用 TypeScript 类型注解
   - 所有组件、函数都要有注释说明

3. 提交规范
   - feat: 新功能
   - fix: 修复bug
   - docs: 文档更新
   - style: 代码格式（不影响代码运行的变动）
   - refactor: 重构（既不是新增功能，也不是修改bug的代码变动）
   - test: 增加测试
   - chore: 构建过程或辅助工具的变动

## 开发指南

1. 安装依赖
```bash
npm install
```

2. 启动开发服务器
```bash
npm run dev
```

3. 构建生产版本
```bash
npm run build
```

4. 代码检查
```bash
npm run lint
``` 