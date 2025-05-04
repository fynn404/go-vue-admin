/**
 * 表单验证工具函数
 * 提供常用的表单验证规则和方法
 */

/**
 * 验证用户名
 * @param {string} username - 用户名
 * @returns {boolean} 是否有效
 */
export function validateUsername(username) {
    // 用户名长度4-20位，只能包含字母、数字和下划线
    const reg = /^[a-zA-Z0-9_]{4,20}$/
    return reg.test(username)
}

/**
 * 验证密码强度
 * @param {string} password - 密码
 * @returns {boolean} 是否有效
 */
export function validatePassword(password) {
    // 密码长度8-20位，必须包含字母和数字
    const reg = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,20}$/
    return reg.test(password)
}

/**
 * 验证邮箱
 * @param {string} email - 邮箱地址
 * @returns {boolean} 是否有效
 */
export function validateEmail(email) {
    const reg = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return reg.test(email)
}

/**
 * 验证手机号
 * @param {string} phone - 手机号
 * @returns {boolean} 是否有效
 */
export function validatePhone(phone) {
    const reg = /^1[3-9]\d{9}$/
    return reg.test(phone)
}

/**
 * 验证URL
 * @param {string} url - URL地址
 * @returns {boolean} 是否有效
 */
export function validateURL(url) {
    const reg = /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/
    return reg.test(url)
}

/**
 * Element Plus表单验证规则
 */
export const rules = {
    // 用户名验证规则
    username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 4, max: 20, message: '长度在4到20个字符', trigger: 'blur' },
        { pattern: /^[a-zA-Z0-9_]{4,20}$/, message: '只能包含字母、数字和下划线', trigger: 'blur' }
    ],
    // 密码验证规则
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 8, max: 20, message: '长度在8到20个字符', trigger: 'blur' },
        { pattern: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,20}$/, message: '必须包含字母和数字', trigger: 'blur' }
    ],
    // 邮箱验证规则
    email: [
        { required: true, message: '请输入邮箱地址', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
    ],
    // 手机号验证规则
    phone: [
        { required: true, message: '请输入手机号', trigger: 'blur' },
        { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
    ]
} 