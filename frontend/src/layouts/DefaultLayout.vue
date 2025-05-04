<!--
  DefaultLayout.vue
  默认布局组件
  
  功能：
  - 整体页面布局
  - 导航菜单
  - 头部导航
  - 内容区域
-->
<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside width="220px" class="aside">
      <div class="logo">
        <h1>课程管理系统</h1>
      </div>
      <nav-menu />
    </el-aside>

    <el-container>
      <!-- 头部 -->
      <el-header height="60px" class="header">
        <nav-header />
      </el-header>

      <!-- 主要内容区 -->
      <el-main class="main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import NavMenu from '@/components/NavMenu.vue'
import NavHeader from '@/components/NavHeader.vue'
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside {
  background-color: #304156;
  color: #fff;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #1f2d3d;
}

.logo h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.main {
  background-color: #f5f7fa;
  padding: 20px;
}

/* 路由过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style> 