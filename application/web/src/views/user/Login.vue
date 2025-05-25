<template>
  <div class="login-container">
    <h2 class="login-title">投资者登录</h2>
    
    <!-- 登录方式选择 -->
    <div class="login-mode-selector">
      <el-radio-group v-model="loginMode" class="mode-group">
        <el-radio-button label="traditional">传统登录</el-radio-button>
        <el-radio-button label="did">DID登录</el-radio-button>
      </el-radio-group>
      <div class="mode-description">
        <small v-if="loginMode === 'traditional'">
          使用身份证号和密码进行传统登录
        </small>
        <small v-if="loginMode === 'did'">
          使用去中心化身份标识(DID)和数字签名进行安全登录
        </small>
      </div>
    </div>

    <!-- 传统登录表单 -->
    <el-form
      v-if="loginMode === 'traditional'"
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

    <!-- DID登录表单 -->
    <el-form
      v-if="loginMode === 'did'"
      ref="didFormRef"
      :model="didLoginForm"
      :rules="didLoginRules"
      @submit.prevent="handleDIDLogin"
      class="login-form"
    >
      <!-- DID输入框 -->
      <el-form-item prop="did">
        <el-input
          v-model="didLoginForm.did"
          placeholder="请输入DID或身份证号"
          prefix-icon="Key"
        />
        <div class="form-help">
          <small>可以输入完整的DID或身份证号自动查找DID</small>
        </div>
      </el-form-item>

      <!-- DID登录进度指示器 -->
      <el-form-item v-if="didLoading">
        <div class="login-progress">
          <el-steps :active="loginStep" finish-status="success" simple>
            <el-step title="验证DID" />
            <el-step title="获取挑战" />
            <el-step title="生成签名" />
            <el-step title="身份验证" />
            <el-step title="登录成功" />
          </el-steps>
          <div class="progress-text">{{ loginProgressText }}</div>
        </div>
      </el-form-item>

      <!-- 密钥管理 -->
      <el-form-item>
        <div class="key-management">
          <el-button 
            v-if="!hasStoredKey" 
            type="info" 
            @click="showKeyImportDialog = true"
            icon="Upload"
          >
            导入密钥
          </el-button>
          <el-button 
            v-if="hasStoredKey" 
            type="success" 
            @click="showKeyInfoDialog = true"
            icon="Key"
          >
            密钥已加载
          </el-button>
          <el-button 
            v-if="hasStoredKey" 
            type="danger" 
            @click="removeStoredKey"
            icon="Delete"
          >
            清除密钥
          </el-button>
        </div>
      </el-form-item>
      
      <!-- DID登录按钮 -->
      <el-form-item>
        <el-button
          type="primary"
          :loading="didLoading"
          :disabled="!hasStoredKey"
          class="w-100"
          @click="handleDIDLogin"
        >
          {{ didLoading ? loginProgressText : 'DID登录' }}
        </el-button>
      </el-form-item>
    </el-form>
      
    <!-- 注册链接 -->
    <el-form-item class="register-link">
      <span>还没有账号？</span>
      <el-button type="text" @click="router.push('/auth/register')">立即注册</el-button>
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
  generateHash,
  generatePublicKeyFromPrivate
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
  organization: 'investor',
  remember: false
})

