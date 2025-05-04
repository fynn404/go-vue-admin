<!--
  NavHeader.vue
  顶部导航栏组件
  
  功能：
  - 显示系统标题
  - 显示当前用户名
  - 提供用户菜单（个人信息、退出登录）
  - 响应式布局设计
-->
<template>
  <!-- 导航栏容器 -->
  <div class="nav-header">
    <!-- 左侧区域：系统标题 -->
    <div class="left">
      <h2>课程管理系统</h2>
    </div>
    <!-- 右侧区域：用户菜单 -->
    <div class="right">
      <!-- 
        下拉菜单组件
        @command - 菜单项命令处理函数
      -->
      <el-dropdown @command="handleCommand">
        <!-- 触发下拉菜单的元素 -->
        <span class="el-dropdown-link">
          {{ username }}
          <el-icon class="el-icon--right"><arrow-down /></el-icon>
        </span>
        <!-- 下拉菜单内容 -->
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">个人信息</el-dropdown-item>
            <el-dropdown-item command="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
// 导入所需的组件和工具
import { computed } from 'vue'                    // Vue组合式API
import { useRouter } from 'vue-router'            // Vue路由
import { useAuthStore } from '@/stores/auth'      // 认证状态管理
import { ArrowDown } from '@element-plus/icons-vue' // 下拉箭头图标

// 获取路由和认证状态实例
const router = useRouter()
const authStore = useAuthStore()

// 计算属性：当前用户名（如果未登录则显示空字符串）
const username = computed(() => authStore.user?.username || '')

/**
 * 处理下拉菜单命令
 * @param {string} command - 菜单命令（profile/logout）
 */
const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')  // 导航到个人信息页面
      break
    case 'logout':
      authStore.logout()       // 执行登出操作
      router.push('/login')    // 导航到登录页面
      break
  }
}
</script>

<style scoped>
/* 导航栏容器样式 */
.nav-header {
  height: 60px;               /* 固定高度 */
  display: flex;              /* 弹性布局 */
  align-items: center;        /* 垂直居中 */
  justify-content: space-between; /* 两端对齐 */
  padding: 0 20px;           /* 水平内边距 */
}

/* 系统标题样式 */
.left h2 {
  margin: 0;                 /* 移除默认外边距 */
  color: #409EFF;           /* 主题蓝色 */
}

/* 下拉菜单触发元素样式 */
.el-dropdown-link {
  cursor: pointer;           /* 鼠标指针样式 */
  display: flex;             /* 弹性布局 */
  align-items: center;       /* 垂直居中 */
  color: #606266;           /* 文字颜色 */
}

/* 下拉菜单触发元素悬停样式 */
.el-dropdown-link:hover {
  color: #409EFF;           /* 主题蓝色 */
}
</style> 