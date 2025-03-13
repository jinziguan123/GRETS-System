<template>
  <div class="mortgage-create-container">
    <div class="page-header">
      <h1>申请抵押贷款</h1>
      <p>请填写以下信息申请房产抵押贷款</p>
    </div>
    
    <el-form
      ref="formRef"
      :model="mortgageForm"
      :rules="formRules"
      label-width="120px"
      class="mortgage-form"
    >
      <el-card class="form-card">
        <template #header>
          <div class="card-header">
            <span>贷款基本信息</span>
          </div>
        </template>
        
        <!-- 房产信息 -->
        <el-form-item label="抵押房产" prop="realtyId">
          <el-select
            v-model="mortgageForm.realtyId"
            placeholder="请选择要抵押的房产"
            filterable
            class="w-100"
            @change="handleRealtyChange"
          >
            <el-option
              v-for="item in realtyOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            >
              <div class="realty-option">
                <span>{{ item.name }}</span>
                <span class="realty-option-price">价值: ¥{{ item.price.toLocaleString() }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item v-if="selectedRealty" label="房产详情">
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="房产编号">{{ selectedRealty.id }}</el-descriptions-item>
            <el-descriptions-item label="房产价值">¥{{ selectedRealty.price.toLocaleString() }}</el-descriptions-item>
            <el-descriptions-item label="房产地址" :span="2">{{ selectedRealty.address }}</el-descriptions-item>
          </el-descriptions>
        </el-form-item>
        
        <!-- 贷款金额 -->
        <el-form-item label="贷款金额" prop="amount">
          <el-input-number
            v-model="mortgageForm.amount"
            :min="minLoanAmount"
            :max="maxLoanAmount"
            :step="10000"
            :precision="2"
            class="w-100"
          />
          <div class="form-help-text">
            最高可贷款额度为房产价值的70%，即 ¥{{ maxLoanAmount.toLocaleString() }}
          </div>
        </el-form-item>
        
        <!-- 贷款期限 -->
        <el-form-item label="贷款期限" prop="term">
          <el-select v-model="mortgageForm.term" class="w-100" placeholder="请选择贷款期限">
            <el-option label="1年" value="1" />
            <el-option label="3年" value="3" />
            <el-option label="5年" value="5" />
            <el-option label="10年" value="10" />
            <el-option label="15年" value="15" />
            <el-option label="20年" value="20" />
            <el-option label="30年" value="30" />
          </el-select>
        </el-form-item>
        
        <!-- 利率类型 -->
        <el-form-item label="利率类型" prop="rateType">
          <el-radio-group v-model="mortgageForm.rateType">
            <el-radio label="fixed">固定利率</el-radio>
            <el-radio label="variable">浮动利率</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <!-- 年利率 -->
        <el-form-item label="年利率(%)" prop="interestRate">
          <el-input-number
            v-model="mortgageForm.interestRate"
            :min="2"
            :max="10"
            :step="0.01"
            :precision="2"
            class="w-100"
          />
        </el-form-item>
      </el-card>
      
      <el-card class="form-card">
        <template #header>
          <div class="card-header">
            <span>申请人信息</span>
          </div>
        </template>
        
        <!-- 申请人姓名 -->
        <el-form-item label="申请人姓名" prop="applicantName">
          <el-input v-model="mortgageForm.applicantName" placeholder="请输入申请人姓名" />
        </el-form-item>
        
        <!-- 身份证号 -->
        <el-form-item label="身份证号" prop="applicantIdCard">
          <el-input v-model="mortgageForm.applicantIdCard" placeholder="请输入申请人身份证号" />
        </el-form-item>
        
        <!-- 联系电话 -->
        <el-form-item label="联系电话" prop="applicantPhone">
          <el-input v-model="mortgageForm.applicantPhone" placeholder="请输入联系电话" />
        </el-form-item>
        
        <!-- 电子邮箱 -->
        <el-form-item label="电子邮箱" prop="applicantEmail">
          <el-input v-model="mortgageForm.applicantEmail" placeholder="请输入电子邮箱" type="email" />
        </el-form-item>
        
        <!-- 工作单位 -->
        <el-form-item label="工作单位" prop="applicantEmployer">
          <el-input v-model="mortgageForm.applicantEmployer" placeholder="请输入工作单位" />
        </el-form-item>
        
        <!-- 年收入 -->
        <el-form-item label="年收入(元)" prop="applicantIncome">
          <el-input-number
            v-model="mortgageForm.applicantIncome"
            :min="0"
            :step="10000"
            :precision="2"
            class="w-100"
          />
        </el-form-item>
      </el-card>
      
      <el-card class="form-card">
        <template #header>
          <div class="card-header">
            <span>贷款银行信息</span>
          </div>
        </template>
        
        <!-- 贷款银行 -->
        <el-form-item label="贷款银行" prop="bankId">
          <el-select
            v-model="mortgageForm.bankId"
            placeholder="请选择贷款银行"
            class="w-100"
          >
            <el-option
              v-for="bank in bankOptions"
              :key="bank.id"
              :label="bank.name"
              :value="bank.id"
            />
          </el-select>
        </el-form-item>
        
        <!-- 还款方式 -->
        <el-form-item label="还款方式" prop="repaymentMethod">
          <el-select v-model="mortgageForm.repaymentMethod" placeholder="请选择还款方式" class="w-100">
            <el-option label="等额本息" value="equal_installment" />
            <el-option label="等额本金" value="equal_principal" />
            <el-option label="先息后本" value="interest_only" />
            <el-option label="一次性还本付息" value="lump_sum" />
          </el-select>
        </el-form-item>
        
        <!-- 还款账号 -->
        <el-form-item label="还款账号" prop="repaymentAccount">
          <el-input v-model="mortgageForm.repaymentAccount" placeholder="请输入还款账号" />
        </el-form-item>
      </el-card>
      
      <el-card class="form-card">
        <template #header>
          <div class="card-header">
            <span>附加信息</span>
          </div>
        </template>
        
        <!-- 附件上传 -->
        <el-form-item label="身份证复印件" prop="idCardFile">
          <el-upload
            class="upload-container"
            action="/api/file/upload"
            :headers="uploadHeaders"
            :on-success="handleIdCardUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
          >
            <el-button type="primary">点击上传</el-button>
            <template #tip>
              <div class="el-upload__tip">请上传身份证正反面扫描件，JPG/PNG格式，不超过5MB</div>
            </template>
          </el-upload>
        </el-form-item>
        
        <!-- 收入证明 -->
        <el-form-item label="收入证明" prop="incomeProofFile">
          <el-upload
            class="upload-container"
            action="/api/file/upload"
            :headers="uploadHeaders"
            :on-success="handleIncomeProofUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
          >
            <el-button type="primary">点击上传</el-button>
            <template #tip>
              <div class="el-upload__tip">请上传最近6个月的收入证明，PDF格式，不超过10MB</div>
            </template>
          </el-upload>
        </el-form-item>
        
        <!-- 备注 -->
        <el-form-item label="备注说明" prop="remark">
          <el-input
            v-model="mortgageForm.remark"
            type="textarea"
            rows="4"
            placeholder="如有其他需要说明的情况，请在此处填写"
          />
        </el-form-item>
      </el-card>
      
      <!-- 按钮组 -->
      <div class="form-actions">
        <el-button @click="$router.push('/mortgage/list')">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">提交申请</el-button>
      </div>
    </el-form>
    
    <!-- 贷款计算器对话框 -->
    <el-dialog v-model="calculatorVisible" title="贷款计算器" width="60%">
      <el-form :model="calcForm" label-width="120px">
        <el-form-item label="贷款金额">
          <el-input-number v-model="calcForm.amount" :min="10000" :step="10000" :precision="2" />
        </el-form-item>
        <el-form-item label="贷款期限(年)">
          <el-input-number v-model="calcForm.term" :min="1" :max="30" :step="1" />
        </el-form-item>
        <el-form-item label="年利率(%)">
          <el-input-number v-model="calcForm.rate" :min="1" :max="10" :step="0.01" :precision="2" />
        </el-form-item>
        <el-form-item label="还款方式">
          <el-select v-model="calcForm.method">
            <el-option label="等额本息" value="equal_installment" />
            <el-option label="等额本金" value="equal_principal" />
          </el-select>
        </el-form-item>
      </el-form>
      
      <div class="calculator-results" v-if="showResults">
        <el-divider content-position="center">计算结果</el-divider>
        <el-descriptions border :column="2">
          <el-descriptions-item label="月供(首月)">¥{{ monthlyPayment.toLocaleString() }}</el-descriptions-item>
          <el-descriptions-item label="总还款额">¥{{ totalPayment.toLocaleString() }}</el-descriptions-item>
          <el-descriptions-item label="总利息">¥{{ totalInterest.toLocaleString() }}</el-descriptions-item>
          <el-descriptions-item label="还款月数">{{ calcForm.term * 12 }}个月</el-descriptions-item>
        </el-descriptions>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="calculateLoan">计算</el-button>
          <el-button type="primary" @click="applyCalculation">应用到表单</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { realtyApi, mortgageApi, userApi } from '@/api'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const submitting = ref(false)
const calculatorVisible = ref(false)
const showResults = ref(false)

// 表单数据
const mortgageForm = reactive({
  realtyId: '',
  amount: 0,
  term: '10',
  rateType: 'fixed',
  interestRate: 4.9,
  applicantName: '',
  applicantIdCard: '',
  applicantPhone: '',
  applicantEmail: '',
  applicantEmployer: '',
  applicantIncome: 0,
  bankId: '',
  repaymentMethod: 'equal_installment',
  repaymentAccount: '',
  idCardFile: '',
  incomeProofFile: '',
  remark: ''
})

// 计算器数据
const calcForm = reactive({
  amount: 500000,
  term: 10,
  rate: 4.9,
  method: 'equal_installment'
})

// 计算结果
const monthlyPayment = ref(0)
const totalPayment = ref(0)
const totalInterest = ref(0)

// 可选房产列表
const realtyOptions = ref([])
const selectedRealty = ref(null)

// 最小和最大贷款额度
const minLoanAmount = ref(0)
const maxLoanAmount = computed(() => {
  if (selectedRealty.value) {
    return selectedRealty.value.price * 0.7
  }
  return 0
})

// 银行选项
const bankOptions = ref([
  { id: 'bank1', name: '中国工商银行' },
  { id: 'bank2', name: '中国建设银行' },
  { id: 'bank3', name: '中国农业银行' },
  { id: 'bank4', name: '中国银行' },
  { id: 'bank5', name: '交通银行' }
])

// 上传相关
const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${userStore.token}`
  }
})

// 验证手机号码
const validatePhone = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入联系电话'))
  } else {
    const reg = /^1[3-9]\d{9}$/
    if (reg.test(value)) {
      callback()
    } else {
      callback(new Error('请输入有效的手机号码'))
    }
  }
}

// 验证身份证号
const validateIdCard = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入身份证号'))
  } else {
    const reg = /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/
    if (reg.test(value)) {
      callback()
    } else {
      callback(new Error('请输入有效的身份证号'))
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

// 表单校验规则
const formRules = reactive({
  realtyId: [
    { required: true, message: '请选择抵押房产', trigger: 'change' }
  ],
  amount: [
    { required: true, message: '请输入贷款金额', trigger: 'blur' },
    { type: 'number', min: 10000, message: '贷款金额不得低于10000元', trigger: 'blur' }
  ],
  term: [
    { required: true, message: '请选择贷款期限', trigger: 'change' }
  ],
  rateType: [
    { required: true, message: '请选择利率类型', trigger: 'change' }
  ],
  interestRate: [
    { required: true, message: '请设置年利率', trigger: 'blur' }
  ],
  applicantName: [
    { required: true, message: '请输入申请人姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在2-20个字符之间', trigger: 'blur' }
  ],
  applicantIdCard: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { validator: validateIdCard, trigger: 'blur' }
  ],
  applicantPhone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { validator: validatePhone, trigger: 'blur' }
  ],
  applicantEmail: [
    { required: true, message: '请输入电子邮箱', trigger: 'blur' },
    { validator: validateEmail, trigger: 'blur' }
  ],
  applicantEmployer: [
    { required: true, message: '请输入工作单位', trigger: 'blur' }
  ],
  applicantIncome: [
    { required: true, message: '请输入年收入', trigger: 'blur' },
    { type: 'number', min: 0, message: '年收入不能为负数', trigger: 'blur' }
  ],
  bankId: [
    { required: true, message: '请选择贷款银行', trigger: 'change' }
  ],
  repaymentMethod: [
    { required: true, message: '请选择还款方式', trigger: 'change' }
  ],
  repaymentAccount: [
    { required: true, message: '请输入还款账号', trigger: 'blur' }
  ]
})

// 初始化
onMounted(async () => {
  try {
    // 获取房产列表
    const res = await realtyApi.getRealtyList({
      owner: userStore.userId,
      status: 'available'
    })
    
    // 处理返回数据
    realtyOptions.value = res.data.map(realty => ({
      id: realty.id,
      name: realty.title || `${realty.type} - ${realty.address}`,
      price: realty.price,
      address: realty.address,
      area: realty.area,
      type: realty.type
    }))
    
    // 获取用户信息填充表单
    const userInfo = await userApi.getUserProfile()
    if (userInfo) {
      mortgageForm.applicantName = userInfo.fullName || userInfo.username
      mortgageForm.applicantPhone = userInfo.phone || ''
      mortgageForm.applicantEmail = userInfo.email || ''
    }
  } catch (error) {
    console.error('初始化数据失败:', error)
    ElMessage.error('获取数据失败，请刷新页面重试')
    
    // 使用模拟数据（仅用于演示）
    realtyOptions.value = [
      { id: 'R001', name: '阳光花园 3室2厅', price: 1200000, address: '杭州市西湖区文三路138号阳光花园3幢1单元601', area: 120, type: '住宅' },
      { id: 'R002', name: '江南名府 4室2厅', price: 2500000, address: '杭州市滨江区滨盛路1509号江南名府12幢2单元1801', area: 180, type: '住宅' },
      { id: 'R003', name: '城市广场 商铺A12', price: 3500000, address: '杭州市江干区凯旋路166号城市广场A区12号', area: 85, type: '商铺' }
    ]
  }
})

// 选择房产时的处理
const handleRealtyChange = (realtyId) => {
  selectedRealty.value = realtyOptions.value.find(item => item.id === realtyId)
  if (selectedRealty.value) {
    minLoanAmount.value = Math.min(100000, selectedRealty.value.price * 0.1)
    mortgageForm.amount = Math.min(selectedRealty.value.price * 0.6, maxLoanAmount.value)
    
    // 更新计算器金额
    calcForm.amount = mortgageForm.amount
  }
}

// 上传文件前的处理
const beforeUpload = (file) => {
  // 检查文件大小和类型
  const isImage = file.type.startsWith('image/')
  const isPDF = file.type === 'application/pdf'
  const isLt10M = file.size / 1024 / 1024 < 10
  
  if (!isImage && !isPDF) {
    ElMessage.error('只能上传图片或PDF文件!')
    return false
  }
  
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过10MB!')
    return false
  }
  
  return true
}

// 身份证上传成功
const handleIdCardUploadSuccess = (response) => {
  if (response.success) {
    ElMessage.success('身份证文件上传成功')
    mortgageForm.idCardFile = response.data.fileId
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

// 收入证明上传成功
const handleIncomeProofUploadSuccess = (response) => {
  if (response.success) {
    ElMessage.success('收入证明上传成功')
    mortgageForm.incomeProofFile = response.data.fileId
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

// 上传出错
const handleUploadError = () => {
  ElMessage.error('文件上传失败，请重试')
}

// 计算贷款
const calculateLoan = () => {
  // 计算月供、总还款额和总利息
  const principal = calcForm.amount
  const months = calcForm.term * 12
  const monthlyRate = calcForm.rate / 100 / 12
  
  if (calcForm.method === 'equal_installment') {
    // 等额本息
    const monthly = principal * monthlyRate * Math.pow(1 + monthlyRate, months) / (Math.pow(1 + monthlyRate, months) - 1)
    monthlyPayment.value = parseFloat(monthly.toFixed(2))
    totalPayment.value = parseFloat((monthly * months).toFixed(2))
    totalInterest.value = parseFloat((totalPayment.value - principal).toFixed(2))
  } else {
    // 等额本金
    const monthlyPrincipal = principal / months
    const firstMonthInterest = principal * monthlyRate
    const firstMonthPayment = monthlyPrincipal + firstMonthInterest
    
    monthlyPayment.value = parseFloat(firstMonthPayment.toFixed(2))
    
    // 计算总利息
    let totalInt = 0
    for (let i = 0; i < months; i++) {
      const remainingPrincipal = principal - (monthlyPrincipal * i)
      totalInt += remainingPrincipal * monthlyRate
    }
    
    totalInterest.value = parseFloat(totalInt.toFixed(2))
    totalPayment.value = parseFloat((principal + totalInterest.value).toFixed(2))
  }
  
  showResults.value = true
}

// 应用计算结果到表单
const applyCalculation = () => {
  mortgageForm.amount = calcForm.amount
  mortgageForm.term = calcForm.term.toString()
  mortgageForm.interestRate = calcForm.rate
  mortgageForm.repaymentMethod = calcForm.method
  
  calculatorVisible.value = false
  ElMessage.success('已应用计算结果到申请表单')
}

// 提交表单
const submitForm = () => {
  if (formRef.value) {
    formRef.value.validate(async (valid) => {
      if (valid) {
        try {
          submitting.value = true
          
          // 确认提交
          await ElMessageBox.confirm(
            '提交后，您的贷款申请将被送交银行审批，请确认信息无误',
            '确认提交',
            {
              confirmButtonText: '确认提交',
              cancelButtonText: '取消',
              type: 'warning'
            }
          )
          
          // 调用API提交贷款申请
          const result = await mortgageApi.createMortgage({
            realty_id: mortgageForm.realtyId,
            amount: mortgageForm.amount,
            term: parseInt(mortgageForm.term),
            rate_type: mortgageForm.rateType,
            interest_rate: mortgageForm.interestRate,
            applicant_name: mortgageForm.applicantName,
            applicant_id_card: mortgageForm.applicantIdCard,
            applicant_phone: mortgageForm.applicantPhone,
            applicant_email: mortgageForm.applicantEmail,
            applicant_employer: mortgageForm.applicantEmployer,
            applicant_income: mortgageForm.applicantIncome,
            bank_id: mortgageForm.bankId,
            repayment_method: mortgageForm.repaymentMethod,
            repayment_account: mortgageForm.repaymentAccount,
            id_card_file: mortgageForm.idCardFile,
            income_proof_file: mortgageForm.incomeProofFile,
            remark: mortgageForm.remark
          })
          
          ElMessage.success('贷款申请提交成功')
          
          // 跳转到贷款详情页
          router.push({
            path: `/mortgage/detail/${result.data.id}`
          })
        } catch (error) {
          if (error === 'cancel') {
            return
          }
          console.error('提交贷款申请失败:', error)
          ElMessage.error(error.message || '提交失败，请稍后重试')
        } finally {
          submitting.value = false
        }
      } else {
        ElMessage.warning('请完善表单信息')
      }
    })
  }
}
</script>

<style scoped>
.mortgage-create-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: bold;
  margin: 0 0 8px;
}

.page-header p {
  margin: 0;
  color: #666;
}

.form-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  font-weight: bold;
}

.w-100 {
  width: 100%;
}

.form-help-text {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.realty-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.realty-option-price {
  color: #409EFF;
  font-weight: bold;
}

.upload-container {
  width: 100%;
}

.form-actions {
  display: flex;
  justify-content: center;
  margin-top: 30px;
  gap: 20px;
}

.calculator-results {
  margin-top: 20px;
}

@media (max-width: 768px) {
  .mortgage-create-container {
    padding: 15px;
  }
  
  .el-form {
    padding: 0;
  }
}
</style> 