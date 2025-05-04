<!--
  NavMenu.vue
  侧边栏导航菜单组件
  
  功能：
  - 根据用户角色显示不同的导航选项
  - 支持路由导航
  - 显示当前激活的菜单项
  - 使用Element Plus的图标和菜单组件
-->
<template>
  <!-- 
    导航菜单容器
    :default-active - 当前激活的菜单项（与路由路径对应）
    background-color - 菜单背景色
    text-color - 文字颜色
    active-text-color - 激活项文字颜色
    router - 启用路由模式（点击菜单项自动导航）
  -->
  <el-menu
    :default-active="activeMenu"
    class="nav-menu"
    background-color="#304156"
    text-color="#bfcbd9"
    active-text-color="#409EFF"
    router
  >
    <!-- 仪表盘 - 所有用户可见 -->
    <el-menu-item index="/dashboard">
      <el-icon><Odometer /></el-icon>
      <span>仪表盘</span>
    </el-menu-item>

    <!-- 学生专属菜单项 -->
    <el-menu-item index="/courses" v-if="isStudent">
      <el-icon><Reading /></el-icon>
      <span>选课中心</span>
    </el-menu-item>

    <el-menu-item index="/my-courses" v-if="isStudent">
      <el-icon><Collection /></el-icon>
      <span>我的课程</span>
    </el-menu-item>

    <!-- 教师专属菜单项 -->
    <el-menu-item index="/teaching" v-if="isTeacher">
      <el-icon><Edit /></el-icon>
      <span>课程管理</span>
    </el-menu-item>

    <el-menu-item index="/students" v-if="isTeacher">
      <el-icon><User /></el-icon>
      <span>学生管理</span>
    </el-menu-item>

    <!-- 管理员专属菜单项 -->
    <el-menu-item index="/users" v-if="isAdmin">
      <el-icon><UserFilled /></el-icon>
      <span>用户管理</span>
    </el-menu-item>

    <el-menu-item index="/settings" v-if="isAdmin">
      <el-icon><Setting /></el-icon>
      <span>系统设置</span>
    </el-menu-item>
  </el-menu>
</template>

<script setup>
// 导入所需的组件和工具
import { computed } from 'vue'                    // Vue组合式API
import { useRoute } from 'vue-router'             // Vue路由
import { useAuthStore } from '@/stores/auth'      // 认证状态管理
import {
  Odometer,    // 仪表盘图标
  Reading,     // 阅读图标
  Collection,  // 收藏图标
  Edit,        // 编辑图标
  User,        // 用户图标
  UserFilled,  // 用户填充图标
  Setting,     // 设置图标
} from '@element-plus/icons-vue'                  // Element Plus图标

// 获取路由和认证状态实例
const route = useRoute()
const authStore = useAuthStore()

// 计算属性：当前激活的菜单项
const activeMenu = computed(() => route.path)

// 计算属性：用户角色判断
const isStudent = computed(() => authStore.user?.role === 'student')  // 是否为学生
const isTeacher = computed(() => authStore.user?.role === 'teacher')  // 是否为教师
const isAdmin = computed(() => authStore.user?.role === 'admin')      // 是否为管理员
</script>

<style scoped>
/* 导航菜单样式 */
.nav-menu {
  height: 100%;          /* 占满容器高度 */
  border-right: none;    /* 移除右边框 */
}

/* 菜单项样式 */
.nav-menu .el-menu-item {
  display: flex;         /* 弹性布局 */
  align-items: center;   /* 垂直居中对齐 */
}

/* 图标样式 */
.nav-menu .el-icon {
  margin-right: 12px;    /* 图标右侧间距 */
  font-size: 18px;       /* 图标大小 */
}
</style> 