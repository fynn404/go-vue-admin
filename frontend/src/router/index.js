/**
 * 路由配置文件
 * 实现前端路由管理，包括：
 * - 路由定义和组件映射
 * - 路由守卫（权限控制）
 * - 路由懒加载
 */

import { createRouter, createWebHistory } from 'vue-router'  // Vue Router核心功能
import { useAuthStore } from '@/stores/auth'                 // 认证状态管理

// 创建路由实例
const router = createRouter({
    // 使用HTML5历史模式（需要服务器配置支持）
    history: createWebHistory(import.meta.env.BASE_URL),
    // 路由配置数组
    routes: [
        // 登录页面 - 无需认证
        {
            path: '/login',
            name: 'login',
            component: () => import('@/views/LoginView.vue'),  // 懒加载组件
            meta: { requiresAuth: false }                      // 元信息：无需认证
        },
        // 注册页面 - 无需认证
        {
            path: '/register',
            name: 'register',
            component: () => import('@/views/RegisterView.vue'),
            meta: { requiresAuth: false }
        },
        // 仪表盘 - 需要认证
        {
            path: '/',
            name: 'dashboard',
            component: () => import('@/views/DashboardView.vue'),
            meta: { requiresAuth: true }                       // 元信息：需要认证
        },
        // 个人信息页面 - 需要认证
        {
            path: '/profile',
            name: 'profile',
            component: () => import('@/views/ProfileView.vue'),
            meta: { requiresAuth: true }
        },
        // 选课中心 - 仅学生可访问
        {
            path: '/courses',
            name: 'courses',
            component: () => import('@/views/CoursesView.vue'),
            meta: { requiresAuth: true, roles: ['student'] }   // 元信息：需要认证且角色为学生
        },
        // 我的课程 - 仅学生可访问
        {
            path: '/my-courses',
            name: 'my-courses',
            component: () => import('@/views/MyCoursesView.vue'),
            meta: { requiresAuth: true, roles: ['student'] }
        },
        // 课程管理 - 仅教师可访问
        {
            path: '/teaching',
            name: 'teaching',
            component: () => import('@/views/TeachingView.vue'),
            meta: { requiresAuth: true, roles: ['teacher'] }   // 元信息：需要认证且角色为教师
        },
        // 学生管理 - 仅教师可访问
        {
            path: '/students',
            name: 'students',
            component: () => import('@/views/StudentsView.vue'),
            meta: { requiresAuth: true, roles: ['teacher'] }
        },
        // 用户管理 - 仅管理员可访问
        {
            path: '/users',
            name: 'users',
            component: () => import('@/views/UsersView.vue'),
            meta: { requiresAuth: true, roles: ['admin'] }     // 元信息：需要认证且角色为管理员
        },
        // 系统设置 - 仅管理员可访问
        {
            path: '/settings',
            name: 'settings',
            component: () => import('@/views/SettingsView.vue'),
            meta: { requiresAuth: true, roles: ['admin'] }
        }
    ]
})

/**
 * 全局前置守卫
 * 在路由跳转前进行权限检查
 * @param {Route} to - 目标路由
 * @param {Route} from - 来源路由
 * @param {Function} next - 确认导航的函数
 */
router.beforeEach((to, from, next) => {
    // 获取认证状态
    const authStore = useAuthStore()
    const isAuthenticated = authStore.isAuthenticated
    const userRole = authStore.user?.role

    // 检查路由是否需要认证
    if (to.meta.requiresAuth && !isAuthenticated) {
        next('/login')  // 未认证时重定向到登录页
        return
    }

    // 检查用户是否具有所需角色
    if (to.meta.roles && !to.meta.roles.includes(userRole)) {
        next('/')      // 角色不匹配时重定向到首页
        return
    }

    // 已登录用户访问登录/注册页面时重定向到首页
    if (isAuthenticated && ['/login', '/register'].includes(to.path)) {
        next('/')
        return
    }

    // 允许导航继续
    next()
})

// 导出路由实例
export default router 