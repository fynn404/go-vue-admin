import request from '@/utils/request'

// 成绩基础操作
export function getGradeList(params) {
    return request({
        url: '/api/grades',
        method: 'get',
        params
    })
}

export function createGrade(data) {
    return request({
        url: '/api/grades',
        method: 'post',
        data
    })
}

export function updateGrade(id, data) {
    return request({
        url: `/api/grades/${id}`,
        method: 'put',
        data
    })
}

export function deleteGrade(id) {
    return request({
        url: `/api/grades/${id}`,
        method: 'delete'
    })
}

// 成绩统计相关
export function getCourseStats(courseId) {
    return request({
        url: `/api/grades/stats/course/${courseId}`,
        method: 'get'
    })
}

export function getStudentStats(studentId) {
    return request({
        url: `/api/grades/stats/student/${studentId}`,
        method: 'get'
    })
}

// 成绩历史记录相关
export function getGradeHistory(gradeId, params) {
    return request({
        url: `/api/grades/${gradeId}/history`,
        method: 'get',
        params
    })
}

export function getCourseGradeHistory(courseId, params) {
    return request({
        url: `/api/grades/history/course/${courseId}`,
        method: 'get',
        params
    })
}

export function getStudentGradeHistory(studentId, params) {
    return request({
        url: `/api/grades/history/student/${studentId}`,
        method: 'get',
        params
    })
}

export function getTeacherGradeHistory(teacherId, params) {
    return request({
        url: `/api/grades/history/teacher/${teacherId}`,
        method: 'get',
        params
    })
} 