<template>
  <div class="register-container">
    <h2 class="register-title">用户注册</h2>
    
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
      
      <!-- 组织选择 -->
      <el-form-item prop="organization">
        <el-select
          v-model="registerForm.organization"
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

      <el-form-item>
        <el-input
            v-model="registerForm.balance"
            placeholder="请输入注册金额"
            prefix-icon="money"
        />
      </el-form-item>
      
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
        >
          {{ loading ? '注册中...' : '注册' }}
        </el-button>
      </el-form-item>
      
      <!-- 登录链接 -->
      <el-form-item class="login-link">
        <span>已有账号？</span>
        <el-button type="text" @click="router.push('login')">立即登录</el-button>
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
        </ul>
        
        <h4>2. 信息使用</h4>
        <p>我们使用收集的信息用于：</p>
        <ul>
          <li>提供、维护和改进我们的服务</li>
          <li>处理交易和请求</li>
          <li>发送通知和更新</li>
          <li>预防欺诈和提升安全性</li>
          <li>分析使用趋势和改进用户体验</li>
        </ul>
        
        <h4>3. 信息共享</h4>
        <p>我们不会出售、出租或交易您的个人信息。但在以下情况下，我们可能会共享您的信息：</p>
        <ul>
          <li>经您明确同意</li>
          <li>完成交易所必需的共享</li>
          <li>遵守法律要求、保护权利或安全</li>
        </ul>
        
        <h4>4. 区块链数据</h4>
        <p>请注意，存储在区块链上的交易数据具有不可篡改性。个人敏感信息不会直接存储在区块链上，但交易记录将被永久保存。</p>
        
        <h4>5. 数据安全</h4>
        <p>我们采取合理措施保护您的个人信息，包括加密技术和访问控制。但互联网传输无法保证绝对安全，您使用服务即表示了解这一风险。</p>
        
        <h4>6. 您的权利</h4>
        <p>您有权访问、更正或删除您的个人信息。如需行使这些权利，请通过系统内的设置或联系我们的客服。</p>
        
        <h4>7. 政策变更</h4>
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { register } from './api.js'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const termsDialogVisible = ref(false)
const privacyDialogVisible = ref(false)

// 注册表单数据
const registerForm = reactive({
  name: '',
  citizenID: '',
  password: '',
  confirmPassword: '',
  email: '',
  phone: '',
  organization: '',
  balance: undefined,
  agreement: false
})

// 可选组织列表
const organizations = [
  { label: '政府监管部门', value: 'government' },
  { label: '投资者/买家', value: 'investor' },
  { label: '银行机构', value: 'bank' },
  { label: '第三方机构', value: 'thirdparty' },
  { label: '审计机构', value: 'audit' }
]

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

// 表单验证规则
const registerRules = reactive({
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 3, max: 20, message: '姓名长度应为3-20个字符', trigger: 'blur' }
  ],
  citizenID: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { min: 18, max: 18, message: '身份证号长度应为18个字符', trigger: 'blur' }
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
  organization: [
    { required: true, message: '请选择组织', trigger: 'change' }
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
        loading.value = true
        try {
          // 调用注册接口
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
}

.register-title {
  font-size: 24px;
  color: #333;
  text-align: center;
  margin-bottom: 30px;
}

.register-form {
  width: 100%;
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
</style> 