// DID登录表单数据
const didLoginForm = reactive({
  did: ''
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

// DID登录进度状态
const loginStep = ref(0)
const loginProgressText = ref('DID认证中...')

// 表单验证规则
const loginRules = reactive<FormRules>({
  citizenID: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { 
      min: 18, 
      max: 18, 
      message: '身份证号应为18位', 
      trigger: 'blur' 
    }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
})

const didLoginRules = reactive<FormRules>({
  did: [
    { required: true, message: '请输入DID或身份证号', trigger: 'blur' }
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

// 处理DID登录 - 按照挑战-响应机制实现
const handleDIDLogin = async (): Promise<void> => {
  if (didFormRef.value) {
    didFormRef.value.validate(async (valid: boolean) => {
      if (valid) {
        didLoading.value = true
        loginStep.value = 0
        
        try {
          let userDID = didLoginForm.did

          // 步骤1: 验证DID格式或根据身份证号查找DID
          loginStep.value = 1
          loginProgressText.value = '正在验证DID...'
          
          if (!validateDID(userDID)) {
            if (userDID.length === 18) {
              try {
                const didResponse = await getDIDByUser(userDID, 'investor')
                userDID = didResponse.did
                didLoginForm.did = userDID
                ElMessage.info(`找到DID: ${userDID}`)
              } catch (error) {
                ElMessage.error('未找到对应的DID，请先注册')
                return
              }
            } else {
              ElMessage.error('请输入有效的DID或18位身份证号')
              return
            }
          }

          // 步骤2: 向RP请求认证挑战
          loginStep.value = 2
          loginProgressText.value = '正在获取认证挑战...'
          
          const challengeResponse = await getChallenge({
            did: userDID,
            domain: window.location.hostname
          })

          const challenge = challengeResponse
          ElMessage.success('获取挑战成功')

          // 步骤3: 获取用户密钥对
          loginProgressText.value = '正在加载密钥...'
          
          const password = await getKeyPassword()          
          if (!password) {
            ElMessage.error('需要密钥密码才能继续')
            return
          }

          const keyPair = loadKeyPair(password)
          if (!keyPair) {
            ElMessage.error('无法加载密钥，请检查密码或重新导入密钥')
            return
          }

          // 步骤4: 使用私钥对挑战进行签名
          loginStep.value = 3
          loginProgressText.value = '正在生成数字签名...'
          
          const authResponse = await createAuthResponse(
            userDID,
            challenge,
            keyPair.privateKey,
            keyPair.publicKey
          )

          // 步骤5: 将签名和公钥信息发送给RP进行验证
          loginStep.value = 4
          loginProgressText.value = '正在验证身份...'
          
          const loginResponse = await didLogin(authResponse)

          // 步骤6: 登录成功处理
          loginStep.value = 5
          loginProgressText.value = '登录成功！'
          
          ElMessage.success('DID身份验证成功！')
          localStorage.setItem('token', loginResponse.data.token)
          
          // 处理用户信息，确保类型兼容
          if (loginResponse.data.user) {
            const userInfo = {
              ...loginResponse.data.user,
              role: loginResponse.data.user.role === 'investor' ? 'user' as const : 'admin' as const
            }
            userStore.updateUserInfo(userInfo)
          }
          
          // 延迟跳转，让用户看到成功状态
          setTimeout(() => {
            const redirect = router.currentRoute.value.query.redirect as string || '/'
            router.push(redirect)
          }, 1000)
          
        } catch (error: any) {
          console.error('DID登录失败:', error)
          
          // 根据错误类型提供更具体的错误信息
          if (error.message?.includes('密码错误')) {
            ElMessage.error('密钥密码错误，请重新输入')
          } else if (error.message?.includes('签名')) {
            ElMessage.error('数字签名验证失败，请检查密钥是否正确')
          } else if (error.message?.includes('挑战')) {
            ElMessage.error('认证挑战已过期，请重试')
          } else {
            ElMessage.error('DID登录失败，请检查网络连接或联系管理员')
          }
        } finally {
          // 延迟重置状态
          setTimeout(() => {
            didLoading.value = false
            loginStep.value = 0
            loginProgressText.value = 'DID认证中...'
          }, loginStep.value === 5 ? 1500 : 500)
        }
      }
    })
  }
}

// 导入密钥
const importKey = async (): Promise<void> => {
  if (keyImportFormRef.value) {
    keyImportFormRef.value.validate(async (valid: boolean) => {
      if (valid) {
        try {
          // 从私钥生成对应的公钥，与后端crypto.go格式保持一致          
          const publicKey = await generatePublicKeyFromPrivate(keyImportForm.privateKey)
          if (!publicKey) {
            throw new Error('无法从私钥生成公钥')
          }
          
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
        } catch (error: any) {
          console.error('密钥导入失败:', error)
          ElMessage.error(`密钥导入失败: ${error.message}`)
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

.login-progress {
  margin: 20px 0;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.progress-text {
  text-align: center;
  margin-top: 10px;
  color: #666;
  font-size: 14px;
}

.mode-description {
  text-align: center;
  margin-top: 8px;
}

.mode-description small {
  color: #999;
  font-size: 12px;
}
</style> 