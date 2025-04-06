<template>
  <div class="transaction-detail-container">
    <el-card class="box-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <h3>交易详情</h3>
          <div class="header-buttons">
            <el-button @click="goBack">返回</el-button>
            <el-button 
              v-if="canCheckTransaction" 
              type="success" 
              @click="handleAudit('APPROVED')"
            >审核通过</el-button>
            <el-button 
              v-if="canCheckTransaction" 
              type="danger" 
              @click="handleAudit('REJECTED')"
            >审核拒绝</el-button>
            <el-button 
              v-if="canCompleteTransaction" 
              type="primary" 
              @click="handleComplete"
            >完成交易</el-button>
            <el-button 
              v-if="canCancelTransaction" 
              type="warning" 
              @click="handleCancel"
            >取消交易</el-button>
          </div>
        </div>
      </template>
      
      <div v-if="transactionInfo">
        <!-- 交易基本信息 -->
        <el-descriptions title="交易基本信息" :column="2" border>
          <el-descriptions-item label="交易编号">{{ transactionInfo.transactionUUID }}</el-descriptions-item>
          <el-descriptions-item label="交易状态">
            <el-tag :type="getStatusType(transactionInfo.status)">
              {{ getStatusText(transactionInfo.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="房产证号">{{ transactionInfo.realtyCert }}</el-descriptions-item>
          <el-descriptions-item label="交易价格">{{ formatPrice(transactionInfo.price) }}</el-descriptions-item>
          <el-descriptions-item label="卖方身份证">{{ transactionInfo.sellerCitizenID }}</el-descriptions-item>
          <el-descriptions-item label="买方身份证">{{ transactionInfo.buyerCitizenID }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(transactionInfo.createTime) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDate(transactionInfo.updateTime) }}</el-descriptions-item>
        </el-descriptions>
        
        <!-- 房产信息 -->
        <el-descriptions v-if="realtyInfo" title="房产信息" :column="3" border class="section-mt">
          <el-descriptions-item label="房产地址">{{ realtyInfo.address }}</el-descriptions-item>
          <el-descriptions-item label="房产类型">{{ realtyInfo.realtyType }}</el-descriptions-item>
          <el-descriptions-item label="建筑面积">{{ realtyInfo.area }} 平方米</el-descriptions-item>
          <el-descriptions-item label="房产状态">{{ getRealtyStatusText(realtyInfo.status) }}</el-descriptions-item>
          <el-descriptions-item label="参考价格">{{ formatPrice(realtyInfo.price) }}</el-descriptions-item>
          <el-descriptions-item>
            <template #label>
              <el-button type="primary" size="small" @click="viewRealtyDetail">查看房产详情</el-button>
            </template>
          </el-descriptions-item>
        </el-descriptions>
        
        <!-- 审核记录 -->
        <div v-if="audits.length > 0" class="section-mt">
          <h4>审核记录</h4>
          <el-timeline>
            <el-timeline-item
              v-for="(audit, index) in audits"
              :key="index"
              :timestamp="formatDate(audit.createTime)"
              :type="getAuditTypeIcon(audit.result)"
              :color="getAuditTypeColor(audit.result)"
            >
              <h5>{{ getAuditResultText(audit.result) }}</h5>
              <p>审核人: {{ audit.auditorName }}</p>
              <p>组织: {{ audit.auditorOrg }}</p>
              <p>意见: {{ audit.opinion || '无' }}</p>
            </el-timeline-item>
          </el-timeline>
        </div>
        
        <!-- 支付记录 -->
        <div v-if="payments.length > 0" class="section-mt">
          <h4>支付记录</h4>
          <el-table :data="payments" border style="width: 100%">
            <el-table-column prop="paymentID" label="支付编号" width="220"></el-table-column>
            <el-table-column prop="paymentType" label="支付类型" width="120">
              <template #default="scope">
                {{ getPaymentTypeText(scope.row.paymentType) }}
              </template>
            </el-table-column>
            <el-table-column prop="amount" label="支付金额" width="150">
              <template #default="scope">
                {{ formatPrice(scope.row.amount) }}
              </template>
            </el-table-column>
            <el-table-column prop="payerCitizenID" label="付款方" width="180"></el-table-column>
            <el-table-column prop="payeeCitizenID" label="收款方" width="180"></el-table-column>
            <el-table-column prop="createTime" label="支付时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.createTime) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <!-- 合同信息 -->
        <div v-if="contracts.length > 0" class="section-mt">
          <h4>合同信息</h4>
          <el-table :data="contracts" border style="width: 100%">
            <el-table-column prop="contractID" label="合同编号" width="220"></el-table-column>
            <el-table-column prop="contractType" label="合同类型" width="120">
              <template #default="scope">
                {{ getContractTypeText(scope.row.contractType) }}
              </template>
            </el-table-column>
            <el-table-column prop="filePath" label="合同文件" width="300">
              <template #default="scope">
                <el-button v-if="scope.row.filePath" type="primary" link @click="viewContract(scope.row)">
                  查看合同
                </el-button>
                <span v-else>无文件</span>
              </template>
            </el-table-column>
            <el-table-column prop="createTime" label="创建时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.createTime) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <!-- 操作区域 -->
        <div v-if="transactionInfo.status === 'APPROVED'" class="action-area section-mt">
          <h4>交易操作</h4>
          <el-row :gutter="20">
            <el-col :span="8">
              <el-card shadow="hover" class="action-card">
                <h5>上传买卖合同</h5>
                <el-upload
                  class="upload-area"
                  action="/api/contract/upload"
                  :headers="uploadHeaders"
                  :data="{ transactionUUID: transactionUUID, contractType: 'SALE' }"
                  :on-success="handleUploadSuccess"
                  :on-error="handleUploadError"
                  :limit="1"
                >
                  <el-button type="primary">上传合同</el-button>
                  <template #tip>
                    <div class="el-upload__tip">请上传PDF格式文件</div>
                  </template>
                </el-upload>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover" class="action-card">
                <h5>缴纳交易税费</h5>
                <el-button type="primary" @click="showPaymentDialog('TAX')">缴纳税费</el-button>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover" class="action-card">
                <h5>支付房款</h5>
                <el-button type="primary" @click="showPaymentDialog('PAYMENT')">支付房款</el-button>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </div>
      
      <el-empty v-else description="未找到交易信息"></el-empty>
    </el-card>
    
    <!-- 支付对话框 -->
    <el-dialog
      v-model="paymentDialogVisible"
      :title="paymentType === 'TAX' ? '缴纳税费' : '支付房款'"
      width="500px"
    >
      <el-form :model="paymentForm" :rules="paymentRules" ref="paymentFormRef" label-width="100px">
        <el-form-item label="支付金额" prop="amount">
          <el-input-number v-model="paymentForm.amount" :min="0" :precision="2" :step="1000" style="width: 100%"></el-input-number>
        </el-form-item>
        <el-form-item label="付款方" prop="payerCitizenID">
          <el-input v-model="paymentForm.payerCitizenID"></el-input>
        </el-form-item>
        <el-form-item label="收款方" prop="payeeCitizenID">
          <el-input v-model="paymentForm.payeeCitizenID"></el-input>
        </el-form-item>
        <el-form-item label="备注" prop="description">
          <el-input v-model="paymentForm.description" type="textarea" :rows="3"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="paymentDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPayment" :loading="submitLoading">确认支付</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 审核对话框 -->
    <el-dialog
      v-model="auditDialogVisible"
      :title="auditAction === 'APPROVED' ? '审核通过' : '审核拒绝'"
      width="500px"
    >
      <el-form :model="auditForm" :rules="auditRules" ref="auditFormRef" label-width="100px">
        <el-form-item label="审核意见" prop="opinion">
          <el-input v-model="auditForm.opinion" type="textarea" :rows="4"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="auditDialogVisible = false">取消</el-button>
          <el-button 
            :type="auditAction === 'APPROVED' ? 'success' : 'danger'" 
            @click="submitAudit" 
            :loading="submitLoading"
          >
            {{ auditAction === 'APPROVED' ? '通过' : '拒绝' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import {getTransactionDetail} from "@/api/transaction.js";

const router = useRouter()
const route = useRoute()
const transactionUUID = ref(route.params.id)
const loading = ref(true)
const transactionInfo = ref(null)
const realtyInfo = ref(null)
const audits = ref([])
const payments = ref([])
const contracts = ref([])
const submitLoading = ref(false)

// 支付对话框
const paymentDialogVisible = ref(false)
const paymentFormRef = ref(null)
const paymentType = ref('')
const paymentForm = reactive({
  amount: 0,
  payerCitizenID: '',
  payeeCitizenID: '',
  description: ''
})
const paymentRules = reactive({
  amount: [
    { required: true, message: '请输入支付金额', trigger: 'blur' },
    { type: 'number', min: 1, message: '金额必须大于0', trigger: 'blur' }
  ],
  payerCitizenID: [
    { required: true, message: '请输入付款方身份证号', trigger: 'blur' },
    { pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, message: '请输入正确的身份证号', trigger: 'blur' }
  ],
  payeeCitizenID: [
    { required: true, message: '请输入收款方身份证号', trigger: 'blur' },
    { pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, message: '请输入正确的身份证号', trigger: 'blur' }
  ]
})

// 审核对话框
const auditDialogVisible = ref(false)
const auditFormRef = ref(null)
const auditAction = ref('')
const auditForm = reactive({
  opinion: ''
})
const auditRules = reactive({
  opinion: [
    { required: true, message: '请输入审核意见', trigger: 'blur' }
  ]
})

// 上传文件的请求头
const uploadHeaders = computed(() => {
  return {
    Authorization: localStorage.getItem('token') || ''
  }
})

// 获取用户信息
const userInfo = computed(() => {
  const userJson = localStorage.getItem('user')
  if (!userJson) return null
  try {
    return JSON.parse(userJson)
  } catch (e) {
    return null
  }
})

// 判断是否可以审核交易
const canCheckTransaction = computed(() => {
  if (!userInfo.value || !transactionInfo.value) return false
  
  // 政府组织且交易状态为待审核
  return userInfo.value.role === 'GOVERNMENT' && transactionInfo.value.status === 'PENDING'
})

// 判断是否可以完成交易
const canCompleteTransaction = computed(() => {
  if (!userInfo.value || !transactionInfo.value) return false
  
  // 政府组织且交易状态为已审核
  return userInfo.value.role === 'GOVERNMENT' && transactionInfo.value.status === 'APPROVED'
})

// 判断是否可以取消交易
const canCancelTransaction = computed(() => {
  if (!userInfo.value || !transactionInfo.value) return false
  
  // 任何角色在交易未完成前均可取消
  return ['PENDING', 'APPROVED'].includes(transactionInfo.value.status)
})

// 获取交易详情
const fetchTransactionDetail = async () => {
  loading.value = true
  try {
    const response = await getTransactionDetail(transactionUUID.value)
    transactionInfo.value = response.data
    
    // 获取房产信息
    if (transactionInfo.value.realtyCert) {
      await fetchRealtyDetail(transactionInfo.value.realtyCert)
    }
    
    // 获取审核记录
    await fetchAudits()
    
    // 获取支付记录
    await fetchPayments()
    
    // 获取合同信息
    await fetchContracts()
  } catch (error) {
    console.error('获取交易详情失败:', error)
    ElMessage.error(error.response?.data?.message || '获取交易详情失败')
  } finally {
    loading.value = false
  }
}

// 获取房产详情
const fetchRealtyDetail = async (realtyCert) => {
  try {
    const response = await axios.get(`/api/realty/cert/${realtyCert}`)
    realtyInfo.value = response.data
  } catch (error) {
    console.error('获取房产详情失败:', error)
  }
}

// 获取审核记录
const fetchAudits = async () => {
  try {
    const response = await axios.get(`/api/audit/transaction/${transactionUUID.value}`)
    audits.value = response.data || []
  } catch (error) {
    console.error('获取审核记录失败:', error)
  }
}

// 获取支付记录
const fetchPayments = async () => {
  try {
    const response = await axios.get(`/api/payment/transaction/${transactionUUID.value}`)
    payments.value = response.data || []
  } catch (error) {
    console.error('获取支付记录失败:', error)
  }
}

// 获取合同信息
const fetchContracts = async () => {
  try {
    const response = await axios.get(`/api/contract/transaction/${transactionUUID.value}`)
    contracts.value = response.data || []
  } catch (error) {
    console.error('获取合同信息失败:', error)
  }
}

// 显示支付对话框
const showPaymentDialog = (type) => {
  paymentType.value = type
  paymentForm.amount = type === 'TAX' ? calculateTax() : transactionInfo.value.price
  
  // 根据支付类型设置默认付款方和收款方
  if (type === 'TAX') {
    paymentForm.payerCitizenID = transactionInfo.value.buyerCitizenID
    paymentForm.payeeCitizenID = 'GOVERNMENT'
    paymentForm.description = '交易税费'
  } else {
    paymentForm.payerCitizenID = transactionInfo.value.buyerCitizenID
    paymentForm.payeeCitizenID = transactionInfo.value.sellerCitizenID
    paymentForm.description = '房产交易款'
  }
  
  paymentDialogVisible.value = true
}

// 计算税费（示例：交易金额的3%）
const calculateTax = () => {
  if (!transactionInfo.value || !transactionInfo.value.price) return 0
  return transactionInfo.value.price * 0.03
}

// 提交支付
const submitPayment = async () => {
  if (!paymentFormRef.value) return
  
  await paymentFormRef.value.validate(async (valid, fields) => {
    if (valid) {
      submitLoading.value = true
      
      try {
        const paymentData = {
          transactionUUID: transactionUUID.value,
          paymentType: paymentType.value,
          amount: paymentForm.amount,
          payerCitizenID: paymentForm.payerCitizenID,
          payeeCitizenID: paymentForm.payeeCitizenID,
          description: paymentForm.description
        }
        
        const response = await axios.post('/api/payment/create', paymentData)
        ElMessage.success('支付成功')
        paymentDialogVisible.value = false
        
        // 刷新支付记录
        await fetchPayments()
      } catch (error) {
        console.error('支付失败:', error)
        ElMessage.error(error.response?.data?.message || '支付失败')
      } finally {
        submitLoading.value = false
      }
    } else {
      console.log('表单验证失败:', fields)
      ElMessage.error('请完善表单信息')
    }
  })
}

// 处理审核
const handleAudit = (action) => {
  auditAction.value = action
  auditForm.opinion = ''
  auditDialogVisible.value = true
}

// 提交审核
const submitAudit = async () => {
  if (!auditFormRef.value) return
  
  await auditFormRef.value.validate(async (valid, fields) => {
    if (valid) {
      submitLoading.value = true
      
      try {
        const auditData = {
          transactionUUID: transactionUUID.value,
          result: auditAction.value,
          opinion: auditForm.opinion
        }
        
        const response = await axios.post('/api/transaction/audit', auditData)
        ElMessage.success(auditAction.value === 'APPROVED' ? '审核通过成功' : '审核拒绝成功')
        auditDialogVisible.value = false
        
        // 刷新交易信息和审核记录
        await fetchTransactionDetail()
      } catch (error) {
        console.error('审核失败:', error)
        ElMessage.error(error.response?.data?.message || '审核失败')
      } finally {
        submitLoading.value = false
      }
    } else {
      console.log('表单验证失败:', fields)
      ElMessage.error('请完善表单信息')
    }
  })
}

// 完成交易
const handleComplete = async () => {
  ElMessageBox.confirm('确定要完成此次交易吗？完成后将不可撤销', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    loading.value = true
    
    try {
      const response = await axios.post(`/api/transaction/${transactionUUID.value}/complete`)
      ElMessage.success('交易已完成')
      await fetchTransactionDetail()
    } catch (error) {
      console.error('完成交易失败:', error)
      ElMessage.error(error.response?.data?.message || '完成交易失败')
    } finally {
      loading.value = false
    }
  }).catch(() => {})
}

// 取消交易
const handleCancel = async () => {
  ElMessageBox.confirm('确定要取消此次交易吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    loading.value = true
    
    try {
      const response = await axios.post(`/api/transaction/${transactionUUID.value}/cancel`)
      ElMessage.success('交易已取消')
      await fetchTransactionDetail()
    } catch (error) {
      console.error('取消交易失败:', error)
      ElMessage.error(error.response?.data?.message || '取消交易失败')
    } finally {
      loading.value = false
    }
  }).catch(() => {})
}

// 上传合同成功处理
const handleUploadSuccess = (response, file) => {
  ElMessage.success('合同上传成功')
  fetchContracts()
}

// 上传合同失败处理
const handleUploadError = (error, file) => {
  console.error('合同上传失败:', error)
  ElMessage.error('合同上传失败')
}

// 查看合同
const viewContract = (contract) => {
  window.open(contract.filePath, '_blank')
}

// 查看房产详情
const viewRealtyDetail = () => {
  if (realtyInfo.value) {
    router.push(`/realty/detail/${realtyInfo.value.id}`)
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 价格格式化
const formatPrice = (price) => {
  return `¥ ${parseFloat(price).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })}`
}

// 日期格式化
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取交易状态标签类型
const getStatusType = (status) => {
  const statusMap = {
    'PENDING': 'info',
    'APPROVED': 'success',
    'REJECTED': 'danger',
    'IN_PROGRESS': 'warning',
    'COMPLETED': 'success',
    'CANCELLED': 'info'
  }
  return statusMap[status] || 'info'
}

// 获取交易状态文本
const getStatusText = (status) => {
  const statusMap = {
    'PENDING': '待审核',
    'APPROVED': '已审核',
    'REJECTED': '已拒绝',
    'IN_PROGRESS': '进行中',
    'COMPLETED': '已完成',
    'CANCELLED': '已取消'
  }
  return statusMap[status] || status
}

// 获取房产状态文本
const getRealtyStatusText = (status) => {
  const statusMap = {
    'NORMAL': '正常',
    'IN_TRANSACTION': '交易中',
    'MORTGAGED': '已抵押',
    'FROZEN': '已冻结'
  }
  return statusMap[status] || status
}

// 获取审核结果图标
const getAuditTypeIcon = (result) => {
  return result === 'APPROVED' ? 'success' : 'danger'
}

// 获取审核结果颜色
const getAuditTypeColor = (result) => {
  return result === 'APPROVED' ? '#67C23A' : '#F56C6C'
}

// 获取审核结果文本
const getAuditResultText = (result) => {
  return result === 'APPROVED' ? '审核通过' : '审核拒绝'
}

// 获取支付类型文本
const getPaymentTypeText = (type) => {
  const typeMap = {
    'PAYMENT': '房款',
    'TAX': '税费',
    'FEE': '手续费'
  }
  return typeMap[type] || type
}

// 获取合同类型文本
const getContractTypeText = (type) => {
  const typeMap = {
    'SALE': '买卖合同',
    'MORTGAGE': '抵押合同',
    'OTHER': '其他合同'
  }
  return typeMap[type] || type
}

onMounted(() => {
  fetchTransactionDetail()
})
</script>

<style scoped>
.transaction-detail-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-mt {
  margin-top: 30px;
}

.action-area {
  margin-bottom: 20px;
}

.action-card {
  height: 150px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 10px;
}

.header-buttons {
  display: flex;
  gap: 10px;
}
</style>
