<template>
  <div class="register-container">
    <h2 class="register-title">投资者注册</h2>
    
    <!-- 注册方式选择 -->
    <div class="register-mode-selector">
      <el-radio-group v-model="registerMode" class="mode-group">
        <el-radio-button label="traditional">传统注册</el-radio-button>
        <el-radio-button label="did">DID注册</el-radio-button>
      </el-radio-group>
      <div class="mode-description">
        <p v-if="registerMode === 'traditional'">
          使用传统的用户名密码方式注册，适合熟悉传统系统的用户
        </p>
        <p v-if="registerMode === 'did'">
          使用去中心化身份(DID)注册，享受更安全的区块链身份认证
        </p>
      </div>
    </div>
    
    <el-form
      ref="formRef"
      :model="registerForm"
      :rules="registerRules"
      @submit.prevent="handleRegister"
      class="register-form"
    >

      <!-- 用户名 -->
      <el-form-item prop="name">
        <el-input
          v-model="registerForm.name"
          placeholder="请输入姓名"
          prefix-icon="User"
        />
      </el-form-item>
      
      <!-- 身份证号 -->
      <el-form-item prop="citizenID">
        <el-input
          v-model="registerForm.citizenID"
          placeholder="请输入身份证号"
          prefix-icon="User"
        />
      </el-form-item>

      <!-- 密码 -->
      <el-form-item prop="password">
        <el-input
          v-model="registerForm.password"
          type="password"
          placeholder="请输入密码"
          prefix-icon="Lock"
          show-password
        />
      </el-form-item>
      
      <!-- 确认密码 -->
      <el-form-item prop="confirmPassword">
        <el-input
          v-model="registerForm.confirmPassword"
          type="password"
          placeholder="请确认密码"
          prefix-icon="Lock"
          show-password
        />
      </el-form-item>
      
      <!-- 电子邮箱 -->
      <el-form-item prop="email">
        <el-input
          v-model="registerForm.email"
          placeholder="请输入电子邮箱"
          prefix-icon="Message"
        />
      </el-form-item>
      
      <!-- 手机号码 -->
      <el-form-item prop="phone">
        <el-input
          v-model="registerForm.phone"
          placeholder="请输入手机号码"
          prefix-icon="Phone"
        />
      </el-form-item>

      <!-- 注册金额 -->
      <el-form-item prop="balance">
        <el-input
          v-model="registerForm.balance"
          placeholder="请输入注册金额"
          prefix-icon="Money"
          type="number"
        />
      </el-form-item>

      <!-- DID密钥管理 -->
      <div v-if="registerMode === 'did'" class="did-key-section">
        <el-divider content-position="left">
          <span class="key-section-title">
            <el-icon><Key /></el-icon>
            DID密钥管理
          </span>
        </el-divider>
        
        <!-- 密钥生成 -->
        <el-form-item>
          <el-button 
            type="primary" 
            @click="generateKeys"
            :disabled="!!keyPair"
            icon="Key"
          >
            {{ keyPair ? '密钥已生成' : '生成密钥对' }}
          </el-button>
          <el-button 
            v-if="keyPair" 
            type="warning" 
            @click="regenerateKeys"
            icon="Refresh"
          >
            重新生成
          </el-button>
        </el-form-item>

        <!-- 密钥显示 -->
        <div v-if="keyPair" class="key-display">
          <el-alert
            title="重要提醒"
            type="warning"
            :closable="false"
            show-icon
          >
            <template #default>
              <p>请务必安全保存您的私钥！私钥丢失将无法恢复您的DID身份。</p>
              <p>建议将私钥保存到安全的地方，如密码管理器或离线存储设备。</p>
            </template>
          </el-alert>

          <el-form-item label="公钥">
            <el-input
              v-model="keyPair.publicKey"
              readonly
              type="textarea"
              :rows="3"
            >
              <template #append>
                <el-button @click="copyToClipboard(keyPair.publicKey)">复制</el-button>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="私钥">
            <el-input
              v-model="keyPair.privateKey"
              readonly
              type="textarea"
              :rows="3"
              :show-password="!showPrivateKey"
            >
              <template #append>
                <el-button @click="showPrivateKey = !showPrivateKey">
                  {{ showPrivateKey ? '隐藏' : '显示' }}
                </el-button>
                <el-button @click="copyToClipboard(keyPair.privateKey)">复制</el-button>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item>
            <el-checkbox v-model="keysSaved">
              我已安全保存私钥，理解私钥丢失的风险
            </el-checkbox>
          </el-form-item>
        </div>
      </div>
      
      <!-- 服务条款 -->
      <el-form-item prop="agreement">
        <el-checkbox v-model="registerForm.agreement">
          我已阅读并同意
          <el-button type="text" @click="showTerms">服务条款</el-button>
          和
          <el-button type="text" @click="showPrivacy">隐私政策</el-button>
        </el-checkbox>
      </el-form-item>
      
      <!-- 注册按钮 -->
      <el-form-item>
        <el-button
          type="primary"
          :loading="loading"
          class="w-100"
          @click="handleRegister"
          :disabled="registerMode === 'did' && (!keyPair || !keysSaved)"
        >
          {{ loading ? '注册中...' : (registerMode === 'did' ? 'DID注册' : '传统注册') }}
        </el-button>
      </el-form-item>
      
      <!-- 登录链接 -->
      <el-form-item class="login-link">
        <span>已有账号？</span>
        <el-button type="text" @click="router.push('/auth/login')">立即登录</el-button>
      </el-form-item>
    </el-form>
    
    <!-- 服务条款对话框 -->
    <el-dialog v-model="termsDialogVisible" title="服务条款" width="50%">
      <div class="terms-content">
        <h3>房地产区块链交易系统服务条款</h3>
        <p>欢迎使用我们的房地产区块链交易系统。本服务条款（"条款"）适用于您对本系统及相关服务的访问和使用。通过注册并使用本系统，您同意接受以下条款的约束。</p>
        
        <h4>1. 服务描述</h4>
        <p>本系统提供一个基于区块链技术的房地产交易平台，用户可通过该平台查看、发布房产信息，进行交易，管理合同、支付等相关服务。</p>
        
        <h4>2. 账号注册与安全</h4>
        <p>您须提供真实、准确、完整的个人信息进行注册。您有责任维护账号安全，保障登录凭证不被泄露，并对使用您账号进行的所有活动负责。</p>
        
        <h4>3. 用户义务</h4>
        <p>您同意：(a) 不进行任何违法或不当活动；(b) 不上传、发布虚假或误导性信息；(c) 不侵犯他人知识产权或隐私；(d) 不干扰系统正常运行。</p>
        
        <h4>4. 隐私保护</h4>
        <p>我们尊重并保护用户隐私，收集和使用个人信息的方式详见隐私政策。使用本服务即表示您同意我们按照隐私政策收集和使用您的相关信息。</p>
        
        <h4>5. 知识产权</h4>
        <p>本系统及相关内容（包括但不限于文本、图像、代码等）的知识产权归我们或相关权利人所有。未经授权，您不得复制、修改、分发或销售这些内容。</p>
        
        <h4>6. 免责声明</h4>
        <p>本系统按"现状"提供，我们不保证服务不会中断或无错误。我们不对因使用本服务而导致的任何损失承担责任，除非法律另有规定。</p>
        
        <h4>7. 条款修改</h4>
        <p>我们保留随时修改这些条款的权利。修改后的条款将在系统上发布，继续使用本服务即表示您接受修改后的条款。</p>
        
        <h4>8. 法律适用</h4>
        <p>本条款受中华人民共和国法律管辖，任何争议应提交至有管辖权的人民法院解决。</p>

        <!-- DID相关条款 -->
        <h4>9. DID身份管理</h4>
        <p>如果您选择使用DID注册：</p>
        <ul>
          <li>您的身份将基于去中心化身份(DID)技术进行管理</li>
          <li>您需要妥善保管您的私钥，私钥丢失将无法恢复您的身份</li>
          <li>我们不会存储您的私钥，无法帮助您恢复丢失的私钥</li>
          <li>您的DID身份信息将记录在区块链上，具有不可篡改性</li>
          <li>您可以使用DID身份在支持的平台间进行身份验证</li>
        </ul>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="termsDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="agreeTerms">同意</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 隐私政策对话框 -->
    <el-dialog v-model="privacyDialogVisible" title="隐私政策" width="50%">
      <div class="terms-content">
        <h3>房地产区块链交易系统隐私政策</h3>
        <p>本隐私政策描述了我们如何收集、使用、存储和保护您的个人信息。我们非常重视您的隐私保护，请仔细阅读本政策。</p>
        
        <h4>1. 信息收集</h4>
        <p>我们可能收集以下类型的信息：</p>
        <ul>
          <li>注册信息：用户名、密码、电子邮箱、手机号码、所属组织等</li>
          <li>交易信息：房产交易记录、合同信息、支付信息等</li>
          <li>使用数据：访问日志、设备信息、IP地址等</li>
          <li>DID信息：DID标识符、公钥、DID文档等（不包括私钥）</li>
        </ul>
        
        <h4>2. 信息使用</h4>
        <p>我们使用收集的信息用于：</p>
        <ul>
          <li>提供、维护和改进我们的服务</li>
          <li>处理交易和请求</li>
          <li>发送通知和更新</li>
          <li>预防欺诈和提升安全性</li>
          <li>分析使用趋势和改进用户体验</li>
          <li>验证DID身份和处理相关凭证</li>
        </ul>
        
        <h4>3. 信息共享</h4>
        <p>我们不会出售、出租或交易您的个人信息。但在以下情况下，我们可能会共享您的信息：</p>
        <ul>
          <li>经您明确同意</li>
          <li>完成交易所必需的共享</li>
          <li>遵守法律要求、保护权利或安全</li>
          <li>DID身份验证所需的公开信息</li>
        </ul>
        
        <h4>4. 区块链数据</h4>
        <p>请注意，存储在区块链上的交易数据具有不可篡改性。个人敏感信息不会直接存储在区块链上，但交易记录将被永久保存。</p>
        
        <h4>5. DID隐私保护</h4>
        <p>对于DID用户：</p>
        <ul>
          <li>我们不会存储您的私钥</li>
          <li>您的DID标识符和公钥可能会公开</li>
          <li>您可以控制哪些信息通过可验证凭证分享</li>
          <li>您的身份验证过程是去中心化的</li>
        </ul>
        
        <h4>6. 数据安全</h4>
        <p>我们采取合理措施保护您的个人信息，包括加密技术和访问控制。但互联网传输无法保证绝对安全，您使用服务即表示了解这一风险。</p>
        
        <h4>7. 您的权利</h4>
        <p>您有权访问、更正或删除您的个人信息。如需行使这些权利，请通过系统内的设置或联系我们的客服。</p>
        
        <h4>8. 政策变更</h4>
        <p>我们可能会不时更新本隐私政策。更新后的政策将在系统内发布，请定期查阅。</p>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="privacyDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="agreePrivacy">同意</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Key } from '@element-plus/icons-vue'
