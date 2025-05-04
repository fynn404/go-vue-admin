/**
 * 认证状态管理
 * 使用 Pinia 管理用户认证状态，包括：
 * - 用户登录状态
 * - 用户信息
 * - Token管理
 * - 登录/登出操作
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, logout as apiLogout } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
    // 状态
    const token = ref(localStorage.getItem('token') || '')
    const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

    // 计算属性
    const isAuthenticated = computed(() => !!token.value)

    // 方法
    const setToken = (newToken) => {
        token.value = newToken
        localStorage.setItem('token', newToken)
    }

    const setUser = (newUser) => {
        user.value = newUser
        localStorage.setItem('user', JSON.stringify(newUser))
    }

    const login = async (credentials) => {
        try {
            const response = await apiLogin(credentials)
            setToken(response.token)
            setUser(response.user)
            return true
        } catch (error) {
            console.error('登录失败:', error)
            return false
        }
    }

    const logout = async () => {
        try {
            await apiLogout()
        } finally {
            // 无论是否成功调用登出API，都清除本地状态
            token.value = ''
            user.value = null
            localStorage.removeItem('token')
            localStorage.removeItem('user')
        }
    }

    return {
        token,
        user,
        isAuthenticated,
        login,
        logout
    }
}) 