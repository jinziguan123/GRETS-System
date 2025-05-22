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

      <!-- 用户名输入框 -->
      <el-form-item prop="citizenID">
        <el-input
          v-model="loginForm.citizenID"
          v-if="loginForm.organization"
          placeholder="请输入身份证号"
          prefix-icon="User"
          autocomplete="citizenID"
        />
      </el-form-item>
      
      <!-- 密码输入框 -->
      <el-form-item prop="password">
        <el-input
          v-model="loginForm.password"
          v-if="loginForm.organization"
          type="password"
          placeholder="请输入密码"
          prefix-icon="Lock"
          autocomplete="current-password"
          show-password
        />
      </el-form-item>
      
      <!-- 记住我选项 -->
      <el-form-item>
        <div class="login-options">
<!--          <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>-->
<!--          <el-button type="text" @click="forgotPassword">忘记密码？</el-button>-->
        </div>
      </el-form-item>
      
      <!-- 登录按钮 -->
      <el-form-item>
        <el-button
          type="primary"
          :loading="loading"
          class="w-100"
          @click="handleLogin"
          @keydown.enter="handleEnter"
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

<script setup lang="ts">
import {useRouter} from 'vue-router'
import {useUserStore} from '@/stores/user.ts'
import type {FormInstance, FormRules} from 'element-plus'
import {ElMessage} from 'element-plus'
import {login} from '@/api/user'
import type {LoginData} from '@/types/api'

interface LoginResponse {
  code: number
  message: string
  data: {
    token: string
    user?: any
  }
}

const router = useRouter()
const formRef = ref<FormInstance | null>(null)
const loading = ref<boolean>(false)
const userStore = useUserStore()

// 登录表单数据
interface LoginForm extends LoginData {
  remember: boolean
}

const loginForm = reactive<LoginForm>({
  citizenID: '',
  password: '',
  organization: '',
  remember: false
})

// 可选组织列表
interface Organization {
  label: string
  value: string
}

const organizations: Organization[] = [
  { label: '政府监管部门', value: 'government' },
  { label: '投资者/买家', value: 'investor' },
  { label: '银行机构', value: 'bank' },
  // { label: '第三方机构', value: 'thirdparty' },
  { label: '审计机构', value: 'audit' }
]

// 表单验证规则
const loginRules = reactive<FormRules>({
  citizenID: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { 
      validator: (rule, value, callback) => {
        if (loginForm.organization === 'investor') {
          if (value.length !== 18) {
            callback(new Error('投资者身份证号长度应为18个字符'));
          } else {
            callback();
          }
        } else {
          // 非投资者组织只验证非空
          callback();
        }
      }, 
      trigger: 'blur' 
    }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ],
  organization: [
    { required: true, message: '请选择组织', trigger: 'change' }
  ]
})

// watch(() => loginForm.organization, (newVal) => {
//   if (newVal !== 'investor') {
//     // 对非投资者组织，设置固定身份证号
//     if (newVal) {
//       const orgName = newVal.charAt(0).toUpperCase() + newVal.slice(1);
//       loginForm.citizenID = `${orgName}Default`;
//     }
//   }else{
//     loginForm.citizenID = ''
//   }
// })

if (localStorage.getItem('token') !== null) {
  router.push({ path: '/' })
}

// 处理登录
const handleLogin = (): void => {
  if (formRef.value) {
    formRef.value.validate(async (valid: boolean) => {
      if (valid) {
        loading.value = true
        try {
          // 调用登录接口
          const res = await login({
            citizenID: loginForm.citizenID,
            password: loginForm.password,
            organization: loginForm.organization
          }) as unknown as LoginResponse

          ElMessage.success('登录成功')
          // 保存token
          localStorage.setItem('token', res.token)
          // 保存用户信息
          if (res.user) {
            userStore.updateUserInfo(res.user)
          }
          // 重定向到仪表盘或之前访问的页面
          const redirect = router.currentRoute.value.query.redirect as string || '/'
          router.push(redirect)
        } catch (error) {
          console.error('登录失败:', error)
          ElMessage.error('登录失败，请检查网络连接')
        } finally {
          loading.value = false
        }
      }
    })
  }
}

// 忘记密码处理
const forgotPassword = (): void => {
  router.push('/forgot-password')
}

onMounted(() => {
  window.addEventListener('keydown', handleEnter);
});

const handleEnter = (e) => {
  if (e.keyCode === 13 || e.keyCode === 108) {
    handleLogin()
  }
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