import { register } from '../../api/user.js'
import { didRegister } from '../../api/did.ts'
import { generateKeyPair } from '../../utils/did.ts'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const termsDialogVisible = ref(false)
const privacyDialogVisible = ref(false)

// 注册模式
const registerMode = ref('traditional')

// DID相关状态
const keyPair = ref(null)
const showPrivateKey = ref(false)
const keysSaved = ref(false)

// 注册表单数据
const registerForm = reactive({
  name: '',
  citizenID: '',
  password: '',
  confirmPassword: '',
  email: '',
  phone: '',
  organization: 'investor', // 固定为投资者
  balance: undefined,
  agreement: false
})

// 监听注册模式变化，重置相关状态
watch(registerMode, (newMode) => {
  if (newMode === 'traditional') {
    keyPair.value = null
    keysSaved.value = false
    showPrivateKey.value = false
  }
})

// 生成密钥对
const generateKeys = async () => {
  try {
    keyPair.value = await generateKeyPair()
    keysSaved.value = false
    showPrivateKey.value = false
    ElMessage.success('密钥对生成成功！请务必保存您的私钥')
    
    // 自动滚动到密钥显示区域
    setTimeout(() => {
      const keyDisplayElement = document.querySelector('.key-display')
      if (keyDisplayElement) {
        keyDisplayElement.scrollIntoView({ 
          behavior: 'smooth', 
          block: 'center' 
        })
      }
    }, 100)
  } catch (error) {
    ElMessage.error('密钥生成失败：' + error.message)
  }
}

