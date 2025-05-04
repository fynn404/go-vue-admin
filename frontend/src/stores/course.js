/**
 * 课程状态管理
 * 使用 Pinia 管理课程相关状态，包括：
 * - 课程列表
 * - 选课记录
 * - 课程详情
 * - 成绩信息
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
    getCourses,
    getMyCourses,
    enrollCourse,
    dropCourse,
    updateCourse
} from '@/api/course'

export const useCourseStore = defineStore('course', () => {
    // 状态
    const courses = ref([])
    const myCourses = ref([])
    const currentCourse = ref(null)

    // 计算属性
    const availableCourses = computed(() => {
        return courses.value.filter(course =>
            !myCourses.value.some(mc => mc.id === course.id)
        )
    })

    // 方法
    const fetchCourses = async () => {
        try {
            const response = await getCourses()
            courses.value = response.data
            return true
        } catch (error) {
            console.error('获取课程列表失败:', error)
            return false
        }
    }

    const fetchMyCourses = async () => {
        try {
            const response = await getMyCourses()
            myCourses.value = response.data
            return true
        } catch (error) {
            console.error('获取我的课程失败:', error)
            return false
        }
    }

    const enroll = async (courseId) => {
        try {
            await enrollCourse(courseId)
            await fetchMyCourses()
            return true
        } catch (error) {
            console.error('选课失败:', error)
            return false
        }
    }

    const drop = async (courseId) => {
        try {
            await dropCourse(courseId)
            await fetchMyCourses()
            return true
        } catch (error) {
            console.error('退课失败:', error)
            return false
        }
    }

    const updateCourseInfo = async (courseId, data) => {
        try {
            await updateCourse(courseId, data)
            await fetchCourses()
            return true
        } catch (error) {
            console.error('更新课程信息失败:', error)
            return false
        }
    }

    return {
        courses,
        myCourses,
        currentCourse,
        availableCourses,
        fetchCourses,
        fetchMyCourses,
        enroll,
        drop,
        updateCourseInfo
    }
}) 