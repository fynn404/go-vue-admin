# 前端开发指南

本文档详细说明 Go-Vue-Admin 项目的前端开发规范和指南。

## 项目结构

```
frontend/
├── public/          # 静态资源
├── src/
│   ├── api/         # API 接口
│   ├── assets/      # 资源文件
│   ├── components/  # 公共组件
│   ├── composables/ # 组合式函数
│   ├── constants/   # 常量定义
│   ├── router/      # 路由配置
│   ├── stores/      # 状态管理
│   ├── utils/       # 工具函数
│   └── views/       # 页面视图
├── .env             # 环境变量
├── index.html       # HTML 模板
├── package.json     # 项目配置
└── vite.config.ts   # Vite 配置
```

## 技术栈

- Vue 3
- Vue Router
- Pinia
- Element Plus
- Axios
- ECharts
- TypeScript
- Vite

## 开发规范

### 命名规范

1. 文件命名
   - 组件文件：PascalCase（如 `UserList.vue`）
   - 工具文件：camelCase（如 `formatDate.ts`）
   - 样式文件：kebab-case（如 `user-list.scss`）

2. 组件命名
   - 使用 PascalCase
   - 有意义的前缀（如 `BaseButton`、`UserForm`）
   - 避免单个单词

3. 变量命名
   - 使用 camelCase
   - 布尔值以 is/has/can 开头
   - 常量使用 UPPER_SNAKE_CASE

### 目录结构规范

1. 组件目录
```
components/
├── base/           # 基础组件
├── layout/         # 布局组件
└── business/       # 业务组件
```

2. 视图目录
```
views/
├── user/           # 用户相关
├── course/         # 课程相关
└── dashboard/      # 仪表盘
```

3. API 目录
```
api/
├── user.ts         # 用户相关接口
├── course.ts       # 课程相关接口
└── types.ts        # 类型定义
```

## 组件开发

### 组件模板

```vue
<template>
  <div class="user-list">
    <!-- 组件内容 -->
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { User } from '@/types'

// 状态定义
const users = ref<User[]>([])

// 方法定义
const fetchUsers = async () => {
  // 实现逻辑
}

// 生命周期
onMounted(() => {
  fetchUsers()
})
</script>

<style scoped lang="scss">
.user-list {
  // 样式定义
}
</style>
```

### Props 定义

```typescript
interface Props {
  title: string
  users?: User[]
  onSelect?: (user: User) => void
}

defineProps<Props>()
```

### 事件处理

```typescript
const emit = defineEmits<{
  (e: 'update', value: string): void
  (e: 'select', item: any): void
}>()
```

## 状态管理

### Store 定义

```typescript
// stores/user.ts
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    currentUser: null,
    token: null
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token
  },
  
  actions: {
    async login(username: string, password: string) {
      // 实现登录逻辑
    }
  }
})
```

### Store 使用

```vue
<script setup lang="ts">
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const handleLogin = async () => {
  await userStore.login(username.value, password.value)
}
</script>
```

## 路由配置

### 路由定义

```typescript
// router/index.ts
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: {
      requiresAuth: true
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})
```

### 路由守卫

```typescript
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})
```

## API 调用

### API 配置

```typescript
// utils/request.ts
import axios from 'axios'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 5000
})

request.interceptors.request.use(
  config => {
    const token = useUserStore().token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

request.interceptors.response.use(
  response => response.data,
  error => Promise.reject(error)
)

export default request
```

### API 定义

```typescript
// api/user.ts
import request from '@/utils/request'
import type { User, LoginRequest, LoginResponse } from './types'

export const login = (data: LoginRequest) =>
  request.post<LoginResponse>('/auth/login', data)

export const getUserInfo = () =>
  request.get<User>('/users/me')
```

## 样式指南

### 样式架构

1. 全局样式
```scss
// styles/index.scss
@import 'variables';
@import 'mixins';
@import 'reset';
@import 'common';
```

2. 组件样式
```scss
.component-name {
  &__element {
    // 样式定义
  }
  
  &--modifier {
    // 样式定义
  }
}
```

### 主题定义

```scss
// styles/variables.scss
:root {
  // 颜色
  --primary-color: #409EFF;
  --success-color: #67C23A;
  --warning-color: #E6A23C;
  --danger-color: #F56C6C;
  
  // 字体
  --font-size-base: 14px;
  --font-size-large: 16px;
  --font-size-small: 12px;
  
  // 间距
  --spacing-base: 8px;
  --spacing-large: 16px;
  --spacing-small: 4px;
}
```

## 工具函数

### 日期格式化

```typescript
// utils/date.ts
export const formatDate = (date: Date, format = 'YYYY-MM-DD'): string => {
  // 实现格式化逻辑
}
```

### 数据验证

```typescript
// utils/validate.ts
export const isEmail = (email: string): boolean => {
  // 实现验证逻辑
}
```

## 测试指南

### 单元测试

```typescript
// __tests__/components/UserList.spec.ts
import { mount } from '@vue/test-utils'
import UserList from '@/components/UserList.vue'

describe('UserList', () => {
  test('renders user list correctly', () => {
    const wrapper = mount(UserList, {
      props: {
        users: [/* 测试数据 */]
      }
    })
    
    expect(wrapper.findAll('.user-item')).toHaveLength(2)
  })
})
```

### E2E 测试

```typescript
// cypress/e2e/login.cy.ts
describe('Login', () => {
  it('should login successfully', () => {
    cy.visit('/login')
    cy.get('input[name=username]').type('admin')
    cy.get('input[name=password]').type('password')
    cy.get('button[type=submit]').click()
    cy.url().should('include', '/dashboard')
  })
})
```

## 性能优化

1. 代码分割
```typescript
// router/index.ts
const routes = [
  {
    path: '/dashboard',
    component: () => import('@/views/Dashboard.vue')
  }
]
```

2. 图片优化
```vue
<template>
  <img
    v-lazy="imageUrl"
    alt="描述"
  >
</template>
```

3. 虚拟列表
```vue
<template>
  <el-virtual-list
    :items="items"
    :item-size="50"
  >
    <template #default="{ item }">
      <!-- 列表项内容 -->
    </template>
  </el-virtual-list>
</template>
```

## 部署配置

### 环境变量

```
# .env.development
VITE_API_BASE_URL=http://localhost:8080/api/v1

# .env.production
VITE_API_BASE_URL=/api/v1
```

### 构建配置

```typescript
// vite.config.ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  build: {
    target: 'es2015',
    minify: 'terser',
    cssCodeSplit: true
  }
})
```

## 最佳实践

1. 组件设计
   - 单一职责
   - 可复用性
   - Props 验证
   - 事件命名

2. 性能优化
   - 懒加载
   - 缓存
   - 防抖节流
   - 组件缓存

3. 代码质量
   - TypeScript
   - ESLint
   - Prettier
   - 单元测试

4. 安全考虑
   - XSS 防护
   - CSRF 防护
   - 敏感信息加密

## 常见问题

1. 状态管理
   - 状态共享
   - 数据持久化
   - 状态同步

2. 路由问题
   - 权限控制
   - 参数传递
   - 路由缓存

3. 性能问题
   - 首屏加载
   - 内存泄漏
   - 渲染优化

## 参考资源

- [Vue 3 文档](https://vuejs.org/)
- [Vue Router 文档](https://router.vuejs.org/)
- [Pinia 文档](https://pinia.vuejs.org/)
- [Element Plus 文档](https://element-plus.org/) 