// 重新生成密钥对
const regenerateKeys = () => {
  ElMessageBox.confirm(
    '重新生成密钥对将覆盖当前密钥，请确保已保存当前私钥。是否继续？',
    '确认重新生成',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    generateKeys()
  }).catch(() => {
    // 用户取消
  })
}

// 复制到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    ElMessage.success('已复制到剪贴板')
  }
}

// 验证密码是否一致
const validatePass = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

// 验证手机号码
const validatePhone = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入手机号码'))
  } else {
    const reg = /^1[3-9]\d{9}$/
    if (reg.test(value)) {
      callback()
    } else {
      callback(new Error('请输入有效的手机号码'))
    }
  }
}

// 验证邮箱
const validateEmail = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入电子邮箱'))
  } else {
    const reg = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
    if (reg.test(value)) {
      callback()
    } else {
      callback(new Error('请输入有效的电子邮箱'))
    }
  }
}

// 验证投资者余额
const validateBalance = (rule, value, callback) => {
  if (value === undefined || value === '') {
    callback(new Error('必须输入注册金额'))
  } else if (isNaN(value) || parseFloat(value) <= 0) {
    callback(new Error('注册金额必须大于0'))
  } else {
    callback()
  }
}

// 表单验证规则
const registerRules = reactive({
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 3, max: 20, message: '姓名长度应为3-20个字符', trigger: 'blur' }
  ],
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
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度应为6-20个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validatePass, trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入电子邮箱', trigger: 'blur' },
    { validator: validateEmail, trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    { validator: validatePhone, trigger: 'blur' }
  ],

  balance: [
    { validator: validateBalance, trigger: 'blur' }
  ],
  agreement: [
    { 
      required: true, 
      validator: (rule, value, callback) => {
        if (value) {
          callback()
        } else {
          callback(new Error('请阅读并同意服务条款和隐私政策'))
        }
      }, 
      trigger: 'change' 
    }
  ]
})

