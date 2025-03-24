<template>
  <div class="login-container">
    <h2 class="login-title">用户登录</h2>
    
    <el-form
      ref="formRef"
      :model="loginForm"
      :rules="loginRules"
      @submit.prevent="handleLogin"
      class="login-form"
    >
      <!-- 用户名输入框 -->
      <el-form-item prop="citizenID">
        <el-input
          v-model="loginForm.citizenID"
          placeholder="请输入身份证号"
          prefix-icon="User"
          autocomplete="citizenID"
        />
      </el-form-item>
      
      <!-- 密码输入框 -->
      <el-form-item prop="password">
        <el-input
          v-model="loginForm.password"
          type="password"
          placeholder="请输入密码"
          prefix-icon="Lock"
          autocomplete="current-password"
          show-password
        />
      </el-form-item>
      
      <!-- 组织选择 -->
      <el-form-item prop="organization">
        <el-select
          v-model="loginForm.organization"
          placeholder="请选择组织"
          class="w-100"
        >
          <el-option
            v-for="org in organizations"
            :key="org.value"
            :label="org.label"
            :value="org.value"
          />
        </el-select>
      </el-form-item>
      
      <!-- 记住我选项 -->
      <el-form-item>
        <div class="login-options">
          <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>
          <el-button type="text" @click="forgotPassword">忘记密码？</el-button>
        </div>
      </el-form-item>
      
      <!-- 登录按钮 -->
      <el-form-item>
        <el-button
          type="primary"
          :loading="loading"
          class="w-100"
          @click="handleLogin"
        >
          {{ loading ? '登录中...' : '登录' }}
        </el-button>
      </el-form-item>
      
      <!-- 注册链接 -->
      <el-form-item class="register-link">
        <span>还没有账号？</span>
        <el-button type="text" @click="router.push('register')">立即注册</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { login } from './api'
import {useLocalStorage} from "@vueuse/core";

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)

// 登录表单数据
const loginForm = reactive({
  citizenID: '',
  password: '',
  organization: '',
  remember: false
})

// 可选组织列表
const organizations = [
  { label: '系统管理员', value: 'administrator' },
  { label: '政府监管部门', value: 'government' },
  { label: '房产中介机构', value: 'agency' },
  { label: '投资者/买家', value: 'investor' },
  { label: '银行机构', value: 'bank' }
]

// 表单验证规则
const loginRules = reactive({
  citizenID: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { min: 18, max: 18, message: '身份证号长度应为18个字符', trigger: 'blur' }
  ],
  organization: [
    { required: true, message: '请选择组织', trigger: 'change' }
  ]
})

if (localStorage.getItem('token') !== null) {
  router.push({ path: '/' })
}

// 处理登录
const handleLogin = () => {
  if (formRef.value) {
    formRef.value.validate(async (valid) => {
      if (valid) {
        loading.value = true
        try {
          // 调用登录接口
          await login({
            citizenID: loginForm.citizenID,
            password: loginForm.password,
            organization: loginForm.organization,
            remember: loginForm.remember
          }).then(res => {
            if (res.code === 200) {
              ElMessage.success('登录成功')
              // 保存token
              localStorage.setItem('token', res.data.token)
              // 保存用户信息
              localStorage.setItem('userInfo', JSON.stringify(res.data.user))
              // 重定向到仪表盘或之前访问的页面
              const redirect = router.currentRoute.value.query.redirect || '/'
              router.push(redirect)
            }else{
              ElMessage.error('用户不存在，请先注册')
            }
          })
        } finally {
          loading.value = false
        }
      }
    })
  }
}

// 忘记密码处理
const forgotPassword = () => {
  router.push('/forgot-password')
}
</script>

<style scoped>
.login-container {
  width: 100%;
}

.login-title {
  font-size: 24px;
  color: #333;
  text-align: center;
  margin-bottom: 30px;
}

.login-form {
  width: 100%;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.w-100 {
  width: 100%;
}

.register-link {
  text-align: center;
  margin-top: 10px;
}

.register-link span {
  font-size: 14px;
  color: #666;
}
</style> 