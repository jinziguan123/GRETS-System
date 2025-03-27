<template>
  <div class="transaction-audit-container">
    <el-card class="box-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <h3>交易审核</h3>
          <div class="header-buttons">
            <el-button @click="goBack">返回</el-button>
          </div>
        </div>
      </template>
      
      <div v-if="transactionInfo">
        <!-- 交易基本信息 -->
        <el-descriptions title="交易信息" :column="2" border>
          <el-descriptions-item label="交易编号">{{ transactionInfo.transactionID }}</el-descriptions-item>
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
        </el-descriptions>
        
        <!-- 房产信息 -->
        <el-descriptions v-if="realtyInfo" title="房产信息" :column="3" border class="section-mt">
          <el-descriptions-item label="房产地址">{{ realtyInfo.address }}</el-descriptions-item>
          <el-descriptions-item label="房产类型">{{ realtyInfo.realtyType }}</el-descriptions-item>
          <el-descriptions-item label="建筑面积">{{ realtyInfo.area }} 平方米</el-descriptions-item>
          <el-descriptions-item label="房产状态">{{ getRealtyStatusText(realtyInfo.status) }}</el-descriptions-item>
          <el-descriptions-item label="参考价格">{{ formatPrice(realtyInfo.price) }}</el-descriptions-item>
          <el-descriptions-item label="当前所有者">{{ realtyInfo.currentOwnerCitizenID }}</el-descriptions-item>
        </el-descriptions>
        
        <!-- 卖方信息 -->
        <el-descriptions v-if="sellerInfo" title="卖方信息" :column="3" border class="section-mt">
          <el-descriptions-item label="姓名">{{ sellerInfo.name }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ sellerInfo.citizenID }}</el-descriptions-item>
          <el-descriptions-item label="电话">{{ sellerInfo.phone }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ sellerInfo.email }}</el-descriptions-item>
          <el-descriptions-item label="组织">{{ sellerInfo.organization }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ getUserStatusText(sellerInfo.status) }}</el-descriptions-item>
        </el-descriptions>
        
        <!-- 买方信息 -->
        <el-descriptions v-if="buyerInfo" title="买方信息" :column="3" border class="section-mt">
          <el-descriptions-item label="姓名">{{ buyerInfo.name }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ buyerInfo.citizenID }}</el-descriptions-item>
          <el-descriptions-item label="电话">{{ buyerInfo.phone }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ buyerInfo.email }}</el-descriptions-item>
          <el-descriptions-item label="组织">{{ buyerInfo.organization }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ getUserStatusText(buyerInfo.status) }}</el-descriptions-item>
        </el-descriptions>
        
        <!-- 审核表单 -->
        <div class="audit-form-container section-mt">
          <h4>审核意见</h4>
          <el-form :model="auditForm" :rules="auditRules" ref="auditFormRef" label-width="120px">
            <el-form-item label="风险评估" prop="riskAssessment">
              <el-rate
                v-model="auditForm.riskAssessment"
                :colors="riskColors"
                :texts="riskTexts"
                show-text
              ></el-rate>
            </el-form-item>
            
            <el-form-item label="审核结果" prop="result">
              <el-radio-group v-model="auditForm.result">
                <el-radio label="APPROVED">通过</el-radio>
                <el-radio label="REJECTED">拒绝</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="审核意见" prop="opinion">
              <el-input
                v-model="auditForm.opinion"
                type="textarea"
                :rows="4"
                placeholder="请输入审核意见"
              ></el-input>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="submitAudit" :loading="submitLoading">提交审核</el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
        
        <!-- 历史审核记录 -->
        <div v-if="audits.length > 0" class="section-mt">
          <h4>历史审核记录</h4>
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
      </div>
      
      <el-empty v-else description="未找到交易信息"></el-empty>
    </el-card>
    
    <!-- 确认对话框 -->
    <el-dialog
      v-model="confirmDialogVisible"
      :title="auditForm.result === 'APPROVED' ? '确认审核通过' : '确认审核拒绝'"
      width="30%"
    >
      <p>{{ auditForm.result === 'APPROVED' ? '您确定要通过此交易申请吗？' : '您确定要拒绝此交易申请吗？' }}</p>
      <p v-if="auditForm.result === 'APPROVED'">通过后，买卖双方将可以进行后续的交易流程。</p>
      <p v-else>拒绝后，交易将终止，申请人需要重新提交申请。</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="confirmDialogVisible = false">取消</el-button>
          <el-button
            :type="auditForm.result === 'APPROVED' ? 'success' : 'danger'"
            @click="confirmAudit"
          >确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const transactionID = ref(route.params.id)
