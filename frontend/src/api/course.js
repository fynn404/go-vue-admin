/**
 * 课程相关API
 * 包含课程管理、选课、成绩等接口
 */

import request from '@/utils/request'

/**
 * 获取所有课程列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.size - 每页数量
 * @returns {Promise} 返回课程列表
 */
export function getCourses(params) {
    return request({
        url: '/courses',
        method: 'get',
        params
    })
}

/**
 * 获取我的课程列表
 * @returns {Promise} 返回我的课程列表
 */
export function getMyCourses() {
    return request({
        url: '/courses/my',
        method: 'get'
    })
}

/**
 * 获取课程详情
 * @param {string} id - 课程ID
 * @returns {Promise} 返回课程详情
 */
export function getCourseDetail(id) {
    return request({
        url: `/courses/${id}`,
        method: 'get'
    })
}

/**
 * 创建课程
 * @param {Object} data - 课程信息
 * @returns {Promise} 返回创建结果
 */
export function createCourse(data) {
    return request({
        url: '/courses',
        method: 'post',
        data
    })
}

/**
 * 更新课程信息
 * @param {string} id - 课程ID
 * @param {Object} data - 课程信息
 * @returns {Promise} 返回更新结果
 */
export function updateCourse(id, data) {
    return request({
        url: `/courses/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除课程
 * @param {string} id - 课程ID
 * @returns {Promise} 返回删除结果
 */
export function deleteCourse(id) {
    return request({
        url: `/courses/${id}`,
        method: 'delete'
    })
}

/**
 * 选课
 * @param {string} courseId - 课程ID
 * @returns {Promise} 返回选课结果
 */
export function enrollCourse(courseId) {
    return request({
        url: `/courses/${courseId}/enroll`,
        method: 'post'
    })
}

/**
 * 退课
 * @param {string} courseId - 课程ID
 * @returns {Promise} 返回退课结果
 */
export function dropCourse(courseId) {
    return request({
        url: `/courses/${courseId}/drop`,
        method: 'post'
    })
}

/**
 * 提交成绩
 * @param {string} courseId - 课程ID
 * @param {Object} data - 成绩信息
 * @returns {Promise} 返回提交结果
 */
export function submitGrade(courseId, data) {
    return request({
        url: `/courses/${courseId}/grades`,
        method: 'post',
        data
    })
} 