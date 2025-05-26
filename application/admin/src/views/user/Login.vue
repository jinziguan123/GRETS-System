<template>
  <div class="login-container">
    <h2 class="login-title">用户登录</h2>
    

    <!-- 传统登录表单 -->
    <el-form
      v-if="loginMode === 'traditional'"
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
    </el-form>

      
    <!-- 注册链接 -->
    <el-form-item class="register-link">
      <span>还没有账号？</span>
      <el-button type="text" @click="router.push('register')">立即注册</el-button>
    </el-form-item>

    <!-- 密钥导入对话框 -->
    <el-dialog v-model="showKeyImportDialog" title="导入DID密钥" width="500px">
      <el-form :model="keyImportForm" :rules="keyImportRules" ref="keyImportFormRef">
        <el-form-item label="私钥" prop="privateKey">
          <el-input
            v-model="keyImportForm.privateKey"
            type="textarea"
            :rows="3"
            placeholder="请输入私钥"
            show-password
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="keyImportForm.password"
            type="password"
            placeholder="设置密钥保护密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="keyImportForm.confirmPassword"
            type="password"
            placeholder="确认密钥保护密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showKeyImportDialog = false">取消</el-button>
          <el-button type="primary" @click="importKey">导入</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 密钥信息对话框 -->
    <el-dialog v-model="showKeyInfoDialog" title="密钥信息" width="500px">
      <div class="key-info">
        <p><strong>公钥:</strong></p>
        <el-input v-model="currentKeyInfo.publicKey" readonly />
        <p style="margin-top: 15px;"><strong>状态:</strong> 已加载</p>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showKeyInfoDialog = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {useRouter} from 'vue-router'
import {useUserStore} from '@/stores/user.ts'
import type {FormInstance, FormRules} from 'element-plus'
import {ElMessage, ElMessageBox} from 'element-plus'
import {login} from '@/api/user'
import {getChallenge, didLogin, getDIDByUser} from '@/api/did'
import type {LoginData} from '@/types/api'
import {
  generateKeyPair,
  saveKeyPair,
  loadKeyPair,
  hasKeyPair,
  removeKeyPair,
  createAuthResponse,
  validateDID,
  generateHash
} from '@/utils/did'

interface LoginResponse {
  code: number
  message: string
  data: {
    token: string
    user?: any
  }
  token: string
  user?: any
}

const router = useRouter()
const formRef = ref<FormInstance | null>(null)
const didFormRef = ref<FormInstance | null>(null)
const keyImportFormRef = ref<FormInstance | null>(null)
const loading = ref<boolean>(false)
const didLoading = ref<boolean>(false)
const userStore = useUserStore()

// 登录模式
const loginMode = ref<'traditional' | 'did'>('traditional')

// 传统登录表单数据
interface LoginForm extends LoginData {
  remember: boolean
}

const loginForm = reactive<LoginForm>({
  citizenID: '',
  password: '',
  organization: '',
  remember: false
})

// DID登录表单数据
const didLoginForm = reactive({
  did: '',
  organization: ''
})

// 密钥导入表单
const keyImportForm = reactive({
  privateKey: '',
  password: '',
  confirmPassword: ''
})

// 对话框状态
const showKeyImportDialog = ref(false)
const showKeyInfoDialog = ref(false)

// 密钥状态
const hasStoredKey = ref(hasKeyPair())
const currentKeyInfo = reactive({
  publicKey: ''
})

// 可选组织列表
interface Organization {
  label: string
  value: string
}