// 处理注册
const handleRegister = () => {
  if (formRef.value) {
    formRef.value.validate(async (valid) => {
      if (valid) {
        // DID注册额外验证
        if (registerMode.value === 'did') {
          if (!keyPair.value) {
            ElMessage.error('请先生成密钥对')
            return
          }
          if (!keysSaved.value) {
            ElMessage.error('请确认已保存私钥')
            return
          }
        }

        loading.value = true
        try {
          if (registerMode.value === 'did') {
            // DID注册
            await didRegister({
              citizenID: registerForm.citizenID,
              name: registerForm.name,
              phone: registerForm.phone,
              email: registerForm.email,
              password: registerForm.password,
              organization: registerForm.organization,
              role: 'user',
              balance: parseFloat(registerForm.balance),
              publicKey: keyPair.value.publicKey
            })
            
            ElMessage.success('DID注册成功！您现在拥有了去中心化身份')
          } else {
            // 传统注册
            await register({
              name: registerForm.name,
              citizenID: registerForm.citizenID,
              password: registerForm.password,
              email: registerForm.email,
              phone: registerForm.phone,
              role: 'user',
              organization: registerForm.organization,
              balance: parseFloat(registerForm.balance),
            })
            
            ElMessage.success('注册成功，请登录')
          }
          
          await router.push('/login')
        } catch (error) {
          ElMessage.error(error.message || '注册失败，请稍后再试')
        } finally {
          loading.value = false
        }
      }
    })
  }
}

// 显示服务条款
const showTerms = () => {
  termsDialogVisible.value = true
}

// 显示隐私政策
const showPrivacy = () => {
  privacyDialogVisible.value = true
}

// 同意服务条款
const agreeTerms = () => {
  registerForm.agreement = true
  termsDialogVisible.value = false
}

