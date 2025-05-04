/**
 * 认证相关API
 * 包含用户认证、注册、个人信息等接口
 */

import request from '@/utils/request'

/**
 * 用户登录
 * @param {Object} data - 登录信息
 * @param {string} data.username - 用户名
 * @param {string} data.password - 密码
 * @returns {Promise} 返回登录结果
 */
export function login(data) {
    return request({
        url: '/auth/login',
        method: 'post',
        data
    })
}

/**
 * 用户注册
 * @param {Object} data - 注册信息
 * @param {string} data.username - 用户名
 * @param {string} data.password - 密码
 * @param {string} data.role - 角色
 * @returns {Promise} 返回注册结果
 */
export function register(data) {
    return request({
        url: '/auth/register',
        method: 'post',
        data
    })
}

/**
 * 用户登出
 * @returns {Promise} 返回登出结果
 */
export function logout() {
    return request({
        url: '/auth/logout',
        method: 'post'
    })
}

/**
 * 获取用户信息
 * @returns {Promise} 返回用户信息
 */
export function getUserInfo() {
    return request({
        url: '/users/profile',
        method: 'get'
    })
}

/**
 * 更新用户信息
 * @param {Object} data - 用户信息
 * @returns {Promise} 返回更新结果
 */
export function updateUserInfo(data) {
    return request({
        url: '/users/profile',
        method: 'put',
        data
    })
}

/**
 * 修改密码
 * @param {Object} data - 密码信息
 * @param {string} data.oldPassword - 原密码
 * @param {string} data.newPassword - 新密码
 * @returns {Promise} 返回修改结果
 */
export function changePassword(data) {
    return request({
        url: '/users/password',
        method: 'put',
        data
    })
} 