const organizations: Organization[] = [
  { label: '政府监管部门', value: 'government' },
  { label: '银行机构', value: 'bank' },
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

const didLoginRules = reactive<FormRules>({
  did: [
    { required: true, message: '请输入DID或身份证号', trigger: 'blur' }
  ],
  organization: [
    { required: true, message: '请选择组织', trigger: 'change' }
  ]
})

const keyImportRules = reactive<FormRules>({
  privateKey: [
    { required: true, message: '请输入私钥', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请设置密钥保护密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== keyImportForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
})

if (localStorage.getItem('token') !== null) {
  router.push({ path: '/' })
}

// 处理传统登录
const handleLogin = (): void => {
  if (formRef.value) {
    formRef.value.validate(async (valid: boolean) => {
      if (valid) {
        loading.value = true
        try {
          const res = await login({
            citizenID: loginForm.citizenID,
            password: loginForm.password,
            organization: loginForm.organization
          }) as unknown as LoginResponse

          ElMessage.success('登录成功')
          localStorage.setItem('token', res.token)
          if (res.user) {
            userStore.updateUserInfo(res.user)
          }
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

// 处理DID登录
const handleDIDLogin = async (): Promise<void> => {
  if (didFormRef.value) {
    didFormRef.value.validate(async (valid: boolean) => {
      if (valid) {
        didLoading.value = true
        try {
          let userDID = didLoginForm.did

          // 如果输入的不是DID格式，尝试根据身份证号查找DID
          if (!validateDID(userDID)) {
            if (userDID.length >= 8) { // 管理员身份证号可能不是18位
              try {
                const didResponse = await getDIDByUser(userDID, didLoginForm.organization)
                userDID = didResponse.data.did
                didLoginForm.did = userDID
              } catch (error) {
                ElMessage.error('未找到对应的DID，请先注册')
                return
              }
            } else {
              ElMessage.error('请输入有效的DID或身份证号')
              return
            }
          }

          // 获取认证挑战
          const challengeResponse = await getChallenge({
            did: userDID,
            domain: window.location.hostname
          })

          const challenge = challengeResponse.data

          // 获取密钥对
          const keyPair = loadKeyPair(await getKeyPassword())
          if (!keyPair) {
            ElMessage.error('无法加载密钥，请重新导入')
            return
          }

          // 创建认证响应
          const authResponse = createAuthResponse(
            userDID,
            challenge,
            keyPair.privateKey,
            keyPair.publicKey
          )

          // 执行DID登录
          const loginResponse = await didLogin(authResponse)

          ElMessage.success('DID登录成功')
          localStorage.setItem('token', loginResponse.data.token)
          if (loginResponse.data.user) {
            userStore.updateUserInfo(loginResponse.data.user)
          }
          const redirect = router.currentRoute.value.query.redirect as string || '/'
          router.push(redirect)
        } catch (error) {
          console.error('DID登录失败:', error)
          ElMessage.error('DID登录失败，请检查网络连接或密钥是否正确')
        } finally {
          didLoading.value = false
        }
      }
    })
  }
}

// 导入密钥
const importKey = (): void => {
  if (keyImportFormRef.value) {
    keyImportFormRef.value.validate((valid: boolean) => {
      if (valid) {
        try {
          // 生成公钥（简化处理）
          const publicKey = generateHash(keyImportForm.privateKey)
          
          const keyPair = {
            privateKey: keyImportForm.privateKey,
            publicKey: publicKey
          }

          saveKeyPair(keyPair, keyImportForm.password)
          hasStoredKey.value = true
          currentKeyInfo.publicKey = publicKey
          
          ElMessage.success('密钥导入成功')
          showKeyImportDialog.value = false
          
          // 清空表单
          keyImportForm.privateKey = ''
          keyImportForm.password = ''
          keyImportForm.confirmPassword = ''
        } catch (error) {
          ElMessage.error('密钥导入失败')
        }
      }
    })
  }
}

// 移除存储的密钥
const removeStoredKey = (): void => {
  removeKeyPair()
  hasStoredKey.value = false
  currentKeyInfo.publicKey = ''
  ElMessage.success('密钥已清除')
}

// 获取密钥密码（简化处理，实际应该有更安全的方式）
const getKeyPassword = async (): Promise<string> => {
  return new Promise((resolve) => {
    ElMessageBox.prompt('请输入密钥保护密码', '密钥验证', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputType: 'password'
    }).then(({ value }) => {
      resolve(value)
    }).catch(() => {
      resolve('')
    })
  })
}

// 初始化时检查密钥状态
onMounted(() => {
  window.addEventListener('keydown', handleEnter)
  
  if (hasStoredKey.value) {
    // 尝试加载公钥信息用于显示
    try {
      const stored = localStorage.getItem('did_keypair')
      if (stored) {
        // 这里只是为了显示，不解密私钥
        currentKeyInfo.publicKey = '已加载（需要密码验证才能查看完整信息）'
      }
    } catch (error) {
      console.error('加载密钥信息失败:', error)
    }
  }
})

const handleEnter = (e: KeyboardEvent) => {
  if (e.keyCode === 13 || e.keyCode === 108) {
    if (loginMode.value === 'traditional') {
      handleLogin()
    } else {
      handleDIDLogin()
    }
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

.login-mode-selector {
  margin-bottom: 20px;
  text-align: center;
}

.mode-group {
  display: inline-flex;
}

.login-form {
  width: 100%;
}

.key-management {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.form-help {
  margin-top: 5px;
}

.form-help small {
  color: #999;
}

.key-info {
  padding: 10px 0;
}

.key-info p {
  margin: 10px 0 5px 0;
  font-weight: bold;
}

.register-link {
  text-align: center;
  margin-top: 20px;
}

.w-100 {
  width: 100%;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 