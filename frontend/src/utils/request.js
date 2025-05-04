/**
 * HTTP请求工具
 * 基于axios的请求封装，提供：
 * - 统一的请求配置
 * - 请求拦截器（添加token等）
 * - 响应拦截器（错误处理等）
 * - 统一的错误处理
 */

import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

// 创建axios实例
const service = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1', // API基础路径
    timeout: 15000, // 请求超时时间
    headers: {
        'Content-Type': 'application/json'
    }
})

// 请求拦截器
service.interceptors.request.use(
    config => {
        // 从状态管理中获取token
        const authStore = useAuthStore()
        if (authStore.token) {
            config.headers['Authorization'] = `Bearer ${authStore.token}`
        }
        return config
    },
    error => {
        console.error('请求错误:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
service.interceptors.response.use(
    response => {
        const res = response.data
        // 如果响应成功但业务状态码不为0，显示错误信息
        if (res.code !== 0) {
            ElMessage({
                message: res.message || '操作失败',
                type: 'error',
                duration: 5 * 1000
            })
            return Promise.reject(new Error(res.message || '操作失败'))
        }
        return res
    },
    error => {
        console.error('响应错误:', error)
        // 处理401未授权错误
        if (error.response?.status === 401) {
            const authStore = useAuthStore()
            authStore.logout()
            window.location.href = '/login'
            return
        }
        // 显示错误消息
        ElMessage({
            message: error.response?.data?.message || '系统错误',
            type: 'error',
            duration: 5 * 1000
        })
        return Promise.reject(error)
    }
)

export default service 