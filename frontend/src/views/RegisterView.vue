<!--
  RegisterView.vue
  注册页面组件
  
  功能：
  - 用户注册表单
  - 表单验证
  - 角色选择
  - 注册成功后自动登录
-->
<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <h2 class="card-header">用户注册</h2>
      </template>
      
      <el-form
        ref="registerForm"
        :model="registerData"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="registerData.username"
            placeholder="请输入用户名"
            prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="registerData.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="registerData.confirmPassword"
            type="password"
            placeholder="请确认密码"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item label="角色" prop="role">
          <el-select v-model="registerData.role" placeholder="请选择角色" class="role-select">
            <el-option label="学生" value="student" />
            <el-option label="教师" value="teacher" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            class="register-button"
            @click="handleRegister"
          >
            注册
          </el-button>
        </el-form-item>

        <div class="form-footer">
          <router-link to="/login">返回登录</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { rules as baseRules } from '@/utils/validate'
import { register } from '@/api/auth'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'

// 路由实例
const router = useRouter()
// 认证状态管理
const authStore = useAuthStore()
// 表单引用
const registerForm = ref(null)
// 加载状态
const loading = ref(false)

// 注册表单数据
const registerData = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  role: ''
})

// 扩展验证规则
const rules = {
  ...baseRules,
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== registerData.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 处理注册
const handleRegister = async () => {
  // 表单验证
  await registerForm.value.validate()
  
  loading.value = true
  try {
    // 调用注册接口
    const response = await register({
      username: registerData.username,
      password: registerData.password,
      role: registerData.role
    })
    
    // 注册成功后自动登录
    await authStore.login({
      username: registerData.username,
      password: registerData.password
    })
    
    ElMessage.success('注册成功')
    router.push('/')
  } catch (error) {
    ElMessage.error('注册失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
}

.register-card {
  width: 400px;
}

.card-header {
  text-align: center;
  margin: 0;
  color: #303133;
}

.role-select {
  width: 100%;
}

.register-button {
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