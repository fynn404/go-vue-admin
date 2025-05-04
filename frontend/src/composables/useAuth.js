/**
 * 认证相关的组合式函数
 * 提供认证状态管理和用户信息处理的逻辑
 */

import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { USER_ROLES, ROLE_LABELS } from '@/constants'

export function useAuth() {
    const router = useRouter()
    const authStore = useAuthStore()

    // 计算属性：是否已认证
    const isAuthenticated = computed(() => authStore.isAuthenticated)

    // 计算属性：当前用户
    const currentUser = computed(() => authStore.user)

    // 计算属性：用户角色标签
    const roleLabel = computed(() => {
        const role = currentUser.value?.role
        return ROLE_LABELS[role] || '未知角色'
    })

    // 计算属性：角色判断
    const isStudent = computed(() => currentUser.value?.role === USER_ROLES.STUDENT)
    const isTeacher = computed(() => currentUser.value?.role === USER_ROLES.TEACHER)
    const isAdmin = computed(() => currentUser.value?.role === USER_ROLES.ADMIN)

    /**
     * 处理登录
     * @param {Object} credentials - 登录凭证
     * @returns {Promise<boolean>} 登录结果
     */
    const handleLogin = async (credentials) => {
        try {
            const success = await authStore.login(credentials)
            if (success) {
                router.push('/')
            }
            return success
        } catch (error) {
            console.error('登录失败:', error)
            return false
        }
    }

    /**
     * 处理注册
     * @param {Object} data - 注册信息
     * @returns {Promise<boolean>} 注册结果
     */
    const handleRegister = async (data) => {
        try {
            const success = await authStore.register(data)
            if (success) {
                // 注册成功后自动登录
                await handleLogin({
                    username: data.username,
                    password: data.password
                })
            }
            return success
        } catch (error) {
            console.error('注册失败:', error)
            return false
        }
    }

    /**
     * 处理登出
     */
    const handleLogout = async () => {
        try {
            await authStore.logout()
            router.push('/login')
        } catch (error) {
            console.error('登出失败:', error)
        }
    }

    /**
     * 检查是否有权限访问
     * @param {string[]} allowedRoles - 允许的角色列表
     * @returns {boolean} 是否有权限
     */
    const hasPermission = (allowedRoles) => {
        const userRole = currentUser.value?.role
        return allowedRoles.includes(userRole)
    }

    return {
        isAuthenticated,
        currentUser,
        roleLabel,
        isStudent,
        isTeacher,
        isAdmin,
        handleLogin,
        handleRegister,
        handleLogout,
        hasPermission
    }
} 