// 同意隐私政策
const agreePrivacy = () => {
  registerForm.agreement = true
  privacyDialogVisible.value = false
}
</script>

<style scoped>
.register-container {
  width: 100%;
  max-height: 100vh;
  overflow-y: auto;
  padding: 20px;
  box-sizing: border-box;
}

/* 为了更好的滚动体验，添加自定义滚动条样式 */
.register-container::-webkit-scrollbar {
  width: 8px;
}

.register-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.register-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.register-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.register-title {
  font-size: 24px;
  color: #333;
  text-align: center;
  margin-bottom: 30px;
  top: 0;
  background: white;
  z-index: 10;
  padding: 10px 0;
}

.register-mode-selector {
  margin-bottom: 30px;
  text-align: center;
  top: 80px;
  background: white;
  z-index: 9;
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
}

.mode-group {
  margin-bottom: 10px;
}

.mode-description {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
}

.register-form {
  width: 100%;
}

.did-key-section {
  margin: 20px 0;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  max-height: 60vh;
  overflow-y: auto;
}

/* DID密钥部分的滚动条样式 */
.did-key-section::-webkit-scrollbar {
  width: 6px;
}

.did-key-section::-webkit-scrollbar-track {
  background: #e9ecef;
  border-radius: 3px;
}

.did-key-section::-webkit-scrollbar-thumb {
  background: #6c757d;
  border-radius: 3px;
}

.did-key-section::-webkit-scrollbar-thumb:hover {
  background: #495057;
}

.key-display {
  margin-top: 20px;
}

.key-display .el-alert {
  margin-bottom: 20px;
}

.w-100 {
  width: 100%;
}

.login-link {
  text-align: center;
  margin-top: 10px;
}

.login-link span {
  font-size: 14px;
  color: #666;
}

.terms-content {
  max-height: 400px;
  overflow-y: auto;
  padding: 0 10px;
}

.terms-content h3 {
  font-size: 18px;
  margin-bottom: 15px;
}

.terms-content h4 {
  font-size: 16px;
  margin: 15px 0 8px;
}

.terms-content p, 
.terms-content ul {
  margin-bottom: 10px;
  font-size: 14px;
  line-height: 1.6;
  color: #666;
}

.terms-content ul {
  padding-left: 20px;
}

/* DID相关样式 */
.el-radio-button {
  margin-right: 0;
}

.el-divider {
  margin: 20px 0;
}

.key-display .el-form-item {
  margin-bottom: 20px;
}

.key-display .el-textarea {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.key-display .el-input-group__append {
  padding: 0;
}

.key-display .el-input-group__append .el-button {
  margin: 0;
  border-radius: 0;
}

.key-display .el-input-group__append .el-button:first-child {
  border-right: 1px solid #dcdfe6;
}

/* 响应式设计 */
@media (max-height: 800px) {
  .register-container {
    max-height: 95vh;
  }
  
  .did-key-section {
    max-height: 50vh;
  }
}

@media (max-height: 600px) {
  .register-container {
    max-height: 90vh;
  }
  
  .did-key-section {
    max-height: 40vh;
  }
  
  .register-title {
    font-size: 20px;
    margin-bottom: 20px;
  }
  
  .register-mode-selector {
    margin-bottom: 20px;
    padding: 10px 0;
  }
}

/* 移动端优化 */
@media (max-width: 768px) {
  .register-container {
    padding: 15px;
    max-height: 100vh;
  }
  
  .did-key-section {
    padding: 15px;
    margin: 15px 0;
  }
  
  .key-display .el-textarea {
    font-size: 11px;
  }
  
  .register-title {
    font-size: 18px;
    margin-bottom: 15px;
  }
}

/* 为表单添加平滑滚动 */
.register-form {
  scroll-behavior: smooth;
}

/* 优化密钥显示区域的布局 */
.key-display {
  margin-top: 20px;
}

.key-display .el-form-item:last-child {
  margin-bottom: 0;
}

/* 密钥部分标题样式 */
.key-section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #409eff;
}

.key-section-title .el-icon {
  font-size: 16px;
}
</style> 