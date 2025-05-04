<!--
  DashboardView.vue
  仪表盘页面组件
  
  功能：
  - 显示用户基本信息
  - 显示课程统计信息
  - 显示最近活动
  - 根据用户角色显示不同内容
-->
<template>
  <div class="dashboard-container">
    <!-- 欢迎信息 -->
    <el-card class="welcome-card">
      <template #header>
        <div class="welcome-header">
          <h3>欢迎回来，{{ authStore.user?.username }}</h3>
          <el-tag>{{ roleLabel }}</el-tag>
        </div>
      </template>
      <p>今天是 {{ currentDate }}，祝您使用愉快！</p>
    </el-card>

    <!-- 统计数据 -->
    <el-row :gutter="20" class="statistics">
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="stat-header">
              <el-icon><Calendar /></el-icon>
              <span>{{ statTitle }}</span>
            </div>
          </template>
          <div class="stat-content">
            <h2>{{ courseCount }}</h2>
            <p>{{ statDescription }}</p>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="stat-header">
              <el-icon><Clock /></el-icon>
              <span>本周活动</span>
            </div>
          </template>
          <div class="stat-content">
            <h2>{{ weeklyActivities }}</h2>
            <p>条记录</p>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="stat-header">
              <el-icon><Bell /></el-icon>
              <span>待办事项</span>
            </div>
          </template>
          <div class="stat-content">
            <h2>{{ todoCount }}</h2>
            <p>项待处理</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近活动 -->
    <el-card class="recent-activities">
      <template #header>
        <div class="card-header">
          <h3>最近活动</h3>
        </div>
      </template>
      
      <el-timeline>
        <el-timeline-item
          v-for="activity in recentActivities"
          :key="activity.id"
          :timestamp="activity.time"
          :type="activity.type"
        >
          {{ activity.content }}
        </el-timeline-item>
      </el-timeline>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useCourseStore } from '@/stores/course'
import { Calendar, Clock, Bell } from '@element-plus/icons-vue'

// 状态管理
const authStore = useAuthStore()
const courseStore = useCourseStore()

// 计算属性
const roleLabel = computed(() => {
  const roleMap = {
    student: '学生',
    teacher: '教师',
    admin: '管理员'
  }
  return roleMap[authStore.user?.role] || '未知角色'
})

const currentDate = computed(() => {
  return new Date().toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
})

const statTitle = computed(() => {
  const role = authStore.user?.role
  return role === 'student' ? '已选课程' : '管理课程'
})

const statDescription = computed(() => {
  const role = authStore.user?.role
  return role === 'student' ? '门课程' : '门课程'
})

// 响应式数据
const courseCount = ref(0)
const weeklyActivities = ref(0)
const todoCount = ref(0)
const recentActivities = ref([
  {
    id: 1,
    content: '登录系统',
    time: new Date().toLocaleTimeString(),
    type: 'primary'
  }
])

// 生命周期钩子
onMounted(async () => {
  // 加载课程数据
  if (authStore.user?.role === 'student') {
    await courseStore.fetchMyCourses()
    courseCount.value = courseStore.myCourses.length
  } else {
    await courseStore.fetchCourses()
    courseCount.value = courseStore.courses.length
  }

  // 模拟其他数据
  weeklyActivities.value = Math.floor(Math.random() * 10) + 1
  todoCount.value = Math.floor(Math.random() * 5)
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.welcome-card {
  margin-bottom: 20px;
}

.welcome-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-header h3 {
  margin: 0;
}

.statistics {
  margin-bottom: 20px;
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-content {
  text-align: center;
}

.stat-content h2 {
  margin: 10px 0;
  font-size: 28px;
  color: #409EFF;
}

.stat-content p {
  margin: 0;
  color: #909399;
}

.recent-activities {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
}
</style> 