<template>
  <el-container class="app-wrapper">
    <el-header v-if="isAuthenticated">
      <nav-header />
    </el-header>
    
    <el-container>
      <el-aside width="200px" v-if="isAuthenticated">
        <nav-menu />
      </el-aside>
      
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import NavHeader from '@/components/NavHeader.vue'
import NavMenu from '@/components/NavMenu.vue'

const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
</script>

<style>
.app-wrapper {
  min-height: 100vh;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0;
}

.el-aside {
  background-color: #304156;
  color: #fff;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style> 