const loading = ref(true)
const submitLoading = ref(false)
const confirmDialogVisible = ref(false)
const auditFormRef = ref(null)

const transactionInfo = ref(null)
const realtyInfo = ref(null)
const sellerInfo = ref(null)
const buyerInfo = ref(null)
const audits = ref([])

// 风险评估颜色和文字提示
const riskColors = ['#67C23A', '#E6A23C', '#F56C6C']
const riskTexts = ['低风险', '中等风险', '高风险']

// 审核表单
const auditForm = reactive({
  result: 'APPROVED',
  opinion: '',
  riskAssessment: 1
})

// 表单验证规则
const auditRules = reactive({
  result: [
    { required: true, message: '请选择审核结果', trigger: 'change' }
  ],
  opinion: [
    { required: true, message: '请输入审核意见', trigger: 'blur' },
    { min: 5, max: 500, message: '审核意见长度应在5到500个字符之间', trigger: 'blur' }
  ],
  riskAssessment: [
    { required: true, message: '请进行风险评估', trigger: 'change' }
  ]
})

// 获取交易详情
const fetchTransactionDetail = async () => {
  loading.value = true
  try {
    const response = await axios.get(`/api/transaction/${transactionID.value}`)
    transactionInfo.value = response.data
    
    // 获取房产信息
    if (transactionInfo.value.realtyCert) {
      await fetchRealtyDetail(transactionInfo.value.realtyCert)
    }
    
    // 获取卖方信息
    if (transactionInfo.value.sellerCitizenID) {
      await fetchUserInfo(transactionInfo.value.sellerCitizenID, 'seller')
    }
    
    // 获取买方信息
    if (transactionInfo.value.buyerCitizenID) {
      await fetchUserInfo(transactionInfo.value.buyerCitizenID, 'buyer')
    }
    
    // 获取审核记录
    await fetchAudits()
    
    // 检查交易状态是否可审核
    if (transactionInfo.value.status !== 'PENDING') {
      ElMessage.warning('该交易当前状态不允许审核')
    }
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

// 获取用户信息
const fetchUserInfo = async (citizenID, userType) => {
  try {
    const response = await axios.get(`/api/user/citizenid/${citizenID}`)
    if (userType === 'seller') {
      sellerInfo.value = response.data
    } else if (userType === 'buyer') {
      buyerInfo.value = response.data
    }
  } catch (error) {
    console.error(`获取${userType}信息失败:`, error)
  }
}

// 获取审核记录
const fetchAudits = async () => {
  try {
    const response = await axios.get(`/api/audit/transaction/${transactionID.value}`)
    audits.value = response.data || []
  } catch (error) {
    console.error('获取审核记录失败:', error)
  }
}

// 提交审核
const submitAudit = async () => {
  if (!auditFormRef.value) return
  
  await auditFormRef.value.validate(async (valid, fields) => {
    if (valid) {
      // 根据风险评估自动设置审核结果
      if (auditForm.riskAssessment >= 4) {
        ElMessageBox.confirm('当前风险评估较高，建议拒绝交易，是否继续?', '风险提示', {
          confirmButtonText: '继续',
          cancelButtonText: '修改审核结果',
          type: 'warning'
        }).then(() => {
          confirmDialogVisible.value = true
        }).catch(() => {})
      } else {
        confirmDialogVisible.value = true
      }
    } else {
      console.log('表单验证失败:', fields)
      ElMessage.error('请完善表单信息')
    }
  })
}

// 确认审核
const confirmAudit = async () => {
  submitLoading.value = true
  
  try {
    const auditData = {
      transactionID: transactionID.value,
      result: auditForm.result,
      opinion: auditForm.opinion,
      riskLevel: auditForm.riskAssessment
    }
    
    const response = await axios.post('/api/transaction/audit', auditData)
    ElMessage.success(auditForm.result === 'APPROVED' ? '审核通过成功' : '审核拒绝成功')
    confirmDialogVisible.value = false
    
    // 跳转到交易详情页
    router.push(`/transaction/detail/${transactionID.value}`)
  } catch (error) {
    console.error('审核失败:', error)
    ElMessage.error(error.response?.data?.message || '审核失败')
  } finally {
    submitLoading.value = false
  }
}

// 重置表单
const resetForm = () => {
  auditFormRef.value.resetFields()
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

// 获取用户状态文本
const getUserStatusText = (status) => {
  const statusMap = {
    'ACTIVE': '正常',
    'INACTIVE': '未激活',
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

onMounted(() => {
  fetchTransactionDetail()
})
</script>

<style scoped>
.transaction-audit-container {
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

.audit-form-container {
  padding: 20px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  background-color: #f5f7fa;
}

.header-buttons {
  display: flex;
  gap: 10px;
}
</style> 