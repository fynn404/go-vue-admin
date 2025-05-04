<!--
  LoginView.vue
  登录页面组件
  
  功能：
  - 用户登录表单
  - 表单验证
  - 登录状态管理
  - 路由导航
-->
<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h2 class="card-header">用户登录</h2>
      </template>
      
      <el-form
        ref="loginForm"
        :model="loginData"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="loginData.username"
            placeholder="请输入用户名"
            prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginData.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>

        <div class="form-footer">
          <router-link to="/register">注册账号</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { rules } from '@/utils/validate'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'

// 路由实例
const router = useRouter()
// 认证状态管理
const authStore = useAuthStore()
// 表单引用
const loginForm = ref(null)
// 加载状态
const loading = ref(false)

// 登录表单数据
const loginData = reactive({
  username: '',
  password: ''
})

// 处理登录
const handleLogin = async () => {
  // 表单验证
  await loginForm.value.validate()
  
  loading.value = true
  try {
    // 调用登录接口
    const success = await authStore.login(loginData)
    if (success) {
      ElMessage.success('登录成功')
      router.push('/')
    }
  } catch (error) {
    ElMessage.error('登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
}

.login-card {
  width: 400px;
}

.card-header {
  text-align: center;
  margin: 0;
  color: #303133;
}

.login-button {
  width: 100%;
}

.form-footer {
  text-align: center;
  margin-top: 20px;
}

.form-footer a {
  color: #409EFF;
  text-decoration: none;
}

.form-footer a:hover {
  color: #66b1ff;
}
</style> 