/**
 * 全局常量定义
 */

// 用户角色
export const USER_ROLES = {
    STUDENT: 'student',
    TEACHER: 'teacher',
    ADMIN: 'admin'
}

// 角色标签类型
export const ROLE_TAG_TYPES = {
    [USER_ROLES.STUDENT]: 'success',
    [USER_ROLES.TEACHER]: 'warning',
    [USER_ROLES.ADMIN]: 'danger'
}

// 角色中文名称
export const ROLE_LABELS = {
    [USER_ROLES.STUDENT]: '学生',
    [USER_ROLES.TEACHER]: '教师',
    [USER_ROLES.ADMIN]: '管理员'
}

// 本地存储键名
export const STORAGE_KEYS = {
    TOKEN: 'token',
    USER: 'user',
    THEME: 'theme',
    LANGUAGE: 'language'
}

// API响应码
export const API_CODE = {
    SUCCESS: 0,
    ERROR: 1,
    UNAUTHORIZED: 401,
    FORBIDDEN: 403,
    NOT_FOUND: 404,
    SERVER_ERROR: 500
}

// 分页配置
export const PAGINATION = {
    PAGE_SIZES: [10, 20, 30, 50],
    DEFAULT_PAGE_SIZE: 10,
    DEFAULT_CURRENT_PAGE: 1
}

// 课程状态
export const COURSE_STATUS = {
    NOT_STARTED: 0,
    IN_PROGRESS: 1,
    FINISHED: 2
}

// 课程状态标签
export const COURSE_STATUS_TAGS = {
    [COURSE_STATUS.NOT_STARTED]: {
        type: 'info',
        label: '未开始'
    },
    [COURSE_STATUS.IN_PROGRESS]: {
        type: 'success',
        label: '进行中'
    },
    [COURSE_STATUS.FINISHED]: {
        type: 'warning',
        label: '已结束'
    }
}

// 成绩等级
export const GRADE_LEVELS = {
    A: 90,
    B: 80,
    C: 70,
    D: 60,
    F: 0
}

// 主题配置
export const THEME = {
    PRIMARY_COLOR: '#409EFF',
    SUCCESS_COLOR: '#67C23A',
    WARNING_COLOR: '#E6A23C',
    DANGER_COLOR: '#F56C6C',
    INFO_COLOR: '#909399',
    BG_COLOR: '#F5F7